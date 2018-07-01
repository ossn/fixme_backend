"""
This file is an intermediate layer between GitHub APIs and Django models.
"""

import requests
from requests.exceptions import ConnectionError

# FIXME: Parse more than 30 issues!!!


def request_github_issues(repository_link):
    """
    Returns a list of all the issues of a repository in `json` format.
    """
    try:
        params = repository_link.split("/")
        api_data = 'https://api.github.com/repos/' + \
            params[len(params)-2] + '/' + \
            params[len(params)-1] + '/issues?state=open'
        response = requests.get(api_data)
        if response.status_code < 400:
            return {'error': False, 'error_type': None, 'data': response.json(),
                    'status_code': response.status_code}
        else:
            return {'error': True, 'error_type': None, 'data': response.json(),
                    'status_code': response.status_code}
    except ConnectionError:
        return {'error': True, 'error_type': ConnectionError, 'data': 'No Internet connection'}
