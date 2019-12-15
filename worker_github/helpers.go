package worker_github

import (
	"fmt"
	"strings"

	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

func split(r rune) bool {
	return r == ' ' || r == ':' || r == '.' || r == ','
}

// Extracts the name and the owner from a git url
func getNameAndOwner(url string) (githubv4.String, githubv4.String, error) {

	tmp := strings.Split(strings.TrimSuffix(url, "/"), "/")
	if len(tmp) < 2 {
		err := errors.New(fmt.Sprintf("Couldn't find repo %s", url))
		fmt.Println(errors.Wrap(err, "failed to find url"))
		return githubv4.String(""), githubv4.String(""), err
	}
	return githubv4.String(tmp[len(tmp)-1]), githubv4.String(tmp[len(tmp)-2]), nil
}

// Searches if a label matches some known labels and updates the model
func searchForMatchingLabels(label *string, model *models.Issue) bool {
	switch strings.ToLower(*label) {
	case "help_wanted", "help wanted", "good first issue", "easyfix", "easy":
		model.ExperienceNeeded = "easy"
		return true
	case "moderate":
		model.ExperienceNeeded = "moderate"
		return true
	case "senior":
		model.ExperienceNeeded = "senior"
		return true
	case "enhancement":
		model.Type = "enhancement"
		return true
	case "bug", "bugfix":
		model.Type = "bugfix"
		return true
	}
	return false
}
