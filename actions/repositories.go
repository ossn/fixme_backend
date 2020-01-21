package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"
)

// RepositoriesResource is the resource for the Repository model
type RepositoriesResource struct {
	buffalo.Resource
}

// List gets all Repositories. This function is mapped to the path
// GET /repositories
func (v RepositoriesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	repositories := &models.Repositories{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Repositories from the DB
	if err := q.Eager("Project").All(repositories); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(repositories))
}

// Show gets the data for one Repository. This function is mapped to
// the path GET /repositories/{repository_id}
func (v RepositoriesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Repository
	repository := &models.Repository{}

	// To find the Repository the parameter repository_id is used.
	if err := tx.Eager("Project").Find(repository, c.Param("repository_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(repository))
}

// New renders the form for creating a new Repository.
// This function is mapped to the path GET /repositories/new
func (v RepositoriesResource) New(c buffalo.Context) error {
	return c.Render(200, r.JSON(&models.Repository{}))
}

// Create adds a Repository to the DB. This function is mapped to the
// path POST /repositories
func (v RepositoriesResource) Create(c buffalo.Context) error {
	// Allocate an empty Repository
	repository := &models.Repository{}

	// Bind repository to the html form elements
	if err := c.Bind(repository); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(repository)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(repository))
	}

	return c.Render(201, r.JSON(repository))
}

// Edit renders a edit form for a Repository. This function is
// mapped to the path GET /repositories/{repository_id}/edit
func (v RepositoriesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Repository
	repository := &models.Repository{}

	if err := tx.Find(repository, c.Param("repository_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(repository))
}

// Update changes a Repository in the DB. This function is mapped to
// the path PUT /repositories/{repository_id}
func (v RepositoriesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Repository
	repository := &models.Repository{}

	if err := tx.Find(repository, c.Param("repository_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Repository to the html form elements
	if err := c.Bind(repository); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(repository)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(repository))
	}

	return c.Render(200, r.JSON(repository))
}

// Destroy deletes a Repository from the DB. This function is mapped
// to the path DELETE /repositories/{repository_id}
func (v RepositoriesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Repository
	repository := &models.Repository{}

	// To find the Repository the parameter repository_id is used.
	if err := tx.Find(repository, c.Param("repository_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(repository); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(repository))
}
