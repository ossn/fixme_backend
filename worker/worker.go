package worker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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
)

var (
	client *githubv4.Client
)

func (w *Worker) Init(ctx context.Context, c <-chan os.Signal) {
	(*w).ctx = ctx
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
	go w.startPolling(c)
}

func (w *Worker) startPolling(c <-chan os.Signal) {
	go func() {
		<-c
		os.Exit(1)

	}()
	go func() {
		for {
			go w.UpdateRepositoryTopics()
			time.Sleep(1 * time.Hour)
		}
	}()
	for {

		//FIXME: Check github limits and run based on those
		go w.getIssues()
		time.Sleep(5 * time.Minute)
	}
}

func (w *Worker) UpdateRepositoryTopics() {
	repos := models.Repositories{}
	err := models.DB.All(&repos)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}
	for _, repo := range repos {
		name, owner, err := getNameAndOwner(repo.RepositoryUrl)
		if err != nil {
			continue
		}
		tags := tagsQuery{}
		err = client.Query((*w).ctx, &tags, map[string]interface{}{"name": name, "owner": owner})
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repos from github"))
			continue
		}

		repoTags := []string{}
		for _, tag := range tags.Repository.RepositoryTopics.Nodes {
			repoTags = append(repoTags, tag.Topic.Name)
		}

		repo.Tags = repoTags
		err = models.DB.Update(&repo)
		if err != nil {
			fmt.Println(errors.Wrap(err, "couldn't load repos from github"))
			continue
		}
	}

	projects := models.Projects{}
	err = models.DB.All(&projects)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get repos"))
		return
	}
	for _, project := range projects {
		repos = models.Repositories{}
		err = models.DB.Where("project_id = ?", project.ID).All(&repos)
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to find repos"))
			continue
		}

		tags := []string{}
		for _, repo := range repos {
			tags = append(tags, repo.Tags...)
		}
		project.Tags = tags
		err = models.DB.Update(&project)
		if err != nil {
			fmt.Println(errors.Wrap(err, "failed to save project"))
			continue
		}
	}
}

func (w *Worker) getIssues() {
	lastUpdatedRepo := models.Repository{}
	err := models.DB.Order("last_parsed asc").First(&lastUpdatedRepo)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to get issues"))
		return

	}
	name, owner, err := getNameAndOwner(lastUpdatedRepo.RepositoryUrl)
	if err != nil {
		return
	}
	variables := map[string]interface{}{"name": name, "owner": owner}
	issueData := initialIssueQuery{}
	err = client.Query((*w).ctx, &issueData, variables)
	if err != nil {
		fmt.Println(errors.Wrap(err, "couldn't load initial issues"))
		return
	}

	languageRequest := language{}
	err = client.Query((*w).ctx, &languageRequest, variables)
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
	err := client.Query((*w).ctx, &issueData, variables)
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

func getNameAndOwner(url string) (githubv4.String, githubv4.String, error) {

	tmp := strings.Split(strings.TrimSuffix(url, "/"), "/")
	if len(tmp) < 2 {
		err := errors.New(fmt.Sprintf("Couldn't find repo %s", url))
		fmt.Println(errors.Wrap(err, "failed to find url"))
		return githubv4.String(""), githubv4.String(""), err
	}
	return githubv4.String(tmp[len(tmp)-1]), githubv4.String(tmp[len(tmp)-2]), nil
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

var WorkerInst = Worker{}
