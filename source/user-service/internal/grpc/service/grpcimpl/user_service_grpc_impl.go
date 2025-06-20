package grpcimpl

import (
	"context"
	"thanhldt060802/internal/dto"
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
	users, err := userServiceGRPC.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &userservicepb.GetAllUsersResponse{}
	res.Users = dto.FromListUserViewToListUserProto(users)
	return res, nil
}
