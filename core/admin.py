# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.contrib import admin

# Register/Unregister models here.
from core.models import UserRepo
from django.contrib.auth.models import *

admin.site.unregister(Group)
admin.site.unregister(User)
admin.site.register(UserRepo)
