package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

type invoiceDetailRepository struct {
}

type InvoiceDetailRepository interface {
	GetAllByInvoiceId(ctx context.Context, invoiceId int64) ([]model.InvoiceDetail, error)
	Create(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error
	DeleteByInvoiceId(ctx context.Context, id int64) error

	// Integrate with Elasticsearch

	GetAll(ctx context.Context) ([]model.InvoiceDetail, error)
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

func (invoiceDetailRepository *invoiceDetailRepository) Create(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(&newInvoiceDetails).Returning("*").Exec(ctx)
	return err
}

func (invoiceDetailRepository *invoiceDetailRepository) DeleteByInvoiceId(ctx context.Context, invoiceId int64) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.InvoiceDetail{}).Where("invoice_id = ?", invoiceId).Exec(ctx)
	return err
}

// Integrate with Elasticsearch

func (invoiceDetailRepository *invoiceDetailRepository) GetAll(ctx context.Context) ([]model.InvoiceDetail, error) {
	var invoiceDetails []model.InvoiceDetail

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).Scan(ctx); err != nil {
		return nil, err
	}

	return invoiceDetails, nil
}
