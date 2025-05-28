package dto

import (
	"thanhldt060802/internal/grpc/client/userservicepb"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
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

// Receive

func FromUserProtoToUserView(userProto *userservicepb.User) *UserView {
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

// Send

func FromUserViewToUserProto(userView *UserView) *elasticsearchservicepb.User {
	return &elasticsearchservicepb.User{
		Id:        userView.Id,
		FullName:  userView.FullName,
		Email:     userView.Email,
		Username:  userView.Username,
		Address:   userView.Address,
		RoleName:  userView.RoleName,
		CreatedAt: timestamppb.New(userView.CreatedAt),
		UpdatedAt: timestamppb.New(userView.UpdatedAt),
	}
}

func FromListUserViewToListUserProto(userViews []UserView) []*elasticsearchservicepb.User {
	userProtos := make([]*elasticsearchservicepb.User, len(userViews))
	for i := range userProtos {
		userProtos[i] = FromUserViewToUserProto(&userViews[i])
	}

	return userProtos
}

type NumberOfUsersCreatedReport struct {
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	TimeInterval string  `json:"time_interval"`
	Total        float64 `json:"total"`
	Average      float64 `json:"average"`
	Details      []struct {
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Total     float64 `json:"total"`
	} `json:"detail"`
}
