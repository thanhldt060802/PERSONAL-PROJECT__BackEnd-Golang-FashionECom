package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"

	"github.com/uptrace/bun"
)

type invoiceDetailRepository struct {
}

type InvoiceDetailRepository interface {
	// Main features
	GetAllByInvoiceId(ctx context.Context, invoiceId int64) ([]model.InvoiceDetail, error)
	CreateMany(ctx context.Context, newInvoiceDetails []model.InvoiceDetail, tx bun.Tx) error
	DeleteByInvoiceId(ctx context.Context, invoiceId int64, tx bun.Tx) error
}

func NewInvoiceDetailRepository() InvoiceDetailRepository {
	return &invoiceDetailRepository{}
}

//
//
// Main features
// ######################################################################################

func (invoiceDetailRepository *invoiceDetailRepository) GetAllByInvoiceId(ctx context.Context, invoiceId int64) ([]model.InvoiceDetail, error) {
	var invoiceDetails []model.InvoiceDetail

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).Where("invoice_id = ?", invoiceId).Scan(ctx); err != nil {
		return nil, err
	}

	return invoiceDetails, nil
}

func (invoiceDetailRepository *invoiceDetailRepository) CreateMany(ctx context.Context, newInvoiceDetails []model.InvoiceDetail, tx bun.Tx) error {
	_, err := tx.NewInsert().Model(&newInvoiceDetails).Exec(ctx)
	return err
}

func (invoiceDetailRepository *invoiceDetailRepository) DeleteByInvoiceId(ctx context.Context, invoiceId int64, tx bun.Tx) error {
	_, err := tx.NewDelete().Model(&model.InvoiceDetail{}).Where("invoice_id = ?", invoiceId).Exec(ctx)
	return err
}
