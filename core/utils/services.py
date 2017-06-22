"""
This file is an intermediate layer between GitHub APIs and Django models.
"""

import requests
from requests.exceptions import ConnectionError

def request_github_issues(user, repo):
    """
    Returns a list of all the issues of a repository in `json` format.
    """
    try:
        api_data = 'https://api.github.com/repos/'+ user +'/' + repo + '/issues?state=all'
        response = requests.get(api_data)
        if response.status_code < 400:
            return {'error': False, 'error_type':None, 'data': response.json(),
                    'status_code': response.status_code}
        else:
            return {'error': True, 'error_type':None, 'data': response.json(),
                    'status_code': response.status_code}
    except ConnectionError:
        return {'error': True, 'error_type':ConnectionError, 'data': 'No Internet connection'}
