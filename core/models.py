# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models
from core.utils.services import request_github_issues

from celery.task.schedules import crontab
from celery.decorators import periodic_task
from datetime import timedelta
from lxml import html

ISSUE_UPDATE_PERIOD = 15 # in minutes

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
    EASYFIX = 'easyfix'
    MODERATE = 'moderate'
    SENIOR = 'senior'
    EXPERIENCE_NEEDED_CHOICES = (
        (EASYFIX, 'easyfix'),
        (MODERATE, 'moderate'),
        (SENIOR, 'senior'),
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
    issue_body = models.TextField()

    class Meta:
        ordering = ('updated_at',) # Ascending order according to updated_at.


@periodic_task(run_every=timedelta(minutes=ISSUE_UPDATE_PERIOD), name="periodic_issues_updater")
def periodic_issues_updater():
    """
    Update `Issue` model in the database in every 
    `ISSUE_UPDATE_PERIOD` minutes.
    """
    list_of_repos = UserRepo.objects.values('user', 'repo',)
    for repo in list_of_repos:
        issue_list = request_github_issues(repo['user'], repo['repo'])
        for issue in issue_list:
            validate_and_store_issue(issue)

def validate_and_store_issue(issue):
    """
    Validate issue:
    if valid - store it into data base,
    else - Do not store in database
    """
    tree = html.fromstring(issue['body'])

    # Parse issue from the tags.
    experience_needed_list = tree.xpath('experience')
    language_list = tree.xpath('language')
    expected_time_list = tree.xpath('expected-time')
    technology_stack_list = tree.xpath('technology-stack')
    description_list = tree.xpath('description')

    if experience_needed_list and language_list and expected_time_list and technology_stack_list and description_list:
        experience_needed = experience_needed_list[0].text_content().strip().lower()
        language = language_list[0].text_content().strip().lower()
        expected_time = expected_time_list[0].text_content().strip().lower()
        technology_stack = technology_stack_list[0].text_content().strip().lower()
        description = description_list[0].text_content()
        issue_instance = Issue(issue_id=issue['id'], title=issue['title'], experience_needed=experience_needed, expected_time=expected_time, language=language, tech_stack=technology_stack, created_at=issue['created_at'], updated_at=issue['updated_at'], issue_number=issue['number'], issue_url=issue['html_url'], issue_body=description )
        issue_instance.save()
    else :
        print 'Issue with id ' + str(issue['id']) + ' is not valid for our system.'
