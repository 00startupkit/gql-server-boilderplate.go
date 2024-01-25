<p align="center">
  <img src="docs/logo-1.png" alt="quickapi logo" height="35" />
</p>
<p align="center">
  <img src="https://github.com/00startupkit/gql-server-boilderplate.go/actions/workflows/build.yml/badge.svg" alt="ci status" />
</p>

This is boilerplate for setting up a GraphQL api server backed by a
mysql database.

# Install
```
git clone https://github.com/00startupkit/gql-server-boilderplate.go.git
go get ./...
go install github.com/cosmtrek/air@latest
# The air binary might be stored in ~/go/bin so make sure its in your path.
```

## Prerequisites
This server configuration assumes a MySQL instance is configured and running. Once you have your MySQL instance running, configure the connection setting in the environment variables in the [steps below](https://github.com/00startupkit/gql-server-boilderplate.go?tab=readme-ov-file#setting-up-environment).

# Setting up Environment
Environment will be automatically loaded by `joho/dotenv` at the start of the
application. Here is an example `.env` that should work for the current config:

```.env
MYSQL_USER=root
MYSQL_PASS=password
MYSQL_HOST=localhost:3306
MYSQL_IDLE_CONNECTIONS=10
MYSQL_OPEN_CONNECTIONS=100
SERVER_HOST=http://localhost
SERVER_PORT=8090
JWT_SECRET=yourtokensecret
```

## Configuring OAuth2
OAuth2 settings can be configured from `oauth2/config.go`. *Google* is defined there by default. For it to work, set the proper `GOOGLE_CLIENT_*` environment variables. Extend the list to define multiple OAuth2 providers. Make sure to also implement the conversion from the user payload from the provider to the *user* model that will be stored in the database.

# Starting the Server
The project is configured with *[cosmtrek/air](https://github.com/cosmtrek/air)* to hot reload. The config is located in `.air.toml`. After downloading the  *air* executable with `go install github.com/cosmtrek/air@latest`, the hot-reloadable server can be started by running `air`.

# Adding Database Models
The database models are defined in `dbmodel/db_model.go`. To add a model to be auto-migrated on startup, define the model struct and add it to the `Models` variable found in that file.

# TODO
- Update chi dependency from deprecated version 1.5.5.
- Enable cors for the server.
- Store the auth token in the response's cookies on oauth success.
