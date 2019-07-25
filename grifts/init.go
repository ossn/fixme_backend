package grifts

import (
	"context"

	"github.com/gobuffalo/buffalo"
	"github.com/ossn/fixme_backend/actions"
)

func init() {
	buffalo.Grifts(actions.App(context.TODO()))
}
