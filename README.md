# issue_parser_backend
Backend api for issue parser app.
## Contributing
Follow below steps for building the project.
- Install Python 2.7 and pip on your system.
- Open a terminal
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
- Now fire up a browser and go to `http://localhost:8000/admin/` for admin view and `http://localhost:8000/issues/` and `http://localhost:8000/metadata/` for issues and metadata view.
