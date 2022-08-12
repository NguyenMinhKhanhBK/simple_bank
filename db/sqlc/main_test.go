package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDirver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDirver, dbSource)
	if err != nil {
		logrus.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
