package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

type invoiceDetailRepository struct {
}

type InvoiceDetailRepository interface {
	// Main features
	GetAllByInvoiceId(ctx context.Context, invoiceId int64) ([]model.InvoiceDetail, error)
	CreateMany(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error
	DeleteByInvoiceId(ctx context.Context, invoiceId int64) error
}

func NewInvoiceDetailRepository() InvoiceDetailRepository {
	return &invoiceDetailRepository{}
}

func (invoiceDetailRepository *invoiceDetailRepository) GetAllByInvoiceId(ctx context.Context, invoiceId int64) ([]model.InvoiceDetail, error) {
	var invoiceDetails []model.InvoiceDetail

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).Where("invoice_id = ?", invoiceId).Scan(ctx); err != nil {
		return nil, err
	}

	return invoiceDetails, nil
}

func (invoiceDetailRepository *invoiceDetailRepository) CreateMany(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(&newInvoiceDetails).Exec(ctx)
	return err
}

func (invoiceDetailRepository *invoiceDetailRepository) DeleteByInvoiceId(ctx context.Context, invoiceId int64) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.InvoiceDetail{}).Where("invoice_id = ?", invoiceId).Exec(ctx)
	return err
}
