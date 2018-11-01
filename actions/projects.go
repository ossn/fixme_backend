package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/ossn/fixme_backend/models"
	"github.com/ossn/fixme_backend/worker"
	"github.com/pkg/errors"
)

// ProjectsResource is the resource for the Project model
type ProjectsResource struct {
	buffalo.Resource
}

// List gets all Projects. This function is mapped to the path
// GET /projects
func (v ProjectsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	projects := &models.Projects{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Projects from the DB
	if err := q.All(projects); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(projects))
}

// Show gets the data for one Project. This function is mapped to
// the path GET /projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Project
	project := &models.Project{}

	// To find the Project the parameter project_id is used.
	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(project))
}

// New renders the form for creating a new Project.
// This function is mapped to the path GET /projects/new
func (v ProjectsResource) New(c buffalo.Context) error {
	return c.Render(200, r.JSON(&models.Project{}))
}

// Create adds a Project to the DB. This function is mapped to the
// path POST /projects
func (v ProjectsResource) Create(c buffalo.Context) error {
	// Allocate an empty Project
	project := &models.Project{}

	// Bind project to the html form elements
	if err := c.Bind(project); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(project)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(project))
	}

	// Force worker to update the topics
	go worker.WorkerInst.UpdateRepositoryTopics()

	repo := models.Repository{RepositoryUrl: project.Link, ProjectID: project.ID}
	count, err := tx.Where("repository_url=?", project.Link).Count(&repo)
	if err != nil {
		return errors.WithStack(err)
	}

	if count < 1 {
		verrs, err = tx.ValidateAndCreate(&repo)

		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the response
			c.Set("errors", verrs)

			return c.Render(422, r.JSON(project))
		}
	}

	return c.Render(201, r.JSON(project))
}

// Edit renders a edit form for a Project. This function is
// mapped to the path GET /projects/{project_id}/edit
func (v ProjectsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Project
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(project))
}

// Update changes a Project in the DB. This function is mapped to
// the path PUT /projects/{project_id}
func (v ProjectsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Project
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(404, err)
	}

	oldProjectUrl := project.Link

	// Bind Project to the html form elements
	if err := c.Bind(project); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(project)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(project))
	}

	// Force worker to update topic list
	go worker.WorkerInst.UpdateRepositoryTopics()

	if oldProjectUrl != project.Link {
		repo := models.Repository{}
		if err = tx.Where("project_id=?", project.ID).Where("repository_url=?", oldProjectUrl).First(&repo); err != nil {
			// if no repo is found just skip update
			return errors.Wrap(err, "Failed to find repo")
		}
		repo.RepositoryUrl = project.Link
		repo.LastParsed = time.Unix(0, 0)
		verrs, err = tx.ValidateAndUpdate(&repo)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the response
			c.Set("errors", verrs)

			return c.Render(422, r.JSON(project))
		}
	}
	return c.Render(200, r.JSON(project))
}

// Destroy deletes a Project from the DB. This function is mapped
// to the path DELETE /projects/{project_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Project
	project := &models.Project{}

	// To find the Project the parameter project_id is used.
	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(project); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(project))
}
