"""
This file is an intermediate layer between GitHub APIs and Django models.
"""

import requests

def request_github_issues(user, repo):
    """
    Returns a list of all the issues of a repository in `json` format.
    """
    api = 'https://api.github.com/repos/'+ user +'/' + repo + '/issues'
    response = requests.get(api)
    return response.json()
