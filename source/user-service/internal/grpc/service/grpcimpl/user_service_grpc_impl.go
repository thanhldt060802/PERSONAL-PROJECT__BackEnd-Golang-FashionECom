package grpcimpl

import (
	"context"
	"net/http"
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

func (grpcUserServiceGRPC *UserServiceGRPCImpl) GetAllUsers(ctx context.Context, req *userservicepb.GetAllUsersRequest) (*userservicepb.GetAllUsersResponse, error) {
	res := &userservicepb.GetAllUsersResponse{}

	userProtos, err := grpcUserServiceGRPC.userService.GetAllUsers(ctx)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all users failed"
	} else {
		res.Status = http.StatusOK
		res.Code = "OK"
		res.Message = "Get all users successful"
		res.Users = userProtos
	}

	return res, nil
}
