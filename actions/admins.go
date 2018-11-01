package actions

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/mw-tokenauth"
	"github.com/gobuffalo/pop"
	"github.com/ossn/fixme_backend/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AdminsResource is the resource for the Admin model
type AdminsResource struct {
	buffalo.Resource
}

// List gets all Admins. This function is mapped to the path
// GET /admins
func (v AdminsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	admins := &models.Admins{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Admins from the DB
	if err := q.All(admins); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.JSON(admins))
}

// Show gets the data for one Admin. This function is mapped to
// the path GET /admins/{admin_id}
func (v AdminsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	// To find the Admin the parameter admin_id is used.
	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(admin))
}

// New renders the form for creating a new Admin.
// This function is mapped to the path GET /admins/new
func (v AdminsResource) New(c buffalo.Context) error {
	return c.Render(200, r.JSON(&models.Admin{}))
}

// Create adds a Admin to the DB. This function is mapped to the
// path POST /admins
func (v AdminsResource) Create(c buffalo.Context) error {
	// Allocate an empty Admin
	admin := &models.Admin{}

	// Bind admin to the html form elements
	if err := c.Bind(admin); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	h, err := generatePasswordHash(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = h
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(admin)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(admin))
	}

	return c.Render(201, r.JSON(admin))
}

// Edit renders a edit form for a Admin. This function is
// mapped to the path GET /admins/{admin_id}/edit
func (v AdminsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(admin))
}

// Update changes a Admin in the DB. This function is mapped to
// the path PUT /admins/{admin_id}
func (v AdminsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return c.Error(404, err)
	}

	oldPassword := admin.Password
	// Bind Admin to the html form elements
	if err := c.Bind(admin); err != nil {
		return errors.WithStack(err)
	}

	if oldPassword != admin.Password {
		h, err := generatePasswordHash(admin.Password)
		if err != nil {
			return err
		}
		admin.Password = h
	}

	verrs, err := tx.ValidateAndUpdate(admin)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the response
		c.Set("errors", verrs)

		return c.Render(422, r.JSON(admin))
	}

	return c.Render(200, r.JSON(admin))
}

// Destroy deletes a Admin from the DB. This function is mapped
// to the path DELETE /admins/{admin_id}
func (v AdminsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	// To find the Admin the parameter admin_id is used.
	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(admin); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(admin))
}

func generatePasswordHash(s string) (string, error) {

	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(hashedBytes[:]), nil
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password`
}

func (v AdminsResource) Login(c buffalo.Context) error {

	// Allocate an empty login
	login := &LoginForm{}

	if err := c.Bind(login); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	user := &models.Admin{}
	err := tx.Where("email = ?", login.Email).First(user)

	if err != nil {
		return c.Error(401, errors.Wrap(err, "Couldn't find user with this email and password"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return c.Error(401, errors.Wrap(err, "Couldn't find user with this email and password"))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	signKey, err := tokenauth.GetHMACKey(jwt.SigningMethodHS256)
	if err != nil {
		return errors.Wrap(err, "Couldn't get hmac key")
	}
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return errors.Wrap(err, "Couldn't sign jwt")
	}
	return c.Render(200, r.JSON(map[string]string{"jwt": tokenString}))

}
