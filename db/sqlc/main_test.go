package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		logrus.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logrus.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
