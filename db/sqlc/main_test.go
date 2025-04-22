package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/olaniyi38/BE/util"
)

var testQueries *Queries
var testDb *sql.DB
var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}

	testDb = conn
	defer conn.Close()
	testQueries = New(testDb)
	testStore = NewStore(testDb)

	os.Exit(m.Run())
}
