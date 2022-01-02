package main

import (
	api "blog/api"
	"blog/util"
	"blog/util/utildb"
	"log"
)

func main() {

	// get config
	config := util.LoadConfig("")

	// open connection to db
	db, err := utildb.Connect(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("Failed to connect to the db: ", err)
	}
	defer db.Close()

	// start server
	api.StartNewServer(db)

}
