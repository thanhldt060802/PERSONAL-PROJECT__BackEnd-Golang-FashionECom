package grpc_server

import (
	"context"
	"thanhldt060802/internal/grpc-server/pb"
	"thanhldt060802/internal/handler"
)

type GRPCUserService struct {
	pb.UnimplementedUserServiceServer
	userHandler *handler.UserHandler
}

func NewGRPCUserService(userHandler *handler.UserHandler) *GRPCUserService {
	return &GRPCUserService{userHandler: userHandler}
}

func (grpcUserService *GRPCUserService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	return grpcUserService.userHandler.GetAllUsers(ctx, req)
}
