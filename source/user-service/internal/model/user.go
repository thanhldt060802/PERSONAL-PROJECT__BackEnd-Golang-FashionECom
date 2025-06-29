package model

import (
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/grpc/service/userservicepb"
	"time"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	bun.BaseModel `bun:"tb_user"`

	Id             string     `bun:"id,pk"`
	FullName       string     `bun:"full_name,notnull"`
	Email          string     `bun:"email,notnull"`
	Username       string     `bun:"username,notnull"`
	HashedPassword string     `bun:"hashed_password,notnull"`
	Address        string     `bun:"address,notnull"`
	RoleName       string     `bun:"role_name,notnull"`
	CreatedAt      *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type UserView struct {
	bun.BaseModel `bun:"tb_user,alias:_user"`

	Id        string    `json:"id" bun:"id,pk"`
	FullName  string    `json:"full_name" bun:"full_name"`
	Email     string    `json:"email" bun:"email"`
	Username  string    `json:"username" bun:"username"`
	Address   string    `json:"address" bun:"address"`
	RoleName  string    `json:"role_name" bun:"role_name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`
}

// View -> Proto

func FromUserViewToUserProto(userView *UserView) *userservicepb.User {
	return &userservicepb.User{
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

func FromListUserViewToListUserProto(userViews []*UserView) []*userservicepb.User {
	userProtos := make([]*userservicepb.User, len(userViews))
	for i, userView := range userViews {
		userProtos[i] = FromUserViewToUserProto(userView)
	}

	return userProtos
}

// Proto -> View

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

func FromListUserProtoToListUserView(userProtos []*elasticsearchservicepb.User) []*UserView {
	userViews := make([]*UserView, len(userProtos))
	for i, userProto := range userProtos {
		userViews[i] = FromUserProtoToUserView(userProto)
	}

	return userViews
}
