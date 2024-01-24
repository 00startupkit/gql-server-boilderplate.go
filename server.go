package main

import (
	"fmt"
	"go-graphql-api/graph"
	oauth "go-graphql-api/oauth2"
	"go-graphql-api/util"
	"go-graphql-api/util/gql_middleware"
	"go-graphql-api/util/logger"
	"net/http"

	"go-graphql-api/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := setup_environment()
	if err != nil {
		panic(fmt.Errorf("failed to setup environment: %v", err))
	}

	db, err := database.GetDbInstance()
	if err != nil {
		panic(fmt.Errorf("failed to instantiate database connection: %v", err))
	}

	router := chi.NewRouter()
	router.Use(gql_middleware.JwtAuthMiddleware())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{
			Database: db,
		}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	oauth.RegisterOauthRoutes(router)

	logger.Info("connect to %s/ for GraphQL playground", util.ServerUri())
	err = http.ListenAndServe(":"+util.ServerPort(), router)
	if err != nil {
		panic(err)
	}
}

func setup_environment() error {
	return godotenv.Load()
}
