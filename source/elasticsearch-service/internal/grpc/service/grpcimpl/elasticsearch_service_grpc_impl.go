package grpcimpl

import (
	"context"
	"net/http"
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
	res := &elasticsearchservicepb.GetUsersResponse{}

	userProtos, err := elasticsearchServiceGRPCImpl.userService.GetUsers(ctx, reqDTO)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get users failed"
		return res, err
	} else {
		res.Status = http.StatusOK
		res.Code = "OK"
		res.Message = "Get users successful"
		res.Users = userProtos
		return res, nil
	}
}
