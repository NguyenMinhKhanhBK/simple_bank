package main

import (
	"database/sql"

	"github.com/NguyenMinhKhanhBK/simple_bank/api"
	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:root@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		logrus.Fatal("cannot connect to db:", err)
	}

	if err := conn.Ping(); err != nil {
		logrus.Fatal("cannot ping db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		logrus.Fatal("cannot start server:", err)
	}
}
