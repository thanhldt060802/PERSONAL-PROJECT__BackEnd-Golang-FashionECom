package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type InvoiceView struct {
	Id          string              `json:"id"`
	UserId      string              `json:"user_id"`
	TotalAmount int64               `json:"total_amount"`
	Stautus     string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Details     []InvoiceDetailView `json:"details,omitempty"`
}

type InvoiceDetailView struct {
	Id                 string `json:"id"`
	ProductId          string `json:"product_id"`
	Price              int64  `json:"price"`
	DiscountPercentage int32  `json:"discount_percentage"`
	Quantity           int32  `json:"quantity"`
	TotalPrice         int64  `json:"total_price"`
	InvoiceDetailExtraInfo
}

type InvoiceDetailExtraInfo struct {
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	ImageURL string `json:"image_url"`

	CategoryId   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	BrandId      string `json:"brand_id"`
	BrandName    string `json:"brand_name"`
}

func ToInvoiceView(invoice *model.Invoice, details []InvoiceDetailView) *InvoiceView {
	return &InvoiceView{
		Id:          invoice.Id,
		UserId:      invoice.UserId,
		TotalAmount: invoice.TotalAmount,
		Stautus:     invoice.Status,
		CreatedAt:   *invoice.CreatedAt,
		UpdatedAt:   *invoice.UpdatedAt,
		Details:     details,
	}
}

func ToListInvoiceView(invoices []model.Invoice, invoiceIdInvoiceDetailsMap map[string][]model.InvoiceDetail) []InvoiceView {
	invoiceViews := make([]InvoiceView, len(invoices))
	for i, invoice := range invoices {
		invoiceViews[i] = *ToInvoiceView(&invoice, &invoiceIdInvoiceDetailsMap[invoice.Id])
	}
	return invoiceViews
}

func ToInvoiceDetailView(invoiceDetail *model.InvoiceDetail, invoiceDetailExtraInfo InvoiceDetailExtraInfo) *InvoiceDetailView {
	return &InvoiceDetailView{
		Id:                     invoiceDetail.Id,
		ProductId:              invoiceDetail.ProductId,
		Price:                  invoiceDetail.Price,
		DiscountPercentage:     invoiceDetail.DiscountPercentage,
		Quantity:               invoiceDetail.Quantity,
		TotalPrice:             invoiceDetail.TotalPrice,
		InvoiceDetailExtraInfo: invoiceDetailExtraInfo,
	}
}

func ToListInvoiceDetailView(invoiceDetails []model.InvoiceDetail, invoiceDetailExtraInfos []InvoiceDetailExtraInfo) []InvoiceDetailView {
	invoiceDetailViews := make([]InvoiceDetailView, len(invoiceDetails))
	for i := range invoiceDetails {
		invoiceDetailViews[i] = *ToInvoiceDetailView(&invoiceDetails[i], invoiceDetailExtraInfos[i])
	}
	return invoiceDetailViews
}
