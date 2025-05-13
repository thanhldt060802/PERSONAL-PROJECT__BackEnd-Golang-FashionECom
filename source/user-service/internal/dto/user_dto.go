package dto

import (
	"thanhldt060802/internal/grpc/pb"
	"thanhldt060802/internal/model"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserView struct {
	Id        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Address   string    `json:"address"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserView(user *model.User) *UserView {
	return &UserView{
		Id:        user.Id,
		FullName:  user.FullName,
		Email:     user.Email,
		Username:  user.Username,
		Address:   user.Address,
		RoleName:  user.RoleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToListUserView(users []model.User) []UserView {
	userViews := make([]UserView, len(users))
	for i, user := range users {
		userViews[i] = *ToUserView(&user)
	}

	return userViews
}

func ToUserProtoFromUserView(user *UserView) *pb.User {
	return &pb.User{
		Id:        user.Id,
		FullName:  user.FullName,
		Email:     user.Email,
		Username:  user.Username,
		Address:   user.Address,
		RoleName:  user.RoleName,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToListUserProtoFromListUserView(users []UserView) []*pb.User {
	userProtos := make([]*pb.User, len(users))
	for i, user := range users {
		userProtos[i] = ToUserProtoFromUserView(&user)
	}

	return userProtos
}
