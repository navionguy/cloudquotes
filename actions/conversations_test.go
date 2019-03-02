package actions

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.wdf.sap.corp/i826414/quotearchive/models"
)

func (as *ActionSuite) Test_Conversations_Show() {
	as.LoadFixture("test quotes")

	cvs := []models.Conversation{}

	err := as.DB.All(&cvs)

	if err != nil {
		as.Fail("conversation show failed to load fixture", err.Error())
	}

	if len(cvs) == 0 {
		as.Fail("no conversations found", "no conversations loaded")
	}

	q := models.Quote{}
	err = as.DB.Where("conversation_id = ?", cvs[0].ID).First(&q)

	if err != nil {
		as.Fail("couldn't load quote", err.Error())
	}

	if len(q.Phrase) == 0 {
		as.Fail("quote note found", "first quote")
	}

	p := fmt.Sprintf("/conversations/%s", cvs[0].ID)
	res := as.HTML(p).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), q.Phrase)
}

func (as *ActionSuite) Test_Conversations_Show_BadID() {
	as.LoadFixture("test quotes")

	q, _ := uuid.NewV4()
	p := fmt.Sprintf("/conversations/%s", q.String())
	res := as.HTML(p).Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_Conversations_Index() {
	res := as.HTML("/conversations").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Conversations")
}

func (as *ActionSuite) Test_Conversations_New() {
	res := as.HTML("/conversations/new").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Quote")
	as.Contains(res.Body.String(), "Publish Quote?")
}

func (as *ActionSuite) Test_Conversations_Create() {
	/*	as.LoadFixture("test quotes")

			authors := []models.Author{}

			err := as.DB.All(&authors)

			if err != nil {
				as.FailNow("error getting authors", err.Error())
			}

			if len(authors) == 0 {
				as.FailNow("no authors found", "no test authors")
			}

		cv := models.Quote{}

		res := as.HTML("/quotes").Post(&cv)
		as.Equal(301, res.Code)*/
}

func (as *ActionSuite) Test_Conversations_Delete() {
	as.LoadFixture("test quotes")

	cvs := []models.Conversation{}

	err := as.DB.All(&cvs)

	if err != nil {
		as.Fail("conversation show failed to load fixture", err.Error())
	}

	if len(cvs) == 0 {
		as.Fail("no conversations found", "no conversations loaded")
	}

	p := fmt.Sprintf("/conversations/%s", cvs[0].ID)
	res := as.HTML(p).Delete()
	as.Equal(302, res.Code)

	q := []models.Quote{}
	err = as.DB.Where("conversation_id = ?", cvs[0].ID).All(&q)

	if err != nil {
		as.Fail("couldn't load quote", err.Error())
	}

	if len(q) > 0 {
		as.Fail("quotes didn't delete", "still found quote")
	}

}

func (as *ActionSuite) Test_Conversations_Delete_BadID() {
	as.LoadFixture("test quotes")

	q, _ := uuid.NewV4()
	p := fmt.Sprintf("/conversations/%s", q.String())
	res := as.HTML(p).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_Conversations_Export() {
	as.LoadFixture("test quotes")

	res := as.HTML("/conversations/export").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Conversations_saveConversation() {
	as.LoadFixture("test quotes")

	authors := []models.Author{}

	err := as.DB.All(&authors)

	if err != nil {
		as.FailNow("error getting authors", err.Error())
	}

	if len(authors) == 0 {
		as.FailNow("no authors found", "no test authors")
	}

	var cv ConversationsResource

	qt := models.Quote{
		SaidOn:   time.Now(),
		Sequence: 0,
		Phrase:   "A test quote.",
		Publish:  true,
		AuthorID: authors[0].ID,
	}

	_, _, verrs, err := cv.saveConversation(&qt)

	if err != nil {
		as.Fail("couldn't save valid quote", err.Error())
	}

	if verrs != nil {
		as.Fail("saveConversation got validation errors", verrs.String())
	}
}

func (as *ActionSuite) Test_Conversations_saveConversationValidationErrors() {
	as.LoadFixture("test quotes")

	qt := models.Quote{}

	var cv ConversationsResource

	_, _, verrs, err := cv.saveConversation(&qt)

	if err != nil {
		as.Fail("couldn't save valid quote", err.Error())
	}

	if verrs == nil {
		as.Fail("saveConversation ignored validation errors", verrs.String())
	}
}
