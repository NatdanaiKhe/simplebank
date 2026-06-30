package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/NatdanaiKhe/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	testDB, err = sql.Open(config.DB_Driver, config.DB_URL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQueries = New(testDB)
	code := m.Run()
	os.Exit(code)
}
