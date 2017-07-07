# issue_parser_backend

[![Build Status](https://travis-ci.org/razat249/issue_parser_backend.svg?branch=master)](https://travis-ci.org/razat249/issue_parser_backend)
[![Issues](https://camo.githubusercontent.com/926d8ca67df15de5bd1abac234c0603d94f66c00/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f636f6e747269627574696f6e732d77656c636f6d652d627269676874677265656e2e7376673f7374796c653d666c6174)](https://github.com/razat249/issue_parser_backend/issues)

Backend API of issue parser app. This repository is required to successfully run the app you can find the front end client [here](https://github.com/mozillacampusclubs/issue_parser_frontend/).


## How to deploy
Find the full deployment guide [here](./docs/deploy.md).
> To successfully run full app, you have to deploy the front end client also. You can find the front end client documentation [here](https://github.com/mozillacampusclubs/issue_parser_frontend/)


## How to Contribute
You can report bugs at the [issue tracker](https://github.com/mozillacampusclubs/issue_parser_backend/issues). Find full guide on how to contribute [here](./docs/contributing.md)

> Help us by creating Pull Requests and solving [issues](https://github.com/mozillacampusclubs/issue_parser_backend/issues).

> For more in-depth knowledge of system read below `How it works`.

## How it works
This image of the system below will give you a rough overview of how api works.
![System Design](./docs/img/system-design.png)

There will be several components in the backend:
1. **Fetcher** : This component will fetch the data from the github API (issues in our case). This component will do only one thing i.e. fetching the data. I suggest to keep the frequency of data fetching in every 15 mins.
2. **Parser/Model** : This component will do 3 things:
Checking if the issue is passing all the quality and documentation criteria.
Parsing the data according to the schema we have for the issues.
And saving the data to database.
3. **Data** : This is database containing two table Repositories and Issues.
4. **Controller** : This is a thin layer between Model and front end. This layer will handle all the routing and REST APIs stuff.
5. **Admin View** : Admin view can be used to add or remove repos from the system.

<!--## License
To Do - discuss with mentor.-->
