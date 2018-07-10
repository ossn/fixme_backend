"""
This file is an intermediate layer between GitHub APIs and Django models.
"""

import requests
from requests.exceptions import ConnectionError

# FIXME: Parse more than 30 issues!!!


def get_repo_and_username(link):
    params = link.split("/")
    return params[len(params)-2], params[len(params)-1]


def request_to_github(repo, username, url_extension):
    """
    Returns a list of allthe issues of a repository in `json` format.
    """
    try:
        api_data = 'https://api.github.com/repos/' + \
            repo + '/' + \
            username + url_extension
        response = requests.get(api_data)
        if response.status_code < 400:
            return {'error': False, 'error_type': None, 'data': response.json(),
                    'status_code': response.status_code}
        else:
            return {'error': True, 'error_type': None, 'data': response.json(),
                    'status_code': response.status_code}
    except ConnectionError:
        return {'error': True, 'error_type': ConnectionError, 'data': 'No Internet connection'}
