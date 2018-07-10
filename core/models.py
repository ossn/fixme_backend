# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from datetime import timedelta
from django.db import models
from django_mysql.models import ListTextField
from core.utils.services import request_to_github, get_repo_and_username
from celery.decorators import periodic_task, task
import logging

ISSUE_UPDATE_PERIOD = 15  # in minutes
"""
 These labels are going count as "easy" for the experience needed
"""
EASY_LABELS = ["help_wanted", "good first issue", "easyfix", "easy"]

"""
 These labels are going count as "moderate" for the experience needed
"""
MODERATE_LABELS = ["moderate"]

"""
These labels are going count as "senior" for the experience needed
"""
SENIOR_LABELS = ["senior"]

"""
These labels are going count as "enhancement" for the issue type
"""
ENHANCEMENT_LABELS = ["enhancement"]

"""
These labels are going count as "bugfix" for the issue type
"""
BUGFIX_LABELS = ["bug", "bugfix"]


class Project(models.Model):
    # TODO: Parse issue_count from related issues
    # TODO: Parse tags from related issues
    """
    Project model is used to store the projects that are
    included in the project.
    """
    created = models.DateTimeField(auto_now_add=True)
    display_name = models.TextField(max_length=150)
    first_color = models.CharField(max_length=14, default="#FF614C")
    second_color = models.CharField(max_length=14, blank=True)
    description = models.TextField()
    logo = models.URLField()
    link = models.URLField(unique=True)
    setup_duration = models.CharField(max_length=100, blank=True)
    tags = ListTextField(base_field=models.CharField(max_length=80), size=100)
    issues_count = models.IntegerField()

    class Meta:
        ordering = ('created',)  # Ascending order according to date created.

    def __str__(self):
        return self.display_name


class Repository(models.Model):
    """
    Repository model is used to store the repository url.
    """
    repository_url = models.URLField(unique=True)
    created = models.DateTimeField(auto_now_add=True)
    project = models.ForeignKey(
        Project, on_delete=models.SET_NULL, null=True, blank=True)

    class Meta:
        ordering = ('created',)  # Ascending order according to date created.

    def __str__(self):
        return self.repository_url


class IssueLabel(models.Model):
    """
    Label model for storing labels of an issue.
    """
    label_id = models.IntegerField(primary_key=True)
    label_url = models.URLField()
    label_name = models.CharField(max_length=100)
    label_color = models.CharField(max_length=6)

    class Meta:
        ordering = ('label_name',)  # Ascending order according to label_name.

    def __str__(self):
        return self.label_name


class Issue(models.Model):
    """
    Issue model is used to store github issues.
    """
    # Setting choices for experience needed to solve an issue.
    EASYFIX = 'easy'
    MODERATE = 'moderate'
    SENIOR = 'senior'
    EXPERIENCE_NEEDED_CHOICES = (
        (EASYFIX, 'easy'),
        (MODERATE, 'moderate'),
        (SENIOR, 'senior'),
    )
    # Model attributes start from here.
    issue_id = models.IntegerField(primary_key=True)
    title = models.TextField()
    experience_needed = models.CharField(
        max_length=50,
        choices=EXPERIENCE_NEEDED_CHOICES,
        default=MODERATE,
    )
    expected_time = models.CharField(max_length=100)
    language = models.TextField()
    tech_stack = models.TextField()
    created_at = models.DateTimeField()
    updated_at = models.DateTimeField()
    issue_number = models.IntegerField()
    issue_labels = models.ManyToManyField(IssueLabel, blank=True)
    issue_url = models.URLField()
    issue_body = models.TextField()
    issue_type = models.TextField(default="")
    repo = models.ForeignKey(
        Repository, on_delete=models.CASCADE, null=True)
    project = models.ForeignKey(
        Project, on_delete=models.SET_NULL, null=True, blank=True)

    class Meta:
        ordering = ('updated_at',)  # Ascending order according to updated_at.

    def __str__(self):
        return self.title


