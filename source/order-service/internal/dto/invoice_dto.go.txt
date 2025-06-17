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

func ToListInvoiceView(invoices []model.Invoice) []InvoiceView {
	invoiceViews := make([]InvoiceView, len(invoices))
	for i, invoice := range invoices {
		invoiceViews[i] = *ToInvoiceView(&invoice, nil)
	}
	return invoiceViews
}

func ToInvoiceDetailView(invoiceDetail *model.InvoiceDetail) *InvoiceDetailView {
	return &InvoiceDetailView{
		Id:                 invoiceDetail.Id,
		ProductId:          invoiceDetail.ProductId,
		Price:              invoiceDetail.Price,
		DiscountPercentage: invoiceDetail.DiscountPercentage,
		Quantity:           invoiceDetail.Quantity,
		TotalPrice:         invoiceDetail.TotalPrice,
	}
}

func ToListInvoiceDetailView(invoiceDetails []model.InvoiceDetail) []InvoiceDetailView {
	invoiceDetailViews := make([]InvoiceDetailView, len(invoiceDetails))
	for i, invoiceDetail := range invoiceDetails {
		invoiceDetailViews[i] = *ToInvoiceDetailView(&invoiceDetail)
	}
	return invoiceDetailViews
}
