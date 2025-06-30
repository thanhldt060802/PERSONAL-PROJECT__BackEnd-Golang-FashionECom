package dto

import (
	"thanhldt060802/internal/grpc/client/orderservicepb"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvoiceView struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	TotalAmount int64     `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Receive

func FromInvoiceProtoToInvoiceView(invoiceProto *orderservicepb.Invoice) *InvoiceView {
	return &InvoiceView{
		Id:          invoiceProto.Id,
		UserId:      invoiceProto.UserId,
		TotalAmount: invoiceProto.TotalAmount,
		Status:      invoiceProto.Status,
		CreatedAt:   invoiceProto.CreatedAt.AsTime(),
		UpdatedAt:   invoiceProto.UpdatedAt.AsTime(),
	}
}

// Send

func FromInvoiceViewToInvoiceProto(invoiceView *InvoiceView) *elasticsearchservicepb.Invoice {
	return &elasticsearchservicepb.Invoice{
		Id:          invoiceView.Id,
		UserId:      invoiceView.UserId,
		TotalAmount: invoiceView.TotalAmount,
		Status:      invoiceView.Status,
		CreatedAt:   timestamppb.New(invoiceView.CreatedAt),
		UpdatedAt:   timestamppb.New(invoiceView.UpdatedAt),
	}
}

func FromListInvoiceViewToListInvoiceProto(invoiceViews []InvoiceView) []*elasticsearchservicepb.Invoice {
	invoiceProtos := make([]*elasticsearchservicepb.Invoice, len(invoiceViews))
	for i := range invoiceProtos {
		invoiceProtos[i] = FromInvoiceViewToInvoiceProto(&invoiceViews[i])
	}

	return invoiceProtos
}
