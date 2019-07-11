package worker_gitlab

import (
	"strings"

	"github.com/ossn/fixme_backend/models"
)

func split(r rune) bool {
	return r == ' ' || r == ':' || r == '.' || r == ','
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
