"""
This file contains all the tests related to `issue parser core` app.
"""

# -*- coding: utf-8 -*-
from __future__ import unicode_literals
from django.test import TestCase

from .models import UserRepo

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
