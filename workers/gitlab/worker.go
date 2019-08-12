package worker_gitlab

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
  "encoding/base64"

	"github.com/ossn/fixme_backend/models"
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
	create_technologies_map()
	for {
		go w.UpdateRepositoryTopics()
		time.Sleep(1 * time.Hour)
	}
}

// Get all the tags repositories and set them to the project
func (w *Worker) UpdateRepositoryTopics() {
	w.waitUntilLimitIsRefreshed()
	projects := models.Projects{}
	err := models.DB.Where("is_github = ?", false).All(&projects)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get projects"))
		return
	}
	for i, project := range projects {

		var filteredTechnologies []string

		projectInfo, _, err := client.Projects.GetProject(project.ProjectID, nil)

		if err != nil {
			fmt.Println(errors.WithMessage(err, "couldn't find languages"))
			return
		} else {
			description := strings.FieldsFunc(projectInfo.Description, split)
			filteredTechnologies = append(filteredTechnologies, searchForMatchingTechnologies(description)...)
		}

		readmeOptions := &gitlab.GetFileOptions{
			Ref: gitlab.String("master"),
		}
		response, _, err := client.RepositoryFiles.GetFile(project.ProjectID, "README.md", readmeOptions)
		if err != nil {
			fmt.Println(err)
		}	else {
			data, err := base64.StdEncoding.DecodeString(response.Content)
			if err != nil {
							fmt.Println("error:", err)
			} else {
				readme := strings.FieldsFunc(string(data), split)
				filteredTechnologies = append(filteredTechnologies, searchForMatchingTechnologies(readme)...)
			}
		}


		//Find languages
		languages, _, err := client.Projects.GetProjectLanguages(project.ProjectID, nil)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "couldn't find languages"))
			return
		}

		var projectLanguages []string
		for language, _ := range *languages {
				projectLanguages = append(projectLanguages, language)
		}

		technologies := append(cleanupArray(filteredTechnologies), projectLanguages...)
		project.Technologies = technologies

		project.IsGitHub = false


		verr, err := project.Validate(models.DB)
		if verr.HasAny() {
			fmt.Println(verr.Error())
			continue
		}
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't save repos from gitlab"))
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
			time.Sleep(time.Second*10)	// we don't know the limits
}

// Get first issues
func (w *Worker) getInitialIssues() {
	w.waitUntilLimitIsRefreshed()
	lastUpdatedProject := models.Project{}
	err := models.DB.Where("is_github = ?", false).Order("last_parsed asc").First(&lastUpdatedProject)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to get issues"))
		return
	}


	// Find issues
	openState := "opened"
	issuesOptions := &gitlab.ListProjectIssuesOptions{
		State: &openState,
	}

	issueData, response, err := client.Issues.ListProjectIssues(lastUpdatedProject.ProjectID, issuesOptions)

	if err != nil {
		fmt.Println(errors.WithMessage(err, "couldn't load initial issues"))
		return
	}


	// full pages (20) of issues
	for len(issueData) == 20 {
		go w.parseAndSaveIssues(issueData, &lastUpdatedProject, true)
		issuesOptions.ListOptions.Page = response.NextPage
		issueData, response, err = client.Issues.ListProjectIssues(lastUpdatedProject.ProjectID, issuesOptions)
	}

	// last page of issues
	if len(issueData) > 0 {
		go w.parseAndSaveIssues(issueData, &lastUpdatedProject, false)
	}
}


// Parse and save gitlab issues
func (w *Worker) parseAndSaveIssues(issueData []*gitlab.Issue, project *models.Project, hasPreviousPage bool) {
	issuesToCreate := models.Issues{}
	issuesToUpdate := models.Issues{}
	for _, node := range issueData {
		gitlabIssue := &models.Issue{
			IssueID:      node.ID,
			Body:         node.Description,
			Title:        node.Title,
			Closed:       node.State == "closed",
			Number:       node.IID,
			URL:          node.WebURL,
			ProjectID:    project.ID,
		}

		// Parse gitlab labels
		var labels []string
		for _, label := range node.Labels {
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


		gitlabIssue.IsGitHub = false
		gitlabIssue.Technologies = project.Technologies
		gitlabIssue.Labels = labels
		gitlabIssue.ExperienceNeeded = difficulty

		// Allocate empty issue if there is no issue with the same id
		dbIssue := models.Issue{}
		if err := models.DB.Where("is_github = ? and issue_id = ?", "false", node.ID).First(&dbIssue); err != nil {
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

// Cleanup project issues that have been closed or couldn't be found in the repo
func (w *Worker) searchForDanglingIssues(project *models.Project) {
	issues := models.Issues{}
	err := models.DB.Where("updated_at < current_timestamp - interval '6 minutes' and closed = false and project_id = ?", project.ID).All(&issues)
	if err != nil {
		fmt.Println(errors.WithMessage(err, "Failed to find unclosed issues"))
		return
	}
	issuesToClose := models.Issues{}
	for _, issue := range issues {
		issueData, _, err := client.Issues.GetIssue(project.ProjectID, issue.Number, nil)

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
