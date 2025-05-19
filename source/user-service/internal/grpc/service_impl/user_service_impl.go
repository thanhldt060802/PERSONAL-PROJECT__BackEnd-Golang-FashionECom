package grpcimpl

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/pb"
	"thanhldt060802/internal/service"
)

type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserServiceImpl(userService service.UserService) *UserServiceImpl {
	return &UserServiceImpl{userService: userService}
}

func (grpcUserService *UserServiceImpl) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := grpcUserService.userService.GetAllUsers(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all users failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &pb.GetAllUsersResponse{}
	res.Users = dto.ToListUserProtoFromListUserView(users)
	return res, nil
}
