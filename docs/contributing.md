# How to Contribute
> Help us by creating Pull Requests and solving [issues](https://github.com/mozillacampusclubs/issue_parser_backend/issues).

**For setting up development environment follow these steps:**
- Install Python 2.7 and pip on your system.
- Open a terminal
- run `sudo apt-get update`
- run `sudo apt-get install python-pip python-dev mysql-server libmysqlclient-dev`
- Install virtualenv using cmd `pip install virtualenv`.
- Clone this repo.
- cd into the repo.
- Create virtual env by running `virtualenv env`.
- run `source env/bin/activate`.
- run `pip install -r requirements.txt`
- **RELEX . . .**
- run `python manage.py makemigration` for making migrations.
- run `python manage.py migrate` for migrating database.
- run `python manage.py createsuperuser` to create a login password for loging in to admin panel.
- For starting dev-server run `python manage.py runserver`.
- **Admin view**: To add repositories head to `/admin/`. Add repositories to the system, which you want to use. The system will fetch the issues of these repositories. you have to fill the username and repo name. Some user-repos that support the system are given below in the `Supported Repositories` section below, use these repos to start off.
- Now you also have to setup worker server (alongside main server) for fetching github issues periodically (15 mins). For this follow these steps:
    - Open another terminal and run `celery -A issue_parser beat -l info`.
    - Open one more terminal and run `celery -A issue_parser worker -l info`
- Now fire up a browser and go to `/issues/` for seeing the json data of issues.
> You have to add repositories to the system or the app will not fetch issues.

## Supported Repositories
Use below data to add repos to the system via admin view.
1. **user:** `mozillacampusclubs`, **repo:** `issue_parser_backend`
2. **user:** `mozillacampusclubs`, **repo:** `issue_parser_frontend`
3. **user:** `razat249`, **repo:** `github-view`

**Issues should follow this template to be valid for the system:**
```
Experience: Easyfix/Moderate/Senior
Expected-time: 1 week/2 months/etc.
Language: python/Javascript/others
Technology-stack: Django/React.js/others

## Description
Write your description here.
```
