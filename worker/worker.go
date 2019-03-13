package worker

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
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	Worker struct {
		ctx context.Context
	}

	/**
	* GraphQL Types
	 */

	PageInfo struct {
		StartCursor     string
		HasPreviousPage bool
	}
	Issues struct {
		Nodes []struct {
			Title      string
			Body       string
			Closed     bool
			Number     int
			URL        string
			CreatedAt  string
			DatabaseID int
			Labels     struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"labels(first:100)"`
		}
		PageInfo PageInfo
	}

	language struct {
		Repository struct {
			PrimaryLanguage struct {
				Name string
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	initialIssueQuery struct {
		Repository struct {
			Issues Issues `graphql:"issues(last: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	issueQueryWithBefore struct {
		Repository struct {
			Issues Issues `graphql:"issues(last: 100, before: $before)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	tagsQuery struct {
		Repository struct {
			RepositoryTopics struct {
				Nodes []struct {
					Topic struct {
						Name string
					}
				}
			} `graphql:"repositoryTopics(first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	issueStatusQuery struct {
		Repository struct {
			Issue struct {
				Closed bool
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	rateLimitQuery struct {
		RateLimit struct {
			Remaining int    `graphql:"remaining"`
			ResetAt   string `graphql:"resetAt"`
		} `graphql:"rateLimit"`
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
	for {
		go w.UpdateRepositoryTopics()
		time.Sleep(1 * time.Hour)
	}
}

// Get all the tags repositories and set them to the project
func (w *Worker) UpdateRepositoryTopics() {
	w.waitUntilLimitIsRefreshed()
	repos := models.Repositories{}
	err := models.DB.All(&repos)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}
	repoIndexMap := make(map[uuid.UUID][]int, len(repos))
	for i, repo := range repos {
		repoIndexMap[repo.ProjectID] = append(repoIndexMap[repo.ProjectID], i)
		name, owner, err := getNameAndOwner(repo.RepositoryUrl)
		if err != nil {
			continue
		}
		tags := tagsQuery{}
		err = client.Query(w.ctx, &tags, map[string]interface{}{"name": name, "owner": owner})
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repos from github"))
			continue
		}

		repoTags := []string{}
		for _, tag := range tags.Repository.RepositoryTopics.Nodes {
			repoTags = append(repoTags, tag.Topic.Name)
		}

		repo.Tags = cleanupArray(repoTags)
		verr, err := repo.Validate(models.DB)
		if verr.HasAny() {
			fmt.Println(verr.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't save repos from github"))
			continue
		}
		repos[i] = repo
	}
	verr, err := models.DB.ValidateAndUpdate(&repos)
		if err != nil || verr.HasAny() {
		fmt.Println(err, verr.Error())
	}

	projects := models.Projects{}
	err = models.DB.All(&projects)
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
	lastUpdatedRepo := models.Repository{}
	err := models.DB.Order("last_parsed asc").First(&lastUpdatedRepo)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to get issues"))
		return
	}

	name, owner, err := getNameAndOwner(lastUpdatedRepo.RepositoryUrl)
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

	languageRequest := language{}
	err = client.Query(w.ctx, &languageRequest, variables)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't find language"))
		return
	}
	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.parseAndSaveIssues(issueQueryWithBefore(issueData), &lastUpdatedRepo, &languageRequest.Repository.PrimaryLanguage.Name, hasPreviousPage)

	if hasPreviousPage {
		w.getExtraIssues(&name, &owner, &issueData.Repository.Issues.PageInfo.StartCursor, &lastUpdatedRepo, &languageRequest.Repository.PrimaryLanguage.Name)

	}

}

// Get next page of issues
func (w *Worker) getExtraIssues(name, owner *githubv4.String, before *string, repository *models.Repository, language *string) {
w.waitUntilLimitIsRefreshed()
	variables := map[string]interface{}{"name": *name, "owner": *owner, "before": githubv4.String(*before)}
	issueData := issueQueryWithBefore{}
	err := client.Query(w.ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to get additional issues"))
		return
	}

	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.parseAndSaveIssues(issueData, repository, language, hasPreviousPage)

	if hasPreviousPage {
		w.getExtraIssues(name, owner, &issueData.Repository.Issues.PageInfo.StartCursor, repository, language)
	}

}

// Parse and save github issues
func (w *Worker) parseAndSaveIssues(issueData issueQueryWithBefore, repository *models.Repository, language *string, hasPreviousPage bool) {
	issuesToCreate := models.Issues{}
	issuesToUpdate := models.Issues{}
	for _, node := range issueData.Repository.Issues.Nodes {
		githubIssue := &models.Issue{
			GithubID:     node.DatabaseID,
			Body:         nulls.String{String: node.Body, Valid: node.Body != ""},
			Title:        nulls.String{String: node.Title, Valid: node.Title != ""},
			Closed:       node.Closed,
			Number:       node.Number,
			URL:          node.URL,
			RepositoryID: repository.ID,
			ProjectID:    repository.ProjectID,
			Language:     nulls.String{String: strings.ToLower(*language), Valid: *language != ""},
		}

		// Parse github labels
		labels := []string{}
		for _, label := range node.Labels.Nodes {
			name := &label.Name
			labels = append(labels, *name)
			// Search for known labels
			matched := searchForMatchingLabels(name, githubIssue)
			// Split name based on known delimeters
			tmp := strings.FieldsFunc(*name, split)
			// If label hasn't been matched try again with the splited string
			if !matched && len(tmp) > 1 {
				for _, label := range tmp {
					searchForMatchingLabels(&label, githubIssue)
				}
			}
		}

		githubIssue.Labels = labels
		// Initialize experience needed with moderate
		if !githubIssue.ExperienceNeeded.Valid {
			githubIssue.ExperienceNeeded = nulls.String{String: "moderate", Valid: true}
		}

		// Allocate empty issue
		dbIssue := models.Issue{}
		if err := models.DB.Where("github_id = ?", node.DatabaseID).First(&dbIssue); err != nil {
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
		w.updateProjectOnFinish(repository)
	}
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
}

// Cleanup project issues that have been deleted or couldn't be found in the repo
func (w *Worker) searchForDanglingIssues(repository *models.Repository) {
	issues := models.Issues{}
	name, owner, err := getNameAndOwner(repository.RepositoryUrl)
	if err != nil {
		return
	}
	err = models.DB.Where("updated_at < current_timestamp - interval '6 minutes' and closed = false and project_id = ?", repository.ProjectID).All(&issues)
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
		if err != nil {
			fmt.Println(errors.WithMessage(err, "couldn't load issue from github "+string(owner)+" "+string(name)))
			issue.Closed = true
			issuesToClose = append(issuesToClose, issue)
			continue
		}

		if issueStatus.Repository.Issue.Closed {
			issue.Closed = true
			issuesToClose = append(issuesToClose, issue)
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
