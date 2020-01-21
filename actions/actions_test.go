package actions

import (
	"context"
	"testing"

	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action := suite.NewAction(App(context.TODO()))

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}
