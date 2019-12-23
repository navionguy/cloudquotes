package actions

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/navionguy/cloudquotes/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Conversation)
// DB Table: Plural (conversations)
// Resource: Plural (Conversations)
// Path: Plural (/conversations)
// View Template Folder: Plural (/templates/conversations/)

var fontScale = map[int]string{
	1: "48px",
	2: "40px",
	3: "30px",
	4: "20px",
	5: "10px",
}

// ConversationsResource is the resource for the Conversation model
type ConversationsResource struct {
	buffalo.Resource
}

// List gets all Conversations. This function is mapped to the path
// GET /conversations
func (v ConversationsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	conversations := &models.Conversations{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	// I only eager load the Quotes because I don't touch data from the
	// other objects in the index page

	q := tx.Eager("Quotes").Eager("Quotes.Author").PaginateFromParams(c.Params()).Order("occurredon DESC")

	// Retrieve all Conversations from the DB
	if err := q.All(conversations); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, conversations))
}

// Show gets the data for one Conversation. This function is mapped to
// the path GET /conversations/{conversation_id}
func (v ConversationsResource) Show(c buffalo.Context) error {

	// Go get my conversation
	conversation, err := v.loadConversation(c)
	if err != nil {
		return c.Error(404, err)
	}

	var fontSize string
	if len(conversation.Quotes) > 1 {
		fontSize = fmt.Sprintf("%s", fontScale[len(conversation.Quotes)])
	} else if len(conversation.Quotes[0].Phrase) > 100 {
		fontSize = fmt.Sprintf("%s", fontScale[2])
	}

	c.Set("fontsize", fontSize)
	return c.Render(200, r.Auto(c, conversation))
}

// New renders the form for creating a new Conversation.
// This function is mapped to the path GET /conversations/new
func (v ConversationsResource) New(c buffalo.Context) error {
	conversation := &models.Conversation{}

	conversation.OccurredOn = time.Now()
	conversation.Publish = true
	err := v.loadForm(conversation, c)

	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.HTML("conversations/new.html"))
}

// Create is mapped to the path POST /conversations.
//
// To support the full process of creating a conversation, the
// browser can POST with a collection of intentions.  This is
// communicated via a hidden field named "option".  This contains
// a single string value that controls what happens with the
// forms data.
//
// "save" - take the conversation data and save it to the database
//
// "reply" - bring up the next quote in the conversation if there
// is one.  If not, he wants to add one.
//
// "addAuthor" - A new author needs to be added to the available
// list of authors.
//
// "prevQuote" - Move to the previous quote in the conversation
//
func (v ConversationsResource) Create(c buffalo.Context) error {

	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return errors.WithStack(err)
	}

	conv, option, err := v.bindToForm(c)
	fmt.Printf("bind complete, err %s, option %s/n", err, *option)

	if err != nil {
		return errors.WithStack(err)
	}

	switch *option {
	case "addAuthor":
		return v.addAuthor(conv, c)

	case "save":
		verrs, err := conv.Create()

		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			err = v.loadForm(conv, c)

			if err != nil {
				return errors.WithStack(err)
			}
			// set the verification errors into the context and send back the quote
			c.Set("errors", verrs)

			return c.Render(422, r.HTML("conversations/new.html"))
		}
		c.Flash().Add("success", "Conversation was created successfully")

		//return c.Redirect(302, fmt.Sprintf("/conversations//%%7B%s%%7D/", conversation.ID.String()))
		return c.Render(201, r.Auto(c, conv))
	}

	return err
}

// Edit renders a edit form for a Conversation. This function is
// mapped to the path GET /conversations/{conversation_id}/edit
func (v ConversationsResource) Edit(c buffalo.Context) error {

	cv, err := v.loadConversation(c)

	if err != nil {
		return c.Error(404, err)
	}

	err = v.loadForm(cv, c)

	if err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, cv))
}

// Update changes a Conversation in the DB. This function is mapped to
// the path PUT /conversations/{conversation_id}
func (v ConversationsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return errors.WithStack(err)
	}

	conv, option, err := v.bindToForm(c)

	if err != nil {
		return err
	}

	// Allocate an empty quote
	quote := &models.Quote{}

	if err := tx.Find(conv, c.Param("conversation_id")); err != nil {
		return c.Error(404, err)
	}
	var verrs *validate.Errors

	switch *option {
	case "addAuthor":
		return v.addAuthor(conv, c)

	case "save":
		conv, quote, verrs, err = v.saveConversation(quote)

		if err != nil {
			return err
		}

		if verrs.HasAny() {
			err = v.loadForm(conv, c)

			if err != nil {
				return errors.WithStack(err)
			}
			// set the verification errors into the context and send back the quote
			c.Set("errors", verrs)

			return c.Render(422, r.HTML("conversations/new.html"))
		}
		c.Flash().Add("success", "Conversation was created successfully")

		//return c.Redirect(302, fmt.Sprintf("/conversations//%%7B%s%%7D/", conversation.ID.String()))
		return c.Render(201, r.Auto(c, conv))
	}

	return err
}

