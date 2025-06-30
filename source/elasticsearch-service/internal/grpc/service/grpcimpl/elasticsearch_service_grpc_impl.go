package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
	"thanhldt060802/internal/service"
)

type ElasticsearchServiceGRPCImpl struct {
	elasticsearchservicepb.UnimplementedElasticsearchServiceGRPCServer
	userService    service.UserService
	catalogService service.CatalogService
	orderService   service.OrderService
}

func NewElasticsearchServiceGRPCImpl(userService service.UserService, catalogService service.CatalogService, orderService service.OrderService) *ElasticsearchServiceGRPCImpl {
	return &ElasticsearchServiceGRPCImpl{
		userService:    userService,
		catalogService: catalogService,
		orderService:   orderService,
	}
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

func (elasticsearchServiceGRPCImpl *ElasticsearchServiceGRPCImpl) GetProducts(ctx context.Context, reqDTO *elasticsearchservicepb.GetProductsRequest) (*elasticsearchservicepb.GetProductsResponse, error) {
	productProtos, err := elasticsearchServiceGRPCImpl.catalogService.GetProducts(ctx, reqDTO)
	if err != nil {
		return nil, err
	}

	res := &elasticsearchservicepb.GetProductsResponse{}
	res.Products = productProtos
	return res, nil
}

func (elasticsearchServiceGRPCImpl *ElasticsearchServiceGRPCImpl) GetInvoices(ctx context.Context, reqDTO *elasticsearchservicepb.GetInvoicesRequest) (*elasticsearchservicepb.GetInvoicesResponse, error) {
	invoiceProtos, err := elasticsearchServiceGRPCImpl.orderService.GetInvoices(ctx, reqDTO)
	if err != nil {
		return nil, err
	}

	res := &elasticsearchservicepb.GetInvoicesResponse{}
	res.Invoices = invoiceProtos
	return res, nil
}
