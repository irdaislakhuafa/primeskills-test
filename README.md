# Todo App
A simple todo list app ever, no more forgot about your tasks.

## Project Structures
```bash
.
├── deploy/ -- Dockerfiles configurations for multiple environment.
├── docker-compose.yaml -- To setup dev env and build docker images.  
├── docs/ -- Store assets and docs for app.
│   ├── sql/ -- Store sql files like schema and queries for app.
│   └── templates/ -- For static file templates like email templates and others.
├── etc/ -- Other configuration for app like app and env configurations for multiple server env.
│   └── cfg/ -- Configurations dir
├── go.mod
├── go.sum
├── README.md
├── sqlc.yaml -- SQLC configuration about schema and queries database that used in app.
└── src/ -- Source code of app.
    ├── business/ -- Business layer of app. Contains any codes that related with business logic.
    │   ├── domain/ -- Layer to communicating with database.
    │   └── usecase/ -- Layer to make business logic of app based on business needs. 
    ├── cmd/
    │   └── main.go -- Entrypoint file.
    ├── connection/ -- Code to setup database connection.
    ├── entity/ -- Store all models or entities that used on app.
    ├── handler/ -- Layer to create handler of endpoints. It can be Rest/GraphQL/gRPC or others.
    │   └── rest/ -- Rest API layer to setup endpoints and create Rest API handlers.
    ├── utils/ -- Store about utilities code for app.
    │   ├── config/ -- Code to parse configuration files to Go struct.
    │   ├── ctxkey/ -- Context keys collections.
    │   ├── entutils/ -- Entity utilities to store related code that use on models or entities.
    │   └── mailtemplates/ -- Collection codes to read and parse email templates.
    └── validation/ -- Collection of validation codes to validate request before processing the data.
```

## Deployment

- `ci/cd`
This repo using `github actions` as `ci/cd` that configured at directory `.github/workflows/`. This automation will triggered by `github release tags`, so if application has new version released then that version will deployed to server automatically. No more repetitive task deployment.

- `local`
First, clone this code to your local env.
```bash
$ git clone https://github.com/irdaislakhuafa/primeskills-test.git
$ cd primeskills-test/
```

You can run this code on local environment with `docker` and `docker compose`. If your don't have `docker` installed on your machine, i recommend you to install it first here [docker install](https://docs.docker.com/get-started/get-docker/).

After you install `docker` then you just need type this command below to run this app without any setups again.
```bash
$ docker composer up -d
```
