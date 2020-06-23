package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/owenjacob/hubblegraphql/graph"
	"github.com/owenjacob/hubblegraphql/graph/generated"
	"github.com/rs/cors"
)

// Defining the Playground handler
func playgroundHandler() httprouter.Handle {
	h := playground.Handler("GraphQL", "/query")

	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, req)
	}
}

// Defining the Graphql handler
func graphqlHandler(dbConnection *gorm.DB) httprouter.Handle {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: dbConnection}}))

	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, req)
	}
}

// Implementation with Httprouter from: https://qiita.com/tlta-bkhn/items/3f8883a7db059aa58717
func main() {
	dbConnection, err := gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port=5432 user="+os.Getenv("DB_USERNAME")+" dbname="+os.Getenv("DB_DATABASE")+" password="+os.Getenv("DB_PASSWORD")+" sslmode=disable")
	if err != nil {
		panic(err)
	}
	dbConnection.LogMode(true)
	defer dbConnection.Close()

	mux := httprouter.New()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})

	mux.GET("/", playgroundHandler())
	mux.POST("/query", graphqlHandler(dbConnection))

	router := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":3000", router))
}
