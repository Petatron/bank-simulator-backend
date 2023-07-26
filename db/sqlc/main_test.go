package db

import (
	"database/sql"
	"github.com/Petatron/bank-simulator-backend/db/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Unable to load project config with error: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to Database with error: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
