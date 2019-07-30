package worker_github

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
  "encoding/json"

	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	Worker struct {
		ctx context.Context
	}
)

var (
	client *githubv4.Client
)

func (w *Worker) Init(ctx context.Context, c <-chan os.Signal) {
	w.ctx = ctx
	token := os.Getenv("GITHUB_TOKEN")
	var src oauth2.TokenSource
	if len(token) < 1 {
		panic("Please provide a github token")
	} else {
		src = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
	}
	httpClient := oauth2.NewClient(ctx, src)

	client = githubv4.NewClient(httpClient)
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

func (w *Worker) checkRateLimitStatus() (bool, time.Time, error) {
	rateLimitQuery := rateLimitQuery{}
	err := client.Query(w.ctx, &rateLimitQuery, nil)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't check the rate limit usage"))
		return true, time.Time{}, errors.WithMessage(err, "couldn't check the rate limit usage")
	}
	rateLimitData := rateLimitQuery.RateLimit
	resetAt, err := time.Parse(time.RFC3339, rateLimitData.ResetAt)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't check the rate limit usage"))
		return true, time.Time{}, errors.WithMessage(err, "couldn't check the rate limit usage")
	}
	if rateLimitData.Remaining < 100 {
		return true, resetAt, nil
	}
	return false, time.Time{}, nil
}

// Func to start repo topics polling
func (w *Worker) repositoryTopicsPolling() {
	create_technologies_map()
	for {
		go w.UpdateRepositoryTopics()
		time.Sleep(1 * time.Hour)
	}
}

// Get all the topics repositories and set them to the project
func (w *Worker) UpdateRepositoryTopics() {
	w.waitUntilLimitIsRefreshed()
	projects := models.Projects{}
	err := models.DB.Where("is_github = ?", true).All(&projects)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}
	for i, project := range projects {
		name, owner, err := getNameAndOwner(project.Link)

		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repos from gitlab"))
			continue
		}

		projectInfo := repoQuery{}
		err = client.Query(w.ctx, &projectInfo, map[string]interface{}{"name": name, "owner": owner})
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repo's info from github"))
			continue
		}

		description := strings.FieldsFunc(projectInfo.Repository.Description, split)
		filteredTechnologies := searchForMatchingTechnologies(description)

		readme := strings.FieldsFunc(projectInfo.Repository.Object.Blob.Text, split)
		filteredTechnologies = append(filteredTechnologies, searchForMatchingTechnologies(readme)...)

		topics := []string{}
		for _, topic := range projectInfo.Repository.RepositoryTopics.Nodes {
			topics = append(topics, topic.Topic.Name)
		}
		filteredTechnologies = append(filteredTechnologies, searchForMatchingTechnologies(topics)...)

		project.IsGitHub = true
		project.Tags = cleanupArray(filteredTechnologies)


		projectLanguages := make(map[string]float64)

		for _, language := range projectInfo.Repository.Languages.Edges {
			percentage := (float64(language.Size)/float64(projectInfo.Repository.Languages.TotalSize))*100
			projectLanguages[language.Node.Name] = percentage
		}

		bytes, err := json.Marshal(projectLanguages)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "failed to convert to bytes"))
			continue
		}

		project.Languages = string(bytes)


		verr, err := project.Validate(models.DB)
		if verr.HasAny() {
			fmt.Println(verr.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't save repos from github"))
			continue
		}
		projects[i] = project
	}

	verr, err := models.DB.ValidateAndUpdate(&projects)
		if err != nil || verr.HasAny() {
		fmt.Println(err, verr.Error())
	}
}

// waitUntilLimitIsRefreshed: A function that waits until the next github query can be executed
func (w *Worker) waitUntilLimitIsRefreshed () {
		limitExceeded, resetAt, err := w.checkRateLimitStatus()
		if err != nil {
			// if there is an issue retry in 5 minutes
			time.Sleep(time.Minute*5)
			w.waitUntilLimitIsRefreshed()
		}
		if limitExceeded {
			time.Sleep(time.Until(resetAt))
			w.waitUntilLimitIsRefreshed()
		}
}

// Get first issues
func (w *Worker) getInitialIssues() {
	w.waitUntilLimitIsRefreshed()
	lastUpdatedProject := models.Project{}
	err := models.DB.Where("is_github = ?", true).Order("last_parsed asc").First(&lastUpdatedProject)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to get issues"))
		return
	}

	name, owner, err := getNameAndOwner(lastUpdatedProject.Link)
	if err != nil {
		return
	}
	variables := map[string]interface{}{"name": name, "owner": owner}
	issueData := initialIssueQuery{}
	err = client.Query(w.ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't load initial issues"))
		return
	}

	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.parseAndSaveIssues(issueQueryWithBefore(issueData), &lastUpdatedProject, lastUpdatedProject.Languages, hasPreviousPage)

	if hasPreviousPage {
		w.getExtraIssues(&name, &owner, &issueData.Repository.Issues.PageInfo.StartCursor, &lastUpdatedProject, lastUpdatedProject.Languages)
	}

}

