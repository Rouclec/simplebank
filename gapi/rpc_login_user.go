package gapi

import (
	"context"
	"errors"

	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/pb"
	"github.com/rouclec/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User with username %s not found", req.GetUsername())
		}
		return nil, status.Errorf(codes.Internal, "An internal error occured: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.Password)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication error: %s", err)
	}

	accessToken, _, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "An internal error occured: %s", err)
	}

	response := &pb.LoginUserResponse{
		AccessToken: accessToken,
		User:        bindGRPCUserResponse(user),
	}

	return response, nil
}
