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
	convertReqDTO := &dto.GetProductsByListIdRequest{}
	convertReqDTO.Ids = req.Ids

	products, err := catalogServiceGRPC.productService.GetProductsByListId(ctx, convertReqDTO)
	if err != nil {
		return nil, err
	}

	res := &catalogservicepb.GetProductsByListIdResponse{}
	res.Products = dto.FromListProductViewToListProductProto(products)
	return res, nil
}

func (catalogServiceGRPC *CatalogServiceGRPCImpl) UpdateProductsByListInvoiceDetail(ctx context.Context, req *catalogservicepb.UpdateProductsByListInvoiceDetailRequest) (*catalogservicepb.UpdateProductsByListInvoiceDetailResponse, error) {
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

	res := &catalogservicepb.UpdateProductsByListInvoiceDetailResponse{}
	return res, nil
}
