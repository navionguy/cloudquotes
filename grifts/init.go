package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/navionguy/cloudquotes/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
