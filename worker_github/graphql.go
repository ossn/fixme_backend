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

	language struct {
		Repository struct {
			PrimaryLanguage struct {
				Name string
			}
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

	tagsQuery struct {
		Repository struct {
			RepositoryTopics struct {
				Nodes []struct {
					Topic struct {
						Name string
					}
				}
			} `graphql:"repositoryTopics(first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	issueStatusQuery struct {
		Repository struct {
			Issue struct {
				Closed bool
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
