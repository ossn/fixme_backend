FROM python:2.7-alpine
RUN apk add --update --no-cache --virtual .build-deps \
  gcc \
  make \
  libc-dev \
  musl-dev \
  linux-headers \
  pcre-dev \
  mysql-dev
WORKDIR /code
ADD  requirements.txt .
RUN pip install -r requirements.txt
ADD . .
RUN mkdir  -p /code/static
RUN chown -R 1000:1000 /code
USER 1000
