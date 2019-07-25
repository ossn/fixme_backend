package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Project struct {
	ID uuid.UUID `json:"id" db:"id"`
	ProjectID int `json:"project_id" db:"project_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	LastParsed time.Time `json:"-" db:"last_parsed"`
	FirstColor string `json:"first_color" db:"first_color"`
	SecondColor string `json:"second_color" db:"second_color"`
	Logo string `json:"logo" db:"logo"`
	Link string `json:"link" db:"link"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Languages slices.String `json:"languages" db:"languages"`
	Tags slices.String `json:"tags" db:"tags"`
	IssuesCount int `json:"issues_count" db:"issues_count"`
	Issues Issues `json:"issues" db:"-" has_many:"issues"`
	IsGitHub bool `json:"is_github" db:"is_github"`
}

type Projects[] Project

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func(p * Project) Validate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent {Field: p.Name, Name: "Name"},
		&validators.StringIsPresent {Field: p.Description,	Name: "Description"},
		&validators.StringIsPresent {Field: p.Logo, Name: "Logo"},
		&validators.StringIsPresent {Field: p.Link,	Name: "Link"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func(p * Project) ValidateCreate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func(p * Project) ValidateUpdate(tx * pop.Connection)( * validate.Errors, error) {
	return validate.NewErrors(), nil
}
