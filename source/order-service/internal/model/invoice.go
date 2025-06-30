package model

import (
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/grpc/service/orderservicepb"
	"time"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Invoice struct {
	bun.BaseModel `bun:"tb_invoice"`

	Id          string     `bun:"id,pk"`
	UserId      string     `bun:"user_id,notnull"`
	TotalAmount int64      `bun:"total_amount,notnull"`
	Status      string     `bun:"status,notnull"`
	CreatedAt   *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type InvoiceDetail struct {
	bun.BaseModel `bun:"tb_invoice_detail"`

	Id                 string `bun:"id,pk"`
	InvoiceId          string `bun:"invoice_id,notnull"`
	ProductId          string `bun:"product_id,notnull"`
	Price              int64  `bun:"price,notnull"`
	DiscountPercentage int32  `bun:"discount_percentage,notnull"`
	Quantity           int32  `bun:"quantity,notnull"`
	TotalPrice         int64  `bun:"total_price,notnull"`
}

type InvoiceView struct {
	bun.BaseModel `bun:"tb_invoice,alias:_invoice"`

	Id          string    `json:"id" bun:"id,pk"`
	UserId      string    `json:"user_id" bun:"user_id"`
	TotalAmount int64     `json:"total_amount" bun:"total_amount"`
	Status      string    `json:"status" bun:"status"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at"`

	InvoiceDetails []*InvoiceDetailView `json:"invoice_details,omitempty" bun:"-"`
}

type InvoiceDetailView struct {
	bun.BaseModel `bun:"tb_invoice_detail,alias:_invoice_detail"`

	Id                 string `json:"id" bun:"id,pk"`
	InvoiceId          string `json:"invoice_id" bun:"invoice_id"`
	ProductId          string `json:"product_id" bun:"product_id"`
	Price              int64  `json:"price" bun:"price"`
	DiscountPercentage int32  `json:"discount_percentage" bun:"discount_percentage"`
	Quantity           int32  `json:"quantity" bun:"quantity"`
	TotalPrice         int64  `json:"total_price" bun:"total_price"`

	ProductName         string `json:"product_name" bun:"product_name"`
	ProductSex          string `json:"product_sex" bun:"product_sex"`
	ProductImageURL     string `json:"product_image_url" bun:"product_image_url"`
	ProductCategoryId   string `json:"product_category_id" bun:"product_category_id"`
	ProductCategoryName string `json:"product_category_name" bun:"product_category_name"`
	ProductBrandId      string `json:"product_brand_id" bun:"product_brand_id"`
	ProductBrandName    string `json:"product_brand_name" bun:"product_brand_name"`
}

// View -> Proto

func FromInvoiceViewToInvoiceProto(invoiceView *InvoiceView) *orderservicepb.Invoice {
	return &orderservicepb.Invoice{
		Id:          invoiceView.Id,
		UserId:      invoiceView.UserId,
		TotalAmount: invoiceView.TotalAmount,
		Status:      invoiceView.Status,
		CreatedAt:   timestamppb.New(invoiceView.CreatedAt),
		UpdatedAt:   timestamppb.New(invoiceView.UpdatedAt),
	}
}

func FromListInvoiceViewToListInvoiceProto(invoiceViews []*InvoiceView) []*orderservicepb.Invoice {
	invoiceProtos := make([]*orderservicepb.Invoice, len(invoiceViews))
	for i, invoiceView := range invoiceViews {
		invoiceProtos[i] = FromInvoiceViewToInvoiceProto(invoiceView)
	}

	return invoiceProtos
}

// Proto -> View

func FromInvoiceProtoToInvoiceView(invoiceProto *elasticsearchservicepb.Invoice) *InvoiceView {
	return &InvoiceView{
		Id:          invoiceProto.Id,
		UserId:      invoiceProto.UserId,
		TotalAmount: invoiceProto.TotalAmount,
		Status:      invoiceProto.Status,
		CreatedAt:   invoiceProto.CreatedAt.AsTime(),
		UpdatedAt:   invoiceProto.UpdatedAt.AsTime(),
	}
}

func FromListInvoiceProtoToListInvoiceView(invoiceProtos []*elasticsearchservicepb.Invoice) []*InvoiceView {
	invoiceViews := make([]*InvoiceView, len(invoiceProtos))
	for i, invoiceProto := range invoiceProtos {
		invoiceViews[i] = FromInvoiceProtoToInvoiceView(invoiceProto)
	}

	return invoiceViews
}
