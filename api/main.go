package main

import (
	"blog/repo/postgres"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "blog"
	schema   = "blog"
)

func main() {

	// define handler for http requests with postgres repository
	database := psqlConnect(host, port, user, password, dbname, schema)
	defer database.Close()

	handler := BlogServer{
		Service: &postgres.PSQLRepository{DB: database},
	}

	// define the associations between endpoints and handlers
	router := mux.NewRouter()

	// define handler for GET on "/articles" endpoint
	router.Handle("/articles", http.HandlerFunc(handler.ListArticles)).Methods(http.MethodGet)

	// define handler for GET on "/articles/id" endpoint
	router.Handle("/articles/{id}", http.HandlerFunc(handler.GetArticleById)).Methods(http.MethodGet)

	// define handler for POST on "/articles" endpoint
	router.Handle("/articles", http.HandlerFunc(handler.AddArticle)).Methods(http.MethodPost)

	// define handler for DELETE on "/articles/id" endpoint
	router.Handle("/articles/{id}", http.HandlerFunc(handler.DeleteArticleById)).Methods(http.MethodDelete)

	// define handler for DELETE on "/authors" endpoint
	router.Handle("/authors", http.HandlerFunc(handler.DeleteAuthorByNameAndEmail)).Methods(http.MethodDelete)

	// define handler for not found endpoint
	router.NotFoundHandler = http.NotFoundHandler()

	// defines handler for not allowed methods
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowed)

	// defines the server instance by specifing the endpoints handler and the address (host:port)
	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting the server...listening on port 8000")

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Connects to a postgres database.
func psqlConnect(host string, port int, user string, password string, dbname string, schema string) *sql.DB {

	psqlConString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbname, schema)

	db, err := sql.Open("postgres", psqlConString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres!")
	return db
}
