# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models
from core.utils.services import request_github_issues

from celery.task.schedules import crontab
from celery.decorators import periodic_task
from datetime import timedelta

class UserRepo(models.Model):
    """
    UserRepo model is used to store the username and repo-name
    for a repository.
    """
    user = models.CharField(max_length=100)
    repo = models.CharField(max_length=100)
    created = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ('created',) # Ascending order according to date created.
        unique_together = ("user", "repo") # Avoid repo duplicates.


class IssueLabel(models.Model):
    """
    Label model for storing labels of an issue.
    """
    label_id = models.IntegerField(primary_key=True)
    label_url = models.URLField()
    label_name = models.CharField(max_length=100)
    label_color = models.CharField(max_length=6)
    
    class Meta:
        ordering = ('label_name',) # Ascending order according to label_name.


class Issue(models.Model):
    """
    Issue model is used to store github issues.
    """
    # Setting choices for experience needed to solve an issue.
    EASYFIX = 'Easyfix'
    MODERATE = 'Moderate'
    SENIOR = 'Senior'
    EXPERIENCE_NEEDED_CHOICES = (
        (EASYFIX, 'Easyfix'),
        (MODERATE, 'Moderate'),
        (SENIOR, 'Senior'),
    )
    # Model attributes start from here.
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
    issue_labels = models.ManyToManyField(IssueLabel, blank=True)
    issue_url = models.URLField()

    class Meta:
        ordering = ('updated_at',) # Ascending order according to updated_at.


@periodic_task(run_every=timedelta(seconds=35), name="pupu")
def pupu():
    list_of_repos = UserRepo.objects.values('user', 'repo',)
    for repo in list_of_repos:
        issue_data = request_github_issues(repo['user'], repo['repo'])
        for i in issue_data:
            issue = Issue(issue_id=i['id'], title=i['title'], expected_time="22 hrs", language="cofeeScript", tech_stack="Django", created_at=i['created_at'], updated_at=i['updated_at'], issue_number=i['number'], issue_url=i['html_url'] )
            issue.save()
