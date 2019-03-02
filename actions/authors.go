package actions

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/navionguy/cloudquotes/models"
)

// AuthorsResource is the resource for the Conversation model
type AuthorsResource struct {
	buffalo.Resource
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
	}
	return c.Render(200, r.HTML("authors/put.html"))
}
