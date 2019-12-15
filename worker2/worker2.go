package worker2

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/ossn/fixme_backend/models"

	"github.com/gobuffalo/pop/nulls"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	//"golang.org/x/oauth2"
)


type (
	Worker struct {
		ctx context.Context
	}
)


var (
	client *gitlab.Client
)

func (w *Worker) Init(ctx context.Context, c <-chan os.Signal) {
	w.ctx = ctx
	token := os.Getenv("GITLAB_TOKEN")	// get GITLAB TOKEN as Environment Variable
	if len(token) < 1 {
		panic("Please provide a gitlab token")
	}

	client = gitlab.NewClient(nil, token)
	go w.startPolling(c)
}


func (w *Worker) startPolling(c <-chan os.Signal) {
	// Handle keyboard interupt
	go func() {
		<-c
		os.Exit(1)
	}()
	// Start topics polling
	go w.repositoryTopicsPolling()

	// Start issue polling
	for {
		w.getInitialIssues()
	}
}


// Func to start repo topics polling
func (w *Worker) repositoryTopicsPolling() {
	for {
		go w.UpdateRepositoryTopics()
		time.Sleep(1 * time.Hour)
	}
}

// Get all the tags repositories and set them to the project
func (w *Worker) UpdateRepositoryTopics() {
	w.waitUntilLimitIsRefreshed()
	repos := models.Repositories{}
	err := models.DB.Where("is_github = ?", false).All(&repos)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}

	// map ProjectID -> RepositoriesArray (One project can contain more than one repositories)
	repoIndexMap := make(map[uuid.UUID][]int, len(repos))

	for i, repo := range repos {
		repoIndexMap[repo.ProjectID] = append(repoIndexMap[repo.ProjectID], i)

		tags, _, err := client.Tags.ListTags(repo.ProjectID, nil)

		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repos from gitlab"))
			continue
		}

		repoTags := []string{}
		for _, tag := range tags {
			repoTags = append(repoTags, tag.Name)
		}
		repo.Tags = cleanupArray(repoTags)

		verr, err := repo.Validate(models.DB)
		if verr.HasAny() {
			fmt.Println(verr.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't save repos from gitlab"))
			continue
		}
		repos[i] = repo
	}

	verr, err := models.DB.ValidateAndUpdate(&repos)
		if err != nil || verr.HasAny() {
		fmt.Println(err, verr.Error())
	}


	projects := models.Projects{}
	err = models.DB.Where("is_github = ?", false).All(&projects)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}

	for i, project := range projects {
		repoIDs, exists := repoIndexMap[project.ID]
		projectRepos := models.Repositories{}
		if !exists {
			err = models.DB.Where("project_id = ?", project.ID).All(&projectRepos)
			if err != nil {
				fmt.Println(errors.Wrap(err, "failed to find repos"))
				continue
			}
		} else {
			for _, index := range repoIDs {
				projectRepos = append(projectRepos, repos[index])
			}
		}

		tags := []string{}
		for _, repo := range projectRepos {
			tags = append(tags, repo.Tags...)
		}

		tags = cleanupArray(tags)
		project.Tags = tags

		verr, err := project.Validate(models.DB)
		if verr.HasAny() {
			fmt.Println(verr.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.WithMessage(err, "failed to save project"))
			continue
		}
		projects[i] = project
	}

	verr, err = models.DB.ValidateAndUpdate(&repos)
		if err != nil || verr.HasAny() {
		fmt.Println(err, verr.Error())
	}
}

// waitUntilLimitIsRefreshed: A function that waits until the next github query can be executed
func (w *Worker) waitUntilLimitIsRefreshed () {
			time.Sleep(time.Second*30)	// we don't know the limits
			w.waitUntilLimitIsRefreshed()
}

// Get first issues
func (w *Worker) getInitialIssues() {
	w.waitUntilLimitIsRefreshed()
	lastUpdatedRepo := models.Repository{}
	err := models.DB.Order("last_parsed asc").First(&lastUpdatedRepo)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to get issues"))
		return
	}


	//Find primary language
	languages, _, err := client.Projects.GetProjectLanguages(lastUpdatedRepo.ProjectID, nil)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't find language"))
		return
	}

	var max float32 = -1.0
	var primaryLanguage = ""

	for language, pososto := range *languages {
		if(pososto > max) {
			primaryLanguage = language
			max = pososto
		}
	}

	// Find issues
	openState := "opened"
	issuesOptions := &gitlab.ListProjectIssuesOptions{
		State: &openState,
	}

	issueData, response, err := client.Issues.ListProjectIssues(lastUpdatedRepo.ProjectID, issuesOptions)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't load initial issues"))
		return
	}

	// full pages (20) of issues
	for len(issueData) == 20 {
		go w.parseAndSaveIssues(issueData, &lastUpdatedRepo, &primaryLanguage)
		issuesOptions.ListOptions.Page = response.NextPage
		issueData, response, err = client.Issues.ListProjectIssues(lastUpdatedRepo.ProjectID, issuesOptions)
	}

	// last page of issues
	if len(issueData) > 0 {
		go w.parseAndSaveIssues(issueData, &lastUpdatedRepo, &primaryLanguage)
	}
}


