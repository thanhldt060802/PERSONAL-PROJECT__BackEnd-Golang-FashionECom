package grpcimpl

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/pb"
	"thanhldt060802/internal/service"
)

type ProductServiceImpl struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewProductServiceImpl(productService service.ProductService) *ProductServiceImpl {
	return &ProductServiceImpl{productService: productService}
}

func (grpcProductService *ProductServiceImpl) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	products, err := grpcProductService.productService.GetAllProducts(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all products failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &pb.GetAllProductsResponse{}
	res.Products = dto.ToListProductProtoFromListProductView(products)
	return res, nil
}
