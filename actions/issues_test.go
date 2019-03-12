package actions

import (
	"encoding/json"
	"math/rand"

	"github.com/gobuffalo/uuid"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/suite"
	"github.com/ossn/fixme_backend/models"
)

type IssuesSuite struct {
	*suite.Action
	Issues *models.Issues
}

func (is *IssuesSuite) CreateIssue(issue models.Issue) {
	is.NoError(is.DB.Create(issue))
}

func (as *IssuesSuite) Test_IssuesResource_List() {
	issue := models.Issue{
		Title: nulls.String{"Test issue", true},
	}
	as.CreateIssue(issue)
	res := as.JSON("/issues").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), issue.Title)
}

func (as *IssuesSuite) Test_IssuesResource_Show() {
	idStr := "asdfasdfasdfasdf"
	var id [16]byte
	copy(id[:], []byte(idStr))
	issue := models.Issue{
		Title: nulls.String{"Test issue", true},
		ID:    uuid.UUID(id),
	}
	as.CreateIssue(issue)
	res := as.JSON("/issues/" + idStr).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), issue.Title)
}

func (as *IssuesSuite) Test_IssuesResource_Count() {
	issue := models.Issue{
		Title: nulls.String{"Test issue", true},
	}
	randCount := rand.Intn(100)
	for i := 0; i < randCount; i++ {
		as.CreateIssue(issue)
	}
	res := as.JSON("/issues-count").Get()
	var count int
	as.NoError(json.Unmarshal(res.Body.Bytes(), &count))
	as.Equal(200, res.Code)
	as.Equal(randCount, count)
}