// Parse and save gitlab issues
func (w *Worker) parseAndSaveIssues(issueData []*gitlab.Issue, repository *models.Repository, language *string) {
	issuesToCreate := models.Issues{}
	issuesToUpdate := models.Issues{}
	for _, node := range issueData {
		gitlabIssue := &models.Issue{
			IssueID:      node.ID,
			Body:         nulls.String{String: node.Description, Valid: node.Description != ""},
			Title:        nulls.String{String: node.Title, Valid: node.Title != ""},
			Closed:       node.State == "closed",
			Number:       node.IID,
			URL:          node.WebURL,
			RepositoryID: repository.ID,
			ProjectID:    repository.ProjectID,
			Language:     nulls.String{String: strings.ToLower(*language), Valid: *language != ""},
		}

		// Parse gitlab labels
		labels := []string{}
		for _, label := range node.Labels {
			name := &label
			labels = append(labels, *name)
			// Search for known labels
			matched := searchForMatchingLabels(name, gitlabIssue)
			// Split name based on known delimeters
			tmp := strings.FieldsFunc(*name, split)
			// If label hasn't been matched try again with the splited string
			if !matched && len(tmp) > 1 {
				for _, label := range tmp {
					searchForMatchingLabels(&label, gitlabIssue)
				}
			}
		}

		gitlabIssue.Labels = labels
		// Initialize experience needed with moderate
		if !gitlabIssue.ExperienceNeeded.Valid {
			gitlabIssue.ExperienceNeeded = nulls.String{String: "moderate", Valid: true}
		}

		// Allocate empty issue if there is no issue with the same id
		dbIssue := models.Issue{}
		if err := models.DB.Where("issue_id = ?", node.ID).First(&dbIssue); err != nil {
			verrs, err := gitlabIssue.Validate(models.DB)
			if verrs.HasAny() {
				fmt.Println(verrs.Error())
				continue
			}
			if err != nil {
				fmt.Println(errors.WithMessage(err, "Issues isn't valid"))
				continue
			}
			issuesToCreate = append(issuesToCreate, *gitlabIssue)
			continue
		}

		// if issue already exists
		gitlabIssue.ID = dbIssue.ID
		gitlabIssue.UpdatedAt = time.Now()
		verrs, err := gitlabIssue.Validate(models.DB)
		if verrs.HasAny() {
			fmt.Println(verrs.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.WithMessage(err, "Issues isn't valid"))
			continue
		}
		issuesToUpdate = append(issuesToUpdate, *gitlabIssue)
	}
	// Create all the new issues
	err := models.DB.Create(&issuesToCreate)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to create issues"))
	}
	// Update all existing issues
	err = models.DB.Update(&issuesToUpdate)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to update issues"))
	}

	w.updateProjectOnFinish(repository)
}


// Update project info when issues have been updated
func (w *Worker) updateProjectOnFinish(repository *models.Repository) {
	go w.searchForDanglingIssues(repository)
	var err error

	repository.IssueCount, err = models.DB.Where("closed=false and repository_id=?", repository.ID).Count(&models.Issue{})
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to count"))
	}

	repository.LastParsed = time.Now()
	verr, err := models.DB.ValidateAndUpdate(repository)
	if verr.HasAny() {
		fmt.Println(verr.Error())
	}
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed update last parsed repo"))
	}

	repos := &models.Repositories{}
	if err = models.DB.Where("project_id=?", repository.ProjectID).All(repos); err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to find repos"))
	}

	count := 0
	for _, repo := range *repos {
		count += repo.IssueCount
	}

	project := &models.Project{}
	if err = models.DB.Find(project, repository.ProjectID); err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to find project"))
	}

	project.IssuesCount = count
	verr, err = models.DB.ValidateAndUpdate(project)
	if verr.HasAny() {
		fmt.Println(verr.Error())
	}
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to update project"))
	}

	project.IsGitHub = false;
}

// Cleanup project issues that have been closed or couldn't be found in the repo
func (w *Worker) searchForDanglingIssues(repository *models.Repository) {
	issues := models.Issues{}
	err := models.DB.Where("updated_at < current_timestamp - interval '6 minutes' and closed = false and project_id = ?", repository.ProjectID).All(&issues)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to find unclosed issues"))
		return
	}
	issuesToClose := models.Issues{}
	for _, issue := range issues {
		issueData, _, err := client.Issues.GetIssue(repository.ProjectID, issue.Number, nil)

		// This might close an issue if there is a network error
		// but it's better to close an issue and reopen it later rather than leaving dangling issues
		if err != nil || issueData.State == "closed" {
			issue.Closed = true
			issuesToClose = append(issuesToClose, issue)
			continue
		}
	}

	verr, err := models.DB.ValidateAndUpdate(&issuesToClose)
	if verr.HasAny() {
		fmt.Println(verr.Error())
	}
	if err != nil {
		fmt.Println(errors.Wrap(err, "couldn't update issue"))
	}
}

var WorkerInst = Worker{}
