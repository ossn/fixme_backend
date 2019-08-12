package worker_github

type (
	/**
	* GraphQL Types
	 */
	PageInfo struct {
		StartCursor     string
		HasPreviousPage bool
	}

	Issues struct {
		Nodes []struct {
			Title      string
			Body       string
			Closed     bool
			Number     int
			URL        string
			CreatedAt  string
			DatabaseID int
			Labels     struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"labels(first:100)"`
		}
		PageInfo PageInfo
	}

	repoQuery struct {
		Repository struct {
			Object struct {
				Blob struct {
					Text string `graphql:"text"`
				} `graphql:"... on Blob"`
			} `graphql:"object(expression: \"master:README.md\")"`
			Description string
			RepositoryTopics struct {
				Nodes []struct {
					Topic struct {
						Name string `graphql:"name"`
					} `graphql:"topic"`
				} `graphql:"nodes"`
			} `graphql:"repositoryTopics(first: 100)"`
			Languages struct {
				Edges []struct{
					Size int `graphql:"size"`
					Node struct {
						Name string `graphql:"name"`
					} `graphql:"node"`
				} `graphql:"edges"`
				TotalSize int `graphql:"totalSize"`
			} `graphql:"languages(first:100, orderBy: {field: SIZE, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	initialIssueQuery struct {
		Repository struct {
			Issues Issues `graphql:"issues(last: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	issueQueryWithBefore struct {
		Repository struct {
			Issues Issues `graphql:"issues(last: 100, before: $before)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}


	issueStatusQuery struct {
		Repository struct {
			Issue struct {
				Closed bool `graphql:"issues(last: 100)"`
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	rateLimitQuery struct {
		RateLimit struct {
			Remaining int    `graphql:"remaining"`
			ResetAt   string `graphql:"resetAt"`
		} `graphql:"rateLimit"`
	}
)
