#!/bin/bash

sleep 5

echo "Check db"
python manage.py check

echo "Prepaire assets"
python manage.py collectstatic --noinput

# Make migrations
echo "Make migrations file"
python manage.py makemigrations --noinput
python manage.py makemigrations core --noinput

# Apply database migrations
echo "Apply database migrations"
python manage.py migrate

# Start server
echo "Starting server"
gunicorn issue_parser.wsgi -c ./config/gunicorn.py  -b 0.0.0.0:8000
