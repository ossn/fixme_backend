package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Repository struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	RepositoryUrl string        `json:"repository_url" db:"repository_url"`
	Project       Project       `json:"project" db:"-" belongs_to:"project"`
	ProjectID     uuid.UUID     `json:"project_id" db:"project_id"`
	IssueCount    int           `json:"issue_count" db:"issue_count"`
	Issues        Issues        `json:"issues" db:"-" has_many:"issues"`
	LastParsed    time.Time     `json:"-" db:"last_parsed"`
	Tags          slices.String `json:"tags" db:"tags"`
}

type Repositories []Repository

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *Repository) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.RepositoryUrl, Name: "RepositoryUrl"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (r *Repository) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (r *Repository) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
