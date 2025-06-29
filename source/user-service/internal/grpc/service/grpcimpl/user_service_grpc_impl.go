package grpcimpl

import (
	"context"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/service/userservicepb"
	"thanhldt060802/internal/model"
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
	res.Users = model.FromListUserViewToListUserProto(users)
	return res, nil
}

func (userServiceGRPC *UserServiceGRPCImpl) GetUserById(ctx context.Context, req *userservicepb.GetUserByIdRequest) (*userservicepb.GetUserByIdResponse, error) {
	convertReqDTO := &dto.GetUserByIdRequest{}
	convertReqDTO.Id = req.Id

	user, err := userServiceGRPC.userService.GetUserById(ctx, convertReqDTO)
	if err != nil {
		return nil, err
	}

	res := &userservicepb.GetUserByIdResponse{}
	res.User = model.FromUserViewToUserProto(user)
	return res, nil
}
