package actions

import (
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"
)

// IssuesResource is the resource for the Issue model
type IssuesResource struct {
	buffalo.Resource
}

// List gets all Issues. This function is mapped to the path
// GET /issues
func (v IssuesResource) ListOpen(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	issues := &models.Issues{}
	params := c.Params()
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(params).Eager()

	whereClause := "closed = false"

	for _, filter := range []string{"language", "experience_needed", "type", "project_id"} {
		param := params.Get(filter)
		if param != "" {
			requestParamToQueryFilter(&whereClause, &param, &filter)
		}
	}
	// Retrieve all Issues from the DB
	if err := q.Where(whereClause).All(issues); err != nil {
		return errors.WithStack(err)
	}
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(issues))
}

// List gets all Issues. This function is mapped to the path
// GET issues without closed
func (v IssuesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	issues := &models.Issues{}
	params := c.Params()
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(params).Eager()

	// Retrieve all Issues from the DB
	if err := q.All(issues); err != nil {
		return errors.WithStack(err)
	}
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(issues))
}

// Show gets the data for one Issue. This function is mapped to
// the path GET /issues/{issue_id}
func (v IssuesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Issue
	issue := &models.Issue{}

	// To find the Issue the parameter issue_id is used.
	if err := tx.Find(issue, c.Param("issue_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(issue))
}

// List gets all Issues. This function is mapped to the path
// GET /issues-count
func (v IssuesResource) Count(c buffalo.Context) error {
	// Get the DB connection from the context
	q, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	issues := &models.Issues{}
	params := c.Params()
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".

	whereClause := "closed = false"
	for _, filter := range []string{"language", "experience_needed", "type", "project_id"} {
		param := params.Get(filter)
		if param != "" {
			requestParamToQueryFilter(&whereClause, &param, &filter)
		}
	}
	count, err := q.Where(whereClause).Count(issues)
	// Count Issues from the DB
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(count))
}

// Update a query to include
func requestParamToQueryFilter(query, paramValue, paramName *string) {
	initialWhereClause := *query
	if *paramValue != "" {
		*paramValue = strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(*paramValue), "[\""), "\"]")
		splitParam := strings.Split(*paramValue, ",")
		for i := range splitParam {
			splitParam[i] = strings.Trim(splitParam[i], "\"")

			switch splitParam[i] {
			case "", "undefined":
				splitParam = append(splitParam[:i], splitParam[i+1:]...)
			case "*":
				*query = initialWhereClause
				return
			}
		}
		if len(splitParam) > 0 {
			*query += " and " + *paramName + " in ("
			for i, t := range splitParam {
				if i > 0 {
					*query += ","
				}
				*query += "'" + strings.TrimSpace(t) + "'"
			}
			*query += ")"
		}
	}
}
