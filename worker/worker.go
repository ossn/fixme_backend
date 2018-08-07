package worker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/nulls"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	Worker struct{}

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
)

var (
	client *githubv4.Client
	ctx    = context.Background()
)

func (w *Worker) Init() {

	token := os.Getenv("GITHUB_TOKEN")
	var src oauth2.TokenSource
	if len(token) < 1 {
		panic("Please provide a github token")
	} else {
		src = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
	}
	httpClient := oauth2.NewClient(context.Background(), src)

	client = githubv4.NewClient(httpClient)
	go w.StartPolling()
}

func (w *Worker) StartPolling() {
	for {

		//FIXME: Check github limits and run based on those
		go w.getIssues()
		time.Sleep(5 * time.Minute)
	}
}

func (w *Worker) getIssues() {
	lastUpdatedRepo := models.Repository{}
	err := models.DB.Order("last_parsed asc").First(&lastUpdatedRepo)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get issues"))
		return

	}

	tmp := strings.Split(strings.TrimSuffix(lastUpdatedRepo.RepositoryUrl, "/"), "/")
	name := githubv4.String(tmp[len(tmp)-1])
	owner := githubv4.String(tmp[len(tmp)-2])
	variables := map[string]interface{}{"name": name, "owner": owner}
	issueData := initialIssueQuery{}
	err = client.Query(ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.Wrap(err, "could load initial issues"))
		return
	}

	languageRequest := language{}
	err = client.Query(ctx, &languageRequest, variables)
	if err != nil {
		fmt.Println(errors.Wrap(err, "couldn't find language"))
		return
	}
	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.saveData(issueQueryWithBefore(issueData), &lastUpdatedRepo, &languageRequest.Repository.PrimaryLanguage.Name, hasPreviousPage)

	if hasPreviousPage {
		go w.getExtraIssues(&name, &owner, &issueData.Repository.Issues.PageInfo.StartCursor, &lastUpdatedRepo, &languageRequest.Repository.PrimaryLanguage.Name)
		return
	}

	updateProjectOnFinish(&lastUpdatedRepo)

}

func (w *Worker) getExtraIssues(name, owner *githubv4.String, before *string, repository *models.Repository, language *string) {
	variables := map[string]interface{}{"name": *name, "owner": *owner, "before": githubv4.String(*before)}
	issueData := issueQueryWithBefore{}
	err := client.Query(ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to get additional issues"))
		return
	}
	hasPreviousPage := issueData.Repository.Issues.PageInfo.HasPreviousPage
	go w.saveData(issueData, repository, language, hasPreviousPage)

	if hasPreviousPage {
		go w.getExtraIssues(name, owner, &issueData.Repository.Issues.PageInfo.StartCursor, repository, language)
		return
	}

}

func (w *Worker) saveData(issueData issueQueryWithBefore, repository *models.Repository, language *string, hasPreviousPage bool) {

	for _, node := range issueData.Repository.Issues.Nodes {
		model := &models.Issue{
			GithubID:     node.DatabaseID,
			Body:         nulls.String{String: node.Body, Valid: node.Body != ""},
			Title:        nulls.String{String: node.Title, Valid: node.Title != ""},
			Closed:       node.Closed,
			Number:       node.Number,
			URL:          node.URL,
			RepositoryID: (*repository).ID,
			ProjectID:    (*repository).ProjectID,
			Language:     nulls.String{String: strings.ToLower(*language), Valid: *language != ""},
		}

		labels := []string{}
		for _, label := range node.Labels.Nodes {
			name := &label.Name
			labels = append(labels, *name)
			matched := searchForMatchingLabels(name, model)
			tmp := strings.FieldsFunc(*name, split)
			if !matched && len(tmp) > 1 {
				for _, label := range tmp {
					searchForMatchingLabels(&label, model)
				}
			}
		}

		model.Labels = labels
		if !model.ExperienceNeeded.Valid {
			model.ExperienceNeeded = nulls.String{String: "moderate", Valid: true}
		}
		if err := models.DB.Where("github_id = ?", node.DatabaseID).First(&models.Issue{}); err != nil {
			err := models.DB.Create(model)
			if err != nil {
				fmt.Println(errors.Wrap(err, "failed to create model"))
			}
			continue
		}
		models.DB.Update(model)
	}
	if !hasPreviousPage {
		updateProjectOnFinish(repository)
	}
}

func split(r rune) bool {
	return r == ' ' || r == ':' || r == '.' || r == ','
}

func searchForMatchingLabels(label *string, model *models.Issue) bool {
	switch strings.ToLower(*label) {
	case "help_wanted", "help wanted", "good first issue", "easyfix", "easy":
		(*model).ExperienceNeeded = nulls.String{String: "easy", Valid: true}
		return true
	case "moderate":
		(*model).ExperienceNeeded = nulls.String{String: "moderate", Valid: true}
		return true
	case "senior":
		(*model).ExperienceNeeded = nulls.String{String: "senior", Valid: true}
		return true
	case "enhancement":
		(*model).Type = nulls.String{String: "enhancement", Valid: true}
		return true
	case "bug", "bugfix":
		(*model).Type = nulls.String{String: "bugfix", Valid: true}
		return true
	}
	return false
}

func updateProjectOnFinish(repository *models.Repository) {
	var err error

	(*repository).IssueCount, err = models.DB.Where("closed=false and repository_id=?", repository.ID).Count(&models.Issue{})
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to count"))
	}

	(*repository).LastParsed = time.Now()
	if err = models.DB.Update(repository); err != nil {
		fmt.Println(errors.Wrap(err, "failed to find last updated repo"))
	}

	repos := &models.Repositories{}
	if err = models.DB.Where("project_id=?", (*repository).ProjectID).All(repos); err != nil {
		fmt.Println(errors.Wrap(err, "Failed to find repos"))
	}

	count := 0
	for _, repo := range *repos {
		count += repo.IssueCount
	}

	project := &models.Project{}
	if err = models.DB.Find(project, (*repository).ProjectID); err != nil {
		fmt.Println(errors.Wrap(err, "Failed to find project"))
	}

	(*project).IssuesCount = count
	if err = models.DB.Update(project); err != nil {
		fmt.Println(errors.Wrap(err, "Failed to update project"))
	}

}
