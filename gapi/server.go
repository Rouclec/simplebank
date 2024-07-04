package gapi

import (
	"fmt"

	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/pb"
	"github.com/rouclec/simplebank/token"
	"github.com/rouclec/simplebank/util"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

// Creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) //Switch between NewPasetoMaker and NewJWTMaker to use either Paseto or JWT tokens respectively
	if err != nil {
		return nil, fmt.Errorf("error creating token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}


	return server, nil
}