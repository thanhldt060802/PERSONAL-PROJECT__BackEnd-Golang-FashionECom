package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/service/catalogservicepb"
	"thanhldt060802/internal/service"
)

type CatalogServiceGRPCImpl struct {
	catalogservicepb.UnimplementedCatalogServiceGRPCServer
	productService service.ProductService
}

func NewCatalogServiceGRPCImpl(productService service.ProductService) *CatalogServiceGRPCImpl {
	return &CatalogServiceGRPCImpl{productService: productService}
}

func (grpcCatalogServiceGRPC *CatalogServiceGRPCImpl) GetAllCatalogs(ctx context.Context, req *catalogservicepb.GetAllProductsRequest) (*catalogservicepb.GetAllProductsResponse, error) {
	productProtos, err := grpcCatalogServiceGRPC.productService.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	res := &catalogservicepb.GetAllProductsResponse{}
	res.Products = productProtos
	return res, nil
}
