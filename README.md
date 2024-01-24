<p align="center">
  <img src="docs/logo-1.png" alt="quickapi logo" height="35" />
</p>
<p align="center">
  <img src="https://github.com/00startupkit/gql-server-boilderplate.go/actions/workflows/build.yml/badge.svg" alt="ci status" />
</p>

This is boilerplate for setting up a GraphQL api server backed by a
mysql database.

# Setting up Environment
Environment will be automatically loaded by `joho/dotenv` at the start of the
application. Here is an example `.env` that should work for the current config:

```.env
MYSQL_USER=root
MYSQL_PASS=password
MYSQL_HOST=localhost:3306
MYSQL_IDLE_CONNECTIONS=10
MYSQL_OPEN_CONNECTIONS=100
SERVER_PORT=8090
```

# Starting the Server
The project is configured with *[cosmtrek/air](https://github.com/cosmtrek/air)* to hot reload. The config is located in `.air.toml`. After downloading the  *air* executable with `go install github.com/cosmtrek/air@latest`, the hot-reloadable server can be started by running `air`.

# Adding Database Models
The database models are defined in `dbmodel/db_model.go`. To add a model to be auto-migrated on startup, define the model struct and add it to the `Models` variable found in that file.
