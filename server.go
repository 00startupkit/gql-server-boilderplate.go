package main

import (
	"fmt"
	"go-graphql-api/graph"
	"go-graphql-api/util"
	"log"
	"net/http"

	"go-graphql-api/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := setup_environment()
	if err != nil {
		panic(fmt.Errorf("failed to setup environment: %v", err))
	}

	port := util.EnvOrDefault("SERVER_PORT", defaultPort)

	err = database.ConnectDB()
	if err != nil {
		panic(fmt.Errorf("failed to instantiate database connection: %v", err))
	}
	database.CreateDB()
	database.MigrateDB()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Database: database.DBInstance,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setup_environment() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
