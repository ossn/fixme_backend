# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from datetime import timedelta
from django.db import models
from django_mysql.models import ListTextField
from core.utils.services import request_github_issues

from celery.decorators import periodic_task, task

ISSUE_UPDATE_PERIOD = 15  # in minutes


class Project(models.Model):
    # TODO: Parse issue_count from related issues
    # TODO: Parse tags from related issues
    """
    Project model is used to store the projects that are
    included in the project.
    """
    created = models.DateTimeField(auto_now_add=True)
    display_name = models.CharField(max_length=100)
    first_color = models.CharField(max_length=14, default="#FF614C")
    second_color = models.CharField(max_length=14, blank=True)
    description = models.TextField()
    logo = models.URLField()
    link = models.URLField()
    setup_duration = models.CharField(max_length=100, blank=True)
    tags = ListTextField(base_field=models.CharField(max_length=50), size=100)
    issues_count = models.IntegerField()

    class Meta:
        ordering = ('created',)  # Ascending order according to date created.
        unique_together = ("link", "display_name")  # Avoid repo duplicates.

    def __str__(self):
        return self.display_name


class UserRepo(models.Model):
    """
    UserRepo model is used to store the username and repo-name
    for a repository.
    """
    user = models.CharField(max_length=100)
    repo = models.CharField(max_length=100)
    created = models.DateTimeField(auto_now_add=True)
    project = models.ForeignKey(
        Project, on_delete=models.SET_NULL, null=True, blank=True)

    class Meta:
        ordering = ('created',)  # Ascending order according to date created.
        unique_together = ("user", "repo")  # Avoid repo duplicates.

    def __str__(self):
        return '/%s/%s' % (self.user, self.repo)


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
    issue_type = models.CharField(max_length=100, default="")

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
    list_of_repos = UserRepo.objects.values('user', 'repo',)
    for repo in list_of_repos:
        issue_list = request_github_issues(repo['user'], repo['repo'])
        if issue_list['error']:
            print "Error" + str(issue_list['data'])
        else:
            for issue in issue_list['data']:
                validate_and_store_issue(issue)


def validate_and_store_issue(issue):
    """
    Validate issue:- if valid - store it into database,
    else - Do not store in database
    """
    if is_issue_state_open(issue):
        if is_issue_valid(issue):
            store_issue_in_db(issue)


def is_issue_state_open(issue):
    """
    Returns true if issue state is open else
    return false and delete the issue from database.
    """
    if issue['state'] == 'open':
        return True
    else:
        delete_closed_issues(issue)  # Delete closed issues from db.
        return False


def is_issue_valid(issue):
    """
    Checks if issue is valid for system or not.
    Return True if valid else return false.
    """
    parsed = parse_issue(issue['body'])
    for item in parsed:
        if not item:
            return False  # issue is not valid
    print 'Issue with id ' + str(issue['id']) + ' is not valid for our system.'
    return True  # issue is valid


def store_issue_in_db(issue):
    """Stores issue in db"""
    experience_needed, language, expected_time, technology_stack = parse_issue(
        issue['body'])
    experience_needed = experience_needed.strip().lower()
    if experience_needed == "easyfix":
        experience_needed = "easy"
    language = language.strip().lower()
    expected_time = expected_time.strip().lower()
    technology_stack = technology_stack.strip().lower()
    issue_type = ""
    issue_instance = Issue(issue_id=issue['id'], title=issue['title'],
                           experience_needed=experience_needed, expected_time=expected_time,
                           language=language, tech_stack=technology_stack,
                           created_at=issue['created_at'], updated_at=issue['updated_at'],
                           issue_number=issue['number'], issue_url=issue['html_url'],
                           issue_body=issue['body'], issue_type=issue_type)
    issue_instance.save()
    for label in issue['labels']:
        try:
            if label['name'].lower() in ['enhancement', 'bugfix', 'task']:
                issue_instance.issue_type = label['name'].lower()
                issue_instance.save()
        except:
            print 'Couldn\'t parse label: ' + label
        label_instance = IssueLabel(label_id=label['id'], label_name=label['name'],
                                    label_url=label['url'], label_color=label['color'])
        label_instance.save()
        issue_instance.issue_labels.add(label_instance)


def delete_closed_issues(issue):
    """Delete issues that are closed on GitHub but present in our db"""
    try:
        issue_instance = Issue.objects.get(issue_id=issue['id'])
        issue_instance.delete()
    except Exception:
        print 'Closed issue with id ' + str(issue['id']) + ' is not present is database.'


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
