package main

import (
	"database/sql"

	"github.com/NguyenMinhKhanhBK/simple_bank/api"
	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		logrus.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logrus.Fatal("cannot connect to db:", err)
	}

	if err := conn.Ping(); err != nil {
		logrus.Fatal("cannot ping db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		logrus.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		logrus.Fatal("cannot start server:", err)
	}
}
