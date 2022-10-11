package gapi

import (
	"fmt"

	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	"github.com/NguyenMinhKhanhBK/simple_bank/pb"
	"github.com/NguyenMinhKhanhBK/simple_bank/token"
	"github.com/NguyenMinhKhanhBK/simple_bank/util"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
	pb.UnimplementedSimpleBankServer
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}