// Destroy deletes a Conversation from the DB. This function is mapped
// to the path DELETE /conversations/{conversation_id}
func (v ConversationsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Conversation
	conversation := &models.Conversation{}

	// To find the Conversation the parameter conversation_id is used.
	if err := tx.Eager("Quotes").Find(conversation, c.Param("conversation_id")); err != nil {
		return c.Error(404, err)
	}

	// loop through all the quotes and delete them
	for i := range conversation.Quotes {
		q := &models.Quote{}
		q.ID = conversation.Quotes[i].ID
		if err := tx.Destroy(q); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := tx.Destroy(conversation); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Conversation was destroyed successfully")

	// Redirect to the conversations index page
	return c.Render(302, r.Auto(c, conversation))

}

// Export dumps the database in JSON.  Maps to the
// path GET /conversations/export
func (v ConversationsResource) Export(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	conversations := &models.Conversations{}

	if err := tx.Eager("Quotes.Conversation").Eager("Quotes").Eager("Quotes.Author").Eager("Quotes.Annotation").All(conversations); err != nil {
		return c.Error(404, err)
	}

	// Redirect to the conversations index page

	//return c.Redirect(301, "/conversations")

	return c.Render(200, r.JSON(conversations))
}

// saveConversation - the user has finished adding quotes and is ready to save the conversation
func (v ConversationsResource) saveConversation(quote *models.Quote) (*models.Conversation, *models.Quote, *validate.Errors, error) {

	conversation := &models.Conversation{}

	if quote.Sequence > 0 {
		verrs, err := conversation.Create()

		if err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}

		if verrs.HasAny() {
			return nil, quote, verrs, nil
		}
	}

	return conversation, nil, nil, nil
}

func (v ConversationsResource) nextQuote(c buffalo.Context) (*models.Conversation, error) {
	return nil, nil
}

// loadConversation handles loading a quote for the Show() function.
// I may push this back into the function unless I figure out a better way to print.
func (v ConversationsResource) loadConversation(c buffalo.Context) (*models.Conversation, error) {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return nil, errors.WithStack(errors.New("no transaction found"))
	}
	conversation := models.Conversation{}

	// Need to tell buffalo to "Eager" load all the objects contained
	// in the conversation object.
	// To find the Conversation the parameter conversation_id is used.

	if err := tx.Eager("Quotes.Conversation").Eager("Quotes").Eager("Quotes.Author").Eager("Quotes.Annotation").Find(&conversation, c.Param("conversation_id")); err != nil {
		return nil, c.Error(404, err)
	}

	// I have not yet figured out how to detect a null pointer in
	// my plush code embedded in the HTML.  Until I do, I build
	// a list of strings that are either empty, or contain any
	// annotations on a quote.

	var notes []string
	for _, quote := range conversation.Quotes {
		if quote.Annotation != nil {
			notes = append(notes, "* "+quote.Annotation.Note)
		} else {
			notes = append(notes, "")
		}
	}

	c.Set("notes", notes)

	return &conversation, nil
}

func (v ConversationsResource) loadForm(conversation *models.Conversation, c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	authors := []models.Author{}
	annotation := models.Annotation{}

	annotation.Note = ""

	// Retrieve all Authors from the DB
	if err := tx.Order("name").All(&authors); err != nil {
		return errors.WithStack(err)
	}

	cvjson, err := conversation.MarshalConversation()

	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("conversation", conversation)
	c.Set("authors", authors)
	c.Set("cvj", cvjson)

	return nil
}

// build a conversation from the standard conversation form
func (v ConversationsResource) bindToForm(c buffalo.Context) (*models.Conversation, *string, error) {
	conv := &models.Conversation{}

	cvjson := c.Request().Form.Get("cvjson")
	err := conv.UnmarshalConversation(cvjson)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	option := c.Request().Form.Get("option")

	return conv, &option, nil
}

// if there is an annotation, add it to the quote
func attachAnnotation(quote *models.Quote, annotation *models.Annotation) error {

	// get out the annotation value

	if len(annotation.Note) == 0 {
		quote.Annotation = nil
		quote.AnnotationID = nil
		return nil
	}

	quote.Annotation = annotation
	err := quote.Annotation.FindByNote() // looks for the Note text in the database, no repeats!

	if err != nil {
		return err
	}

	// the quotes annotationID needs to either point to the annotation's id, or be nil
	// can't use uuid.Nil because buffalo tries to write that to the database as
	// a zero guid
	if quote.Annotation.ID == uuid.Nil {
		quote.AnnotationID = nil
	} else {
		quote.AnnotationID = &quote.Annotation.ID
	}

	return nil
}

// User wants to add an author to the database.  Save where we are and go do that.
func (v ConversationsResource) addAuthor(conv *models.Conversation, c buffalo.Context) error {
	author := &models.Author{}

	cvjson, err := conv.MarshalConversation()

	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("author", author)
	c.Set("cvj", cvjson)

	return c.Render(200, r.HTML("authors/new"))
}
