# FixMe Backend

> This app is built with [Buffalo](https://gobuffalo.io/en)

## Running the app

- Install docker-compose
- Generate a github token, there is a guide [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
- Open the docker-compose.yml file and add values to the env variable
- Run `docker-compose up -d`

## Dev enviroment

> Note: It's recommended to read the getting started guide and a few things regarding buffalo from [here](https://gobuffalo.io/en/docs/installation)

### Using docker-compose

- Install docker-compose
- Generate a github token, there is a guide [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
- Open the dev.docker-compose.yml file and add values to the env variable
- Run `docker-compose -f dev.docker-compose.yml up --build` (Note: every time that you save a file the app will be recompiled and restarted)
- The app should be up and running at http://localhost:3000

### Installing locally

#### Requirments

- [Postgresql](https://www.postgresql.org/download/)
- [Go](https://golang.org/doc/install)
- [Buffalo](https://gobuffalo.io/en/docs/installation)

#### Starting the app

- Configure `database.yml` file in order to connect the app to your PostgreSQL instance
- Generate a github token, there is a guide [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)
- Set your github token to an environment variable called `GITHUB_TOKEN`
- Set a random jwt secret key to an environment variable called `JWT_SECRET`
- Run `buffalo db create -a`
- Run `buffalo db migrate`
- Run `buffalo task db:seed`
- Run `buffalo dev`(Note: This will watch the current directory and it will recompile and restart the app every time there is a change in your files)
- The app should be up and running at http://localhost:3000
