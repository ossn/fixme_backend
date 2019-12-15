package worker2

import (
	//"fmt"
	"strings"

	"github.com/gobuffalo/pop/nulls"
	"github.com/ossn/fixme_backend/models"
	//"github.com/pkg/errors"
)

func split(r rune) bool {
	return r == ' ' || r == ':' || r == '.' || r == ','
}

// Remove empty and duplicate strings from an array
func cleanupArray(s []string) (r []string) {
	seen := make(map[string]bool, len(s))
	seen[""] = true
	for _, str := range s {
		if _, exists := seen[str]; !exists {
			seen[str] = true
			r = append(r, str)
		}
	}
	return
}

// Searches if a label matches some known labels and updates the model
func searchForMatchingLabels(label *string, model *models.Issue) bool {
	switch strings.ToLower(*label) {
	case "help_wanted", "help wanted", "good first issue", "easyfix", "easy":
		model.ExperienceNeeded = nulls.String{String: "easy", Valid: true}
		return true
	case "moderate":
		model.ExperienceNeeded = nulls.String{String: "moderate", Valid: true}
		return true
	case "senior":
		model.ExperienceNeeded = nulls.String{String: "senior", Valid: true}
		return true
	case "enhancement":
		model.Type = nulls.String{String: "enhancement", Valid: true}
		return true
	case "bug", "bugfix":
		model.Type = nulls.String{String: "bugfix", Valid: true}
		return true
	}
	return false
}
