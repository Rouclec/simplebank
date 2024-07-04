package gapi

import (
	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func bindGRPCUserResponse(user db.Users) *pb.User {
	return &pb.User{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}