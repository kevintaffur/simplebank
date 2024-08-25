package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kvgtl/simplebank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Panic("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
