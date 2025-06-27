package dto

import (
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/model"
	"time"
)

type InvoiceView struct {
	Id             string              `json:"id"`
	UserId         string              `json:"user_id"`
	TotalAmount    int64               `json:"total_amount"`
	Status         string              `json:"status"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	InvoiceDetails []InvoiceDetailView `json:"invoice_details"`
}

type InvoiceDetailView struct {
	Id                 string `json:"id"`
	ProductId          string `json:"product_id"`
	Price              int64  `json:"price"`
	DiscountPercentage int32  `json:"discount_percentage"`
	Quantity           int32  `json:"quantity"`
	TotalPrice         int64  `json:"total_price"`

	ProductName         string `json:"product_name"`
	ProductSex          string `json:"product_sex"`
	ProductImageURL     string `json:"product_image_url"`
	ProductCategoryId   string `json:"product_category_id"`
	ProductCategoryName string `json:"product_category_name"`
	ProductBrandId      string `json:"product_brand_id"`
	ProductBrandName    string `json:"product_brand_name"`
}

func ToInvoiceView(invoice *model.Invoice, invoiceDetails []*model.InvoiceDetail, productProtos []*catalogservicepb.Product) *InvoiceView {
	return &InvoiceView{
		Id:             invoice.Id,
		UserId:         invoice.UserId,
		TotalAmount:    invoice.TotalAmount,
		Status:         invoice.Status,
		CreatedAt:      *invoice.CreatedAt,
		UpdatedAt:      *invoice.UpdatedAt,
		InvoiceDetails: ToListInvoiceDetailView(invoiceDetails, productProtos),
	}
}

func ToListInvoiceView(invoices []*model.Invoice, invoiceDetailsMap map[string][]*model.InvoiceDetail, productProtosMap map[string][]*catalogservicepb.Product) []*InvoiceView {
	invoiceViews := make([]*InvoiceView, len(invoices))
	for i, invoice := range invoices {
		invoiceViews[i] = ToInvoiceView(invoice, invoiceDetailsMap[invoice.Id], productProtosMap[invoice.Id])
	}
	return invoiceViews
}

func ToInvoiceDetailView(invoiceDetail *model.InvoiceDetail, productProto *catalogservicepb.Product) *InvoiceDetailView {
	return &InvoiceDetailView{
		Id:                 invoiceDetail.Id,
		ProductId:          invoiceDetail.ProductId,
		Price:              invoiceDetail.Price,
		DiscountPercentage: invoiceDetail.DiscountPercentage,
		Quantity:           invoiceDetail.Quantity,
		TotalPrice:         invoiceDetail.TotalPrice,

		ProductName:         productProto.Name,
		ProductSex:          productProto.Sex,
		ProductImageURL:     productProto.ImageUrl,
		ProductCategoryId:   productProto.CategoryId,
		ProductCategoryName: productProto.CategoryName,
		ProductBrandId:      productProto.BrandId,
		ProductBrandName:    productProto.BrandName,
	}
}

func ToListInvoiceDetailView(invoiceDetails []*model.InvoiceDetail, productProtos []*catalogservicepb.Product) []InvoiceDetailView {
	productProtoMap := make(map[string]*catalogservicepb.Product)
	for _, productProto := range productProtos {
		productProtoMap[productProto.Id] = productProto
	}

	invoiceDetailViews := make([]InvoiceDetailView, len(invoiceDetails))
	for i := range invoiceDetails {
		invoiceDetailViews[i] = *ToInvoiceDetailView(invoiceDetails[i], productProtoMap[invoiceDetails[i].ProductId])
	}

	return invoiceDetailViews
}
