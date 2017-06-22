# Deploying the Back-End
> This document only contains how to deploy the back-end part of the app for running full app you have to deploy front end client also. You can find the front end client documentation [here](https://github.com/mozillacampusclubs/issue_parser_frontend/)

So lets start with the basic steps on how to deploy this back-end API.
- Install Python 2.7 and pip on the system.
- Open a terminal
- run `pip install -r requirements.txt`
- **RELEX . . .**
- run `python manage.py makemigration` for making migrations.
- run `python manage.py migrate` for migrating database.
- run `python manage.py createsuperuser` to create a login password for logging in to admin panel.
- For starting dev-server run `python manage.py runserver`.
- Now you also have to setup worker server (alongside main server) for fetching github issues periodically. For this follow these steps:
    - Open another terminal and run `celery -A issue_parser beat -l info`.
    - Open one more terminal and run `celery -A issue_parser worker -l info`
- Now fire up a browser and go to `/admin/` for admin view and `/issues/` and `/metadata/` for issues and metadata view.

**Admin view**: Add repositories to the system, which you want to use. The system will fetch the issues of these repositories. To add repositories head to `/admin/`.


