package gapi

import (
	"context"

	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/pb"
	"github.com/rouclec/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error hashing password: %s", err)
	}

	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		FullName: req.GetFullName(),
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "User with username %s already exits", req.GetUsername())
		}
		return nil, status.Errorf(codes.Internal, "An internal error occured: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: bindGRPCUserResponse(user),
	}

	return rsp, nil
}
