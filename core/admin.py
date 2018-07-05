# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.contrib import admin

# Register/Unregister models here.
from core.models import Repository, Issue, IssueLabel, Project
from django.contrib.auth.models import *

admin.site.unregister(Group)
admin.site.unregister(User)
admin.site.register(Repository)
admin.site.register(Project)
admin.site.register(Issue)
