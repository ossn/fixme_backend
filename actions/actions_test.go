package actions

import (
	"testing"

	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action := suite.NewAction(App())

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}
