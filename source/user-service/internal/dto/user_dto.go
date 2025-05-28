package dto

import (
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/model"
	"time"
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

// Receive

func FromUserProtoToUserView(userProto *elasticsearchservicepb.User) *UserView {
	return &UserView{
		Id:        userProto.Id,
		FullName:  userProto.FullName,
		Email:     userProto.Email,
		Username:  userProto.Username,
		Address:   userProto.Address,
		RoleName:  userProto.RoleName,
		CreatedAt: userProto.CreatedAt.AsTime(),
		UpdatedAt: userProto.UpdatedAt.AsTime(),
	}
}

func FromListUserProtoToListUserView(userProtos []*elasticsearchservicepb.User) []UserView {
	userViews := make([]UserView, len(userProtos))
	for i, userProto := range userProtos {
		userViews[i] = *FromUserProtoToUserView(userProto)
	}

	return userViews
}
