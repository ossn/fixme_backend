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