@periodic_task(run_every=timedelta(minutes=ISSUE_UPDATE_PERIOD), name="periodic_issues_updater")
def periodic_issues_updater():
    """
    Update `Issue` model in the database in every
    `ISSUE_UPDATE_PERIOD` minutes.
    """
    list_of_repos = Repository.objects.values('repository_url',)
    for repo in list_of_repos:
        repository, username = get_repo_and_username(repo['repository_url'])
        language_request = request_to_github(repository, username, "")
        if language_request['error']:
            logging.error("Error" + str(language_request['data']))
            continue
        language = language_request['data']['language']
        issue_list = request_to_github(
            repository, username, '/issues?state=all')
        if issue_list['error']:
            logging.error("Error" + str(issue_list['data']))
            continue
        for issue in issue_list['data']:
            validate_and_store_issue(issue, language, repo)


def validate_and_store_issue(issue, language, repo):
    """
    Validate issue:- if valid - store it into database,
    else - Do not store in database
    """
    if is_issue_state_open(issue):
        repo_instace = Repository.objects.get(
            repository_url=repo.get("repository_url"))
        project = repo_instace.project
        store_issue_in_db(
            issue, language, project, repo_instace)


def is_issue_state_open(issue):
    """
    Returns true if issue state is open else
    return false and delete the issue from database.
    """
    if issue.get('pull_request'):
        return False

    if issue['state'] == 'open':
        return True
    else:
        delete_closed_issues(issue)  # Delete closed issues from db.
        return False


# FIXME: Rewrite and cleanup this


def store_issue_in_db(issue, language, project, repo):
    """Stores issue in db"""
    # experience_needed, bodyLanguage, expected_time, technology_stack = parse_issue(
    #     issue['body'])
    # experience_needed = parse_experience(experience_needed.strip().lower())
    # expected_time = expected_time.strip().lower()
    # technology_stack = technology_stack.strip().lower()
    try:
        language = language.strip().lower()
        issue_type = ""
        issue_instance = Issue(issue_id=issue['id'], title=issue['title'], language=language, created_at=issue['created_at'],
                               updated_at=issue['updated_at'], issue_number=issue['number'], issue_url=issue['html_url'], issue_body=issue['body'], issue_type=issue_type, project=project, repo=repo)
        issue_instance.save()
        for label in issue['labels']:
            try:
                name = label['name'].lower()
                if name in ENHANCEMENT_LABELS + BUGFIX_LABELS:
                    issue_instance.issue_type = parse_type(name)
                    issue_instance.save()
                elif name in SENIOR_LABELS + EASY_LABELS + MODERATE_LABELS:
                    issue_instance.experience_needed = parse_experience(name)
                    issue_instance.save()

            except Exception as e:
                logging.error("Couldn't parse label: " + label)
                logging.error(e)
            label_instance = IssueLabel(label_id=label['id'], label_name=label['name'],
                                        label_url=label['url'], label_color=label['color'])
            label_instance.save()
            issue_instance.issue_labels.add(label_instance)
    except Exception as e:
        logging.error("Couldn't parse issue: ")
        logging.error(e)


def delete_closed_issues(issue):
    """Delete issues that are closed on GitHub but present in our db"""
    try:
        issue_instance = Issue.objects.get(issue_id=issue['id'])
        issue_instance.delete()
    except Exception:
        logging.warning('Closed issue with id ' +
                        str(issue['id']) + ' is not present is database.')


def parse_issue(issue_body):
    """
    Parse the issue body and return `experience_needed`, `language`,
    `expected_time` and `technology_stack`.
    """
    issue_body = issue_body.lower()
    experience_needed = find_between(issue_body, 'experience', '\r\n')
    language = find_between(issue_body, 'language', '\r\n')
    expected_time = find_between(issue_body, 'expected-time', '\r\n')
    technology_stack = find_between(issue_body, 'technology-stack', '\r\n')
    return experience_needed, language, expected_time, technology_stack


def find_between(string, first, last):
    """
    Return string between two substrings `first` and `last`.
    """
    try:
        start = string.index(first) + len(first)
        end = string.index(last, start)
        return string[start:end].replace(': ', '')
    except ValueError:
        return ""


def parse_experience(experience_needed):
    if experience_needed in EASY_LABELS:
        return "easy"
    elif experience_needed in MODERATE_LABELS:
        return "moderate"
    elif experience_needed in SENIOR_LABELS:
        return "senior"
    return experience_needed


def parse_type(issue_type):
    if issue_type in BUGFIX_LABELS:
        return "bugfix"
    elif issue_type in ENHANCEMENT_LABELS:
        return "enhancement"
    return issue_type
