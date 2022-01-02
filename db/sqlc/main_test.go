package sqlc_db

import (
	"blog/util"
	"blog/util/utildb"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

// TestMain is the entry point of all the tests in the package.
func TestMain(m *testing.M) {

	config := util.LoadConfig("../../dev_config.json")

	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = utildb.ExecFile(conn, "./db/blog_migrate.sql")
	if err != nil {
		log.Fatal("could not create tables:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
