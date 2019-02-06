package main

import (
	log "log"
	http "net/http"
	os "os"

	"github.com/gorilla/mux"

	"github.com/aneri/graphql-dataloaden/dataloader"

	handler "github.com/99designs/gqlgen/handler"
	"github.com/aneri/gqlgen-dataloader/middleware"
	resolever "github.com/aneri/graphql-dataloaden/api/handler"
	graph "github.com/aneri/graphql-dataloaden/graph"
)

const defaultPort = "8000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := mux.NewRouter()
	queryHandler := handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &resolever.Resolver{}}))
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", middleware.MultipleMiddleware(queryHandler, resolever.DbMiddleware, dataloader.LoadMiddleware))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
