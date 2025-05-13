package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/pb"
	"thanhldt060802/internal/handler"
)

type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
	userHandler *handler.UserHandler
}

func NewUserServiceImpl(userHandler *handler.UserHandler) *UserServiceImpl {
	return &UserServiceImpl{userHandler: userHandler}
}

func (grpcUserService *UserServiceImpl) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	return grpcUserService.userHandler.GetAllUsers(ctx, req)
}
