package main

import (
	"blog/repo/postgres"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// const (
// 	host     = "/cloudsql/lunar-outlet-334614:us-central1:blog"
// 	port     = 5432
// 	user     = "bemihai"
// 	password = "bemihai"
// 	dbname   = "blog"
// 	schema   = "blog"
// )

const (
	DBDriver = "postgres"
	DBSource = "postgresql://postgres:postgres@localhost:5432/blog?sslmode=disable"
)

func main() {

	// get config
	// config, err := util.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }

	// define handler for http requests with postgres repository
	database := dbConnect(DBDriver, DBSource)
	defer database.Close()

	handler := BlogServer{
		Service: &postgres.PSQLRepository{DB: database},
	}

	// define the associations between endpoints and handlers
	r := mux.NewRouter()

	// define handler for GET on "/"
	r.HandleFunc("/", CheckHealth)

	// define handler for GET on "/articles" endpoint
	r.Handle("/articles", http.HandlerFunc(handler.ListArticles)).Methods(http.MethodGet)

	// define handler for GET on "/articles/id" endpoint
	r.Handle("/articles/{id}", http.HandlerFunc(handler.GetArticleById)).Methods(http.MethodGet)

	// define handler for POST on "/articles" endpoint
	r.Handle("/articles", http.HandlerFunc(handler.AddArticle)).Methods(http.MethodPost)

	// define handler for DELETE on "/articles/id" endpoint
	r.Handle("/articles/{id}", http.HandlerFunc(handler.DeleteArticleById)).Methods(http.MethodDelete)

	// define handler for DELETE on "/authors" endpoint
	r.Handle("/authors", http.HandlerFunc(handler.DeleteAuthorByNameAndEmail)).Methods(http.MethodDelete)

	// define handler for not found endpoint
	r.NotFoundHandler = http.NotFoundHandler()

	// defines handler for not allowed methods
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowed)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port: %s", port)

	//  Start HTTP
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), r); err != nil {
		log.Fatal("Failed starting http server: ", err)
	}

}

// dbConnect creates a connection to a database.
func dbConnect(dbDriver string, dbSource string) *sql.DB {

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to db!")
	return db
}
