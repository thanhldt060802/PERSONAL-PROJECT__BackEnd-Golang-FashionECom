package grpcimpl

import (
	"context"
	"thanhldt060802/internal/grpc/service/orderservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/service"
)

type OrderServiceGRPCImpl struct {
	orderservicepb.UnimplementedOrderServiceGRPCServer
	invoiceService service.InvoiceService
}

func NewOrderServiceGRPCImpl(invoiceService service.InvoiceService) *OrderServiceGRPCImpl {
	return &OrderServiceGRPCImpl{invoiceService: invoiceService}
}

func (orderServiceGRPC *OrderServiceGRPCImpl) GetAllInvoices(ctx context.Context, req *orderservicepb.GetAllInvoicesRequest) (*orderservicepb.GetAllInvoicesResponse, error) {
	invoices, err := orderServiceGRPC.invoiceService.GetAllInvoices(ctx)
	if err != nil {
		return nil, err
	}

	res := &orderservicepb.GetAllInvoicesResponse{}
	res.Invoices = model.FromListInvoiceViewToListInvoiceProto(invoices)
	return res, nil
}
