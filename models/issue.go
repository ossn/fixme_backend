package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Issue struct {
	ID uuid.UUID `json:"id" db:"id"`
	IssueID int `json:"issue_id" db:"issue_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title string `json:"title" db:"title"`
	ExperienceNeeded string `json:"experience_needed" db:"experience_needed"`
	URL string `json:"url" db:"url"`
	Body string `json:"body" db:"body"`
	Type string `json:"type" db:"type"`
	Project Project `json:"project" db:"-" belongs_to:"project"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Number int `json:"number" db:"number"`
	Closed bool `json:"closed" db:"closed"`
	Technologies slices.String `json:"technologies" db:"technologies"`
	Labels slices.String `json:"labels" db:"labels"`
	IsGitHub bool `json:"is_github" db:"is_github"`
}

type Issues[] Issue

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func(i * Issue) Validate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.Validate(
		& validators.IntIsPresent {Field: i.IssueID, Name: "IssueID"},
		& validators.IntIsPresent {Field: i.Number,	Name: "Number"},
		& validators.StringIsPresent {Field: i.URL,	Name: "URL"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func(i * Issue) ValidateCreate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func(i * Issue) ValidateUpdate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.NewErrors(), nil
}
