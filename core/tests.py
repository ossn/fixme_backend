"""
This file contains all the tests related to `issue parser core` app.
"""

# -*- coding: utf-8 -*-
from __future__ import unicode_literals
from django.test import TestCase
from rest_framework.test import APIClient
from rest_framework import status

from .models import UserRepo, parse_issue, validate_and_store_issue, Issue, delete_closed_issues
from .utils.services import request_github_issues

SAMPLE_ISSUE = {
    "html_url": "https://github.com/mozillacampusclubs/issue_parser_backend/issues/7",
    "id": 233564738,
    "number": 7,
    "title": "Dockerize Project",
    "labels": [
        {
            "id": 613678729,
            "url": "https://api.github.com/repos/labels/enhancement",
            "name": "enhancement",
            "color": "84b6eb",
            "default": True
        }
    ],
    "state": "open",
    "created_at": "2017-06-05T11:47:01Z",
    "updated_at": "2017-06-06T10:35:57Z",
    "body": """
            Experience: Easyfix\r\nExpected-time: 3 hours\r\nLanguage: Python\r\n
            Technology-stack Django\r\n\r\n## Description\r\n
            Dockerize this backend project for development and deployment purposes.
            """
}

class UserRepoModelTestCase(TestCase):
    """This class defines the test suite for the `UserRepo` model."""

    def setUp(self):
        """Define the test client and other test variables."""
        self.user = 'razat249'
        self.repo = 'github-view'
        self.user_repo = UserRepo(user=self.user, repo=self.repo)

    def test_user_repo_model_can_create_a_userrepo(self):
        """Test the `UserRepo` model can create a `user_repo`."""
        old_count = UserRepo.objects.count()
        self.user_repo.save()
        new_count = UserRepo.objects.count()
        self.assertNotEqual(old_count, new_count)

    def test_user_repo_model_can_delete_a_userrepo(self):
        """Test the `UserRepo` model can delete a `user_repo`."""
        old_count = UserRepo.objects.count()
        self.user_repo.save()
        self.user_repo.delete()
        new_count = UserRepo.objects.count()
        self.assertEqual(old_count, new_count)


class IssueModelAndFetcherTestCase(TestCase):
    """This class defines the test suite for the `issue fetcher` component."""

    def setUp(self):
        """Initial setup for running tests."""
        pass

    def test_api_can_request_issues(self):
        """Test the request function"""
        payload = request_github_issues('razat249', 'github-view')
        self.assertEqual(payload['error'], False)
        self.assertLess(payload['status_code'], 400)

    def test_api_request_can_handle_errors(self):
        """Test the request function can handle errors"""
        # wrong repo name to test error handling.
        payload = request_github_issues('razat249', 'wrong_repo')
        self.assertEqual(payload['error'], True)
        self.assertGreaterEqual(payload['status_code'], 400)

    def test_correct_issue_parsing(self):
        """Test for correct parsing of issues"""
        issue = SAMPLE_ISSUE.copy()
        parsed = parse_issue(issue['body'])
        for item in parsed:
            self.assertTrue(item)

    def test_validate_and_store_issue(self):
        """Test for validating and storing issues."""
        old_count = Issue.objects.count()
        validate_and_store_issue(SAMPLE_ISSUE)
        new_count = Issue.objects.count()
        self.assertNotEqual(old_count, new_count)

    def test_api_can_delete_closed_issues_in_db(self):
        """Test for checking if issues are deleted when closed online but present in db"""
        issue = SAMPLE_ISSUE.copy()
        validate_and_store_issue(issue)
        issue['state'] = 'closed'
        old_count = Issue.objects.count()
        delete_closed_issues(issue)
        new_count = Issue.objects.count()
        self.assertLess(new_count, old_count)

class ViewTestCase(TestCase):
    """This class defines the test suite for the api views."""

    def setUp(self):
        """Define the test client and other test variables."""
        self.client = APIClient()
        issue = SAMPLE_ISSUE.copy()
        validate_and_store_issue(issue)

    def test_api_can_get_metadata(self):
        """Test the api can get given metadata."""
        response = self.client.get('/metadata/', format="json")
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_api_can_get_issues_list(self):
        """Test the api can get given issues list."""
        response = self.client.get('/issues/', format="json")
        self.assertEqual(response.status_code, status.HTTP_200_OK)
