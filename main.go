package main

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"net"
	"net/http"
	"sync"

	_ "embed"

	"github.com/NguyenMinhKhanhBK/simple_bank/api"
	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	"github.com/NguyenMinhKhanhBK/simple_bank/gapi"
	"github.com/NguyenMinhKhanhBK/simple_bank/pb"
	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed doc/swagger
var docFS embed.FS
var swaggerFS, _ = fs.Sub(docFS, "doc/swagger")

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

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go runGRPCServer(wg, config, store)
	go runGatewayServer(wg, config, store)
	wg.Wait()

	logrus.Info("Exiting ...")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		logrus.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logrus.Fatal("cannot start server:", err)
	}
}

func runGRPCServer(wg *sync.WaitGroup, config util.Config, store db.Store) {
	defer wg.Done()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		logrus.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		logrus.Fatal("cannot create listener:", err)
	}

	logrus.Infof("start gRPC server at %v", config.GRPCServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		logrus.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(wg *sync.WaitGroup, config util.Config, store db.Store) {
	defer wg.Done()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		logrus.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		logrus.Fatal("cannot register handler server:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fsHandler := http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFS)))
	mux.Handle("/swagger/", fsHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		logrus.Fatal("cannot create listener:", err)
	}

	logrus.Infof("start HTTP gateway server at %v", config.HTTPServerAddress)

	err = http.Serve(listener, mux)
	if err != nil {
		logrus.Fatal("cannot start gRPC server:", err)
	}
}
