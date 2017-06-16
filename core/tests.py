"""
This file contains all the tests related to `issue parser core` app.
"""

# -*- coding: utf-8 -*-
from __future__ import unicode_literals
from django.test import TestCase

from .models import UserRepo
from .utils.services import request_github_issues

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


class IssueFetcherTestCase(TestCase):
    """This class defines the test suite for the `issue fetcher` component."""

    def setUp(self):
        pass

    def test_api_can_request_issues(self):
        """Test the request function"""
        payload = request_github_issues('razat249', 'github-view')
        self.assertEqual(payload['error'], False)
        self.assertLess(payload['status_code'], 400)
    
    def test_api_can_handle_errors_while_requesting_issues(self):
        """Test the request function can handle errors"""
        payload = request_github_issues('razat249', 'wrong_repo') # wrong repo name to test error handling.
        self.assertEqual(payload['error'], True)
        self.assertGreaterEqual(payload['status_code'], 400)
