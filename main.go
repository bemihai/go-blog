package main

import (
	api "blog/api"
	"blog/util"
	"blog/util/utildb"
	"log"
)

func main() {

	// get config
	config := util.LoadConfig("prod_config.json")
	source := util.GetDBSource(config)

	log.Println(source)

	// open connection to db
	db, err := utildb.Connect(config.DB_DRIVER, source)
	if err != nil {
		log.Fatal("Start app: Failed to connect to the db: ", err)
	}
	defer db.Close()

	// start server
	api.StartNewServer(db)

}
