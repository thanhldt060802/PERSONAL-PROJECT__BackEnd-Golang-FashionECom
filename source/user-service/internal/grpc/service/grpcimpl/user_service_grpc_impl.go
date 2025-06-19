package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/service/userservicepb"
	"thanhldt060802/internal/service"
)

type UserServiceGRPCImpl struct {
	userservicepb.UnimplementedUserServiceGRPCServer
	userService service.UserService
}

func NewUserServiceGRPCImpl(userService service.UserService) *UserServiceGRPCImpl {
	return &UserServiceGRPCImpl{userService: userService}
}

func (userServiceGRPC *UserServiceGRPCImpl) GetAllUsers(ctx context.Context, req *userservicepb.GetAllUsersRequest) (*userservicepb.GetAllUsersResponse, error) {
	userProtos, err := userServiceGRPC.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &userservicepb.GetAllUsersResponse{}
	res.Users = userProtos
	return res, nil
}
