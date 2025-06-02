package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
	"thanhldt060802/internal/service"
)

type ElasticsearchServiceGRPCImpl struct {
	elasticsearchservicepb.UnimplementedElasticsearchServiceGRPCServer
	userService service.UserService
}

func NewElasticsearchServiceGRPCImpl(userService service.UserService) *ElasticsearchServiceGRPCImpl {
	return &ElasticsearchServiceGRPCImpl{userService: userService}
}

func (elasticsearchServiceGRPCImpl *ElasticsearchServiceGRPCImpl) GetUsers(ctx context.Context, reqDTO *elasticsearchservicepb.GetUsersRequest) (*elasticsearchservicepb.GetUsersResponse, error) {
	userProtos, err := elasticsearchServiceGRPCImpl.userService.GetUsers(ctx, reqDTO)
	if err != nil {
		return nil, err
	}

	res := &elasticsearchservicepb.GetUsersResponse{}
	res.Users = userProtos
	return res, nil
}
