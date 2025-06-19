package grpcimpl

import (
	"context"
	"thanhldt060802/internal/dto"
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

func (catalogServiceGRPC *CatalogServiceGRPCImpl) GetAllProducts(ctx context.Context, req *catalogservicepb.GetAllProductsRequest) (*catalogservicepb.GetAllProductsResponse, error) {
	products, err := catalogServiceGRPC.productService.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	res := &catalogservicepb.GetAllProductsResponse{}
	res.Products = dto.FromListProductViewToListProductProto(products)
	return res, nil
}

func (catalogServiceGRPC *CatalogServiceGRPCImpl) GetProductsByListId(ctx context.Context, req *catalogservicepb.GetProductsByListIdRequest) (*catalogservicepb.GetProductsByListIdResponse, error) {
	productProtos := []*catalogservicepb.Product{}
	for _, id := range req.ListId {
		reqDTO := &dto.GetProductByIdRequest{}
		reqDTO.Id = id

		product, err := catalogServiceGRPC.productService.GetProductById(ctx, reqDTO)
		if err != nil {
			return nil, err
		}
		productProtos = append(productProtos, dto.FromProductViewToProductProto(product))
	}

	res := &catalogservicepb.GetProductsByListIdResponse{}
	res.Products = productProtos
	return res, nil
}

func (catalogServiceGRPC *CatalogServiceGRPCImpl) UpdateProductsByListInvoiceDetail(ctx context.Context, req *catalogservicepb.UpdateProductsByListInvoiceDetailRequest) (*catalogservicepb.UpdateProductsByListInvoiceDetailResponse, error) {
	for _, invoiceDetail := range req.InvoiceDetails {
		reqDTO := &dto.UpdateProductByIdRequest{}
		reqDTO.Id = invoiceDetail.ProductId
		stock := *reqDTO.Body.Stock - invoiceDetail.Quantity
		reqDTO.Body.Stock = &stock

		if err := catalogServiceGRPC.productService.UpdateProductById(ctx, reqDTO); err != nil {
			return nil, err
		}
	}

	res := &catalogservicepb.UpdateProductsByListInvoiceDetailResponse{}
	return res, nil
}