// Get next page of issues
func (w *Worker) getExtraIssues(name, owner *githubv4.String, before *string, project *models.Project, languages string) {
w.waitUntilLimitIsRefreshed()
	variables := map[string]interface{}{"name": *name, "owner": *owner, "before": githubv4.String(*before)}
	issueData := issueQueryWithBefore{}
	err := client.Query(w.ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to get additional issues"))
		return
	}

	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.parseAndSaveIssues(issueData, project, languages, hasPreviousPage)

	if hasPreviousPage {
		w.getExtraIssues(name, owner, &issueData.Repository.Issues.PageInfo.StartCursor, project, languages)
	}

}

// Parse and save github issues
func (w *Worker) parseAndSaveIssues(issueData issueQueryWithBefore, project *models.Project, languages string, hasPreviousPage bool) {
	issuesToCreate := models.Issues{}
	issuesToUpdate := models.Issues{}
	for _, node := range issueData.Repository.Issues.Nodes {
		githubIssue := &models.Issue{
			IssueID:	    node.DatabaseID,
			Body:         node.Body,
			Title:        node.Title,
			Closed:       node.Closed,
			Number:       node.Number,
			URL:          node.URL,
			ProjectID:    project.ID,
		}


		// Parse gitlab labels
		labels := []string{}
		for _, node := range node.Labels.Nodes {
			label := node.Name
			labels = append(labels, label)
		}

		difficulty := searchForMatchingLabels(labels)

		if difficulty == "unknown" {
			for _, label := range labels {
				// Split label based on known delimeters
				parts := strings.FieldsFunc(label, split)
				// If label hasn't been matched try again with the splited string
				if len(parts) > 1 {
						difficulty = searchForMatchingLabels(parts)
				}
				if difficulty == "easy" {
					break
				}
			}
		}

		githubIssue.IsGitHub = true
		githubIssue.Labels = labels
		githubIssue.ExperienceNeeded = difficulty

		// Allocate empty issue
		dbIssue := models.Issue{}
		if err := models.DB.Where("is_github = ? and issue_id = ?", "true", node.DatabaseID).First(&dbIssue); err != nil {
			verrs, err := githubIssue.Validate(models.DB)
			if verrs.HasAny() {
				fmt.Println(verrs.Error())
				continue
			}
			if err != nil {
				fmt.Println(errors.WithMessage(err, "Issues isn't valid"))
				continue
			}
			issuesToCreate = append(issuesToCreate, *githubIssue)
			continue
		}

		githubIssue.ID = dbIssue.ID
		githubIssue.UpdatedAt = time.Now()
		verrs, err := githubIssue.Validate(models.DB)
		if verrs.HasAny() {
			fmt.Println(verrs.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.WithMessage(err, "Issues isn't valid"))
			continue
		}
		issuesToUpdate = append(issuesToUpdate, *githubIssue)
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

	// Update repo record once all the github issues have been parsed
	if !hasPreviousPage {
		w.updateProjectOnFinish(project)
	}
}

// Update project info when issues have been updated
func (w *Worker) updateProjectOnFinish(project *models.Project) {
	go w.searchForDanglingIssues(project)
	var err error

	project.IssuesCount, err = models.DB.Where("closed = ? and project_id = ?", "false", project.ID).Count(&models.Issue{})
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to count"))
	}

	project.LastParsed = time.Now()
	verr, err := models.DB.ValidateAndUpdate(project)
	if verr.HasAny() {
		fmt.Println(verr.Error())
	}
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed update last parsed project"))
	}
}

// Cleanup project issues that have been deleted or couldn't be found in the repo
func (w *Worker) searchForDanglingIssues(project *models.Project) {
	issues := models.Issues{}
	name, owner, err := getNameAndOwner(project.Link)
	if err != nil {
		return
	}
	err = models.DB.Where("updated_at < current_timestamp - interval '6 minutes' and closed = false and project_id = ?", project.ID).All(&issues)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to find unclosed issues"))
		return
	}
	issuesToClose := models.Issues{}
	for _, issue := range issues {
		issueStatus := issueStatusQuery{}
		requestParams := map[string]interface{}{
			"name": name, "owner": owner, "number": githubv4.Int(issue.Number)}
		err = client.Query(w.ctx, &issueStatus, requestParams)
		// This might close an issue if there is a network error
		// but it's better to close an issue and reopen it later rather than leaving dangling issues
		if err != nil || issueStatus.Repository.Issue.Closed {
			fmt.Println("couldn't load issue from github "+string(owner)+" "+string(name))
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
