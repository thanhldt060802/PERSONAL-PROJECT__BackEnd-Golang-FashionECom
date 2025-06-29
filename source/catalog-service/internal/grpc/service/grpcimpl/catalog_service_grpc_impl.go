package grpcimpl

import (
	"context"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/service/catalogservicepb"
	"thanhldt060802/internal/model"
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
	res.Products = model.FromListProductViewToListProductProto(products)
	return res, nil
}

func (catalogServiceGRPC *CatalogServiceGRPCImpl) GetProductById(ctx context.Context, req *catalogservicepb.GetProductByIdRequest) (*catalogservicepb.GetProductByIdResponse, error) {
	convertReqDTO := &dto.GetProductByIdRequest{}
	convertReqDTO.Id = req.Id

	product, err := catalogServiceGRPC.productService.GetProductById(ctx, convertReqDTO)
	if err != nil {
		return nil, err
	}

	res := &catalogservicepb.GetProductByIdResponse{}
	res.Product = model.FromProductViewToProductProto(product)
	return res, nil
}

func (catalogServiceGRPC *CatalogServiceGRPCImpl) UpdateProductStocksByListInvoiceDetail(ctx context.Context, req *catalogservicepb.UpdateProductStocksByListInvoiceDetailRequest) (*catalogservicepb.UpdateProductStocksByListInvoiceDetailResponse, error) {
	convertReqDTO := &dto.UpdateProductStocksByListInvoiceDetailRequest{}
	convertReqDTO.InvoiceDetails = make([]dto.InvoiceDetail, len(req.InvoiceDetails))
	for _, invoiceDetailProto := range req.InvoiceDetails {
		convertReqDTO.InvoiceDetails = append(convertReqDTO.InvoiceDetails, dto.InvoiceDetail{
			ProductId: invoiceDetailProto.ProductId,
			Quantity:  invoiceDetailProto.Quantity,
		})
	}

	if err := catalogServiceGRPC.productService.UpdateProductStocksByListInvoiceDetail(ctx, convertReqDTO); err != nil {
		return nil, err
	}

	res := &catalogservicepb.UpdateProductStocksByListInvoiceDetailResponse{}
	return res, nil
}
