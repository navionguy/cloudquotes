package actions

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/navionguy/cloudquotes/models"
	"github.com/pkg/errors"
)

// AuthorsResource is the resource for the Conversation model
type AuthorsResource struct {
	buffalo.Resource
}

// List all the known authors
func (v AuthorsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	authors := &models.Authors{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".

	q := tx.Eager("Authors").PaginateFromParams(c.Params()).Order("name")

	// Retrieve all Authors from the DB
	if err := q.All(authors); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, authors))
}

// Create default implementation.
func (v AuthorsResource) Create(c buffalo.Context) error {
	s := c.Session()

	tcv := s.Get("conversation")
	escv, ok := tcv.(string)

	if ok {
		scv, _ := url.QueryUnescape(escv)
		cv := &models.Conversation{}

		_ = json.Unmarshal([]byte(scv), cv)

		fmt.Println(cv)
		fmt.Println("take that")
	}

	auth := &models.Author{}

	// Bind quote to the html form elements
	if err := c.Bind(auth); err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("author:%s\n", auth.Name)
	if nil == auth.FindByName() {
		fmt.Println("found him")
	}

	return c.Render(200, r.HTML("conversations/new"))
}
