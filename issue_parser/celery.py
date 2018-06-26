from __future__ import absolute_import
import os
from celery import Celery
from django.conf import settings

# set the default Django settings module for the 'celery' program.
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'issue_parser.settings')
broker = "amqp://" + os.environ.get('RABBIT_USERNAME') +\
    ":" + os.environ.get('RABBIT_PASSWORD') + "@" + \
    os.environ.get('RABBIT_HOST')+":5672"

app = Celery('issue_parser', broker=broker)

# Using a string here means the worker will not have to
# pickle the object when using Windows.
app.config_from_object('django.conf:settings')
app.autodiscover_tasks(lambda: settings.INSTALLED_APPS)


@app.task(bind=True)
def debug_task(self):
    print('Request: {0!r}'.format(self.request))
