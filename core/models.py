# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models

# UserRepo model is used to store the username and repo-name for a repository.
class UserRepo(models.Model):
    user = models.CharField(max_length=100)
    repo = models.CharField(max_length=100)
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ('created',)

# Label model for storing labels of an issue.
class IssueLabel(models.Model):
    label_id = models.IntegerField(primary_key=True)
    label_url = models.URLField()
    label_name = models.CharField(max_length=100)
    label_color = models.CharField(max_length=6)
    
    class Meta:
        ordering = ('label_name',)

# Issue model is used to store github issues.
class Issue(models.Model):
    EASYFIX = 'Easyfix'
    MODERATE = 'Moderate'
    SENIOR = 'Senior'
    EXPERIENCE_NEEDED_CHOICES = (
        (EASYFIX, 'Easyfix'),
        (MODERATE, 'Moderate'),
        (SENIOR, 'Senior'),
    )


    issue_id = models.IntegerField(primary_key=True)
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
    issue_number = models.IntegerField()
    issue_labels = models.ManyToManyField(IssueLabel)
    issue_url = models.URLField()

    class Meta:
        ordering = ('updated_at',)
