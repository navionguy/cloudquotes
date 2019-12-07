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

	cv, ok := v.unMarshalConversation(s)

	speaker := &models.Author{}

	// Bind quote to the html form elements
	if err := c.Bind(speaker); err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("new speaker %s\n", speaker.Name)

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	if nil != speaker.FindByName() {
		verrs, err := tx.ValidateAndCreate(speaker)

		if err != nil {
			return err
		}

		if verrs.HasAny() {
			c.Set("author", speaker)
			c.Set("gotoPage", "new")

			// set the verification errors into the context and send back the author
			c.Set("errors", verrs)

			return c.Render(422, r.Auto(c, speaker))
		}
		c.Flash().Add("success", "Speaker created successfully!")
	}

	// put the conversation back into the form
	quote := models.Quote{}
	authors := []models.Author{}
	annotation := models.Annotation{}

	annotation.Note = ""

	if quote.Annotation != nil {
		annotation.Note = quote.Annotation.Note
	}

	// Retrieve all Authors from the DB
	if err := tx.Order("name").All(&authors); err != nil {
		return errors.WithStack(err)
	}

	quote = cv.Quotes[0]

	if quote.Annotation != nil {
		annotation.Note = quote.Annotation.Note
	}

	c.Set("conversation", cv)
	c.Set("quote", quote)
	c.Set("authors", authors)
	c.Set("annotation", annotation)
	c.Set("option", "save")
	//c.Set("cvjson", scv)

	return c.Render(200, r.HTML("conversations/new"))
}

func (v AuthorsResource) unMarshalConversation(s *buffalo.Session) (*models.Conversation, bool) {
	cvv := s.Get("conversation")
	cvesc, ok := cvv.(string)

	// if no conversations in the session, can't unMarshal it
	if !ok {
		return nil, false
	}

	cvs, err := url.QueryUnescape(cvesc)

	if err != nil {
		return nil, false
	}

	cv := &models.Conversation{}
	err = json.Unmarshal([]byte(cvs), cv)

	if err != nil {
		return nil, false
	}

	return cv, true
}
