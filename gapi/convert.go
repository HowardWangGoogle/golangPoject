package gapi

import (
	"github.com/techschool/simplebank/dbsq"
	"github.com/techschool/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user dbsq.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
