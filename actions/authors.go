package actions

import (
	"fmt"

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
	fmt.Println("List Authors")

	authors := &models.Authors{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".

	// Get all the authors names and their quote count
	// SELECT authors.name, COUNT(DISTINCT quotes.id) FROM authors LEFT JOIN quotes ON quotes.author_id = authors.id GROUP BY authors.id ORDER BY authors.name;
	//models.DB.Q().LeftJoin("authors", "authors.id").LeftJoin("quotes", "quotes.author_id").
	//Where("quotes.author_id = authors.id").PaginateFromParams()
	q := tx.PaginateFromParams(c.Params()).Order("name")
	tq := tx.PaginateFromParams(c.Params()).Order("name")
	cq := tq.RawQuery("SELECT authors.id, authors.name, COUNT(DISTINCT quotes.id) FROM authors LEFT JOIN quotes ON quotes.author_id = authors.id GROUP BY authors.id ORDER BY authors.name")
	// SELECT authors.name, COUNT(DISTINCT quotes.id) FROM authors LEFT JOIN quotes ON quotes.author_id = authors.id GROUP BY authors.id ORDER BY authors.name;
	//models.DB.Q().LeftJoin("authors", "authors.id").LeftJoin("quotes", "quotes.author_id").
	//Where("quotes.author_id = authors.id").PaginateFromParams()
	//tx.Select("authors.name", "COUNT(DISTINCT quotes.id)")
	//models.DB.LeftJoin("roles", "roles.id=user_roles.role_id").LeftJoin("users u", "u.id=user_roles.user_id").
	//Where(`roles.name like ?`, name).Paginate(page, perpage)

	authorCredits := &models.AuthorCredits{}

	if err := cq.All(authorCredits); err != nil {
		fmt.Println("Couldn't get credits")
		fmt.Println(err.Error())
		return errors.WithStack(err)
	} else {
		ac := *authorCredits
		fmt.Printf("%s %s has %d\n", ac[0].ID, ac[0].Name, ac[0].Count)
	}

	// Retrieve all Authors from the DB
	if err := q.All(authors); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", tq.Paginator)

	return c.Render(200, r.Auto(c, authorCredits))
}

// New author about to be entered
func (v AuthorsResource) New(c buffalo.Context) error {
	fmt.Println("Author->New")

	spkr := &models.Author{}

	c.Set("author", spkr)
	c.Set("cvj", "")

	return c.Render(200, r.HTML("authors/new.html"))
}

// Create default implementation.
func (v AuthorsResource) Create(c buffalo.Context) error {

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

	cvjson := c.Request().Form.Get("cvjson")

	fmt.Printf("json length %d, %s\n", len(cvjson), cvjson)

	if len(cvjson) == 0 {
		return c.Redirect(201, "authors")
	}

	// put the conversation back into the form
	conv := models.Conversation{}
	err := conv.UnmarshalConversation(cvjson)

	if err != nil {
		return errors.WithStack(err)
	}

	authors := []models.Author{}

	// Retrieve all Authors from the DB
	if err := tx.Order("name").All(&authors); err != nil {
		return errors.WithStack(err)
	}

	c.Set("conversation", conv)
	c.Set("authors", authors)
	c.Set("cvj", cvjson)

	return c.Render(200, r.HTML("conversations/new"))
}

func (v AuthorsResource) unMarshalConversation(c buffalo.Context) (*models.Conversation, bool) {
	cvv := c.Request().Form.Get("cvjson")
	conv := &models.Conversation{}
	err := conv.UnmarshalConversation(cvv)

	if err != nil {
		return nil, false
	}

	return conv, true
}
