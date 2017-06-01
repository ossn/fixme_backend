# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models
import uuid

# UserRepo model is used to store the username and repo-name for a repository.
class UserRepo(models.Model):
    user = models.CharField(max_length=100)
    repo = models.CharField(max_length=100)
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ('created',)

# Label model for storing labels of an issue.
class IssueLabel(models.Model):
    label_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    label_url = models.CharField(max_length=2000)
    label_name = models.CharField(max_length=100)
    label_color = models.CharField(max_length=10)

# Issue model is used to store github issues.
class Issue(models.Model):
    EASYFIX = 0
    MODERATE = 1
    SENIOR = 2
    EXPERIENCE_NEEDED_CHOICES = (
        (EASYFIX, 'Easyfix'),
        (MODERATE, 'Moderate'),
        (SENIOR, 'Senior'),
    )
    issue_id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    title = models.CharField(max_length=100)
    experience_needed = models.CharField(
        max_length=10,
        choices=EXPERIENCE_NEEDED_CHOICES,
        default=MODERATE,
    )
    expected_time = models.CharField(max_length=100)
    language = models.CharField(max_length=100)
    tech_stack = models.CharField(max_length=100)
    created_at = models.DateTimeField()
    updated_at = models.DateTimeField()
    issue_details = models.TextField()
    issue_number = models.IntegerField()
    issue_labels = models.ForeignKey(IssueLabel, on_delete=models.CASCADE)
    issue_url = models.CharField(max_length=2000)
