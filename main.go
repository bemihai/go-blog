package main

import (
	api "blog/api"
	"blog/util/utildb"
	"log"
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
	database, err := utildb.Connect(DBDriver, DBSource)
	if err != nil {
		log.Fatal("Failed to connect to the db: ", err)
	}
	defer database.Close()

	api.StartNewServer(database)

}
