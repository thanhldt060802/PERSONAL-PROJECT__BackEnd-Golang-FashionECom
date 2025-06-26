package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

type invoiceRepository struct {
}

type InvoiceRepository interface {
	GetById(ctx context.Context, id string) (*model.Invoice, []*model.InvoiceDetail, error)
	Create(ctx context.Context, newInvoice *model.Invoice, newInvoiceDetails []*model.InvoiceDetail) error
	Update(ctx context.Context, updatedInvoice *model.Invoice) error
	DeleteById(ctx context.Context, id string) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAll(ctx context.Context) ([]*model.Invoice, map[string][]*model.InvoiceDetail, error)
}

func NewInvoiceRepository() InvoiceRepository {
	return &invoiceRepository{}
}

func (invoiceRepository *invoiceRepository) GetById(ctx context.Context, id string) (*model.Invoice, []*model.InvoiceDetail, error) {
	var invoice model.Invoice
	var invoiceDetails []*model.InvoiceDetail

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoice).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, nil, err
	}

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).Where("invoice_id = ?", id).Scan(ctx); err != nil {
		return nil, nil, err
	}

	return &invoice, invoiceDetails, nil
}

func (invoiceRepository *invoiceRepository) Create(ctx context.Context, newInvoice *model.Invoice, newInvoiceDetails []*model.InvoiceDetail) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.NewInsert().Model(newInvoice).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(&newInvoiceDetails).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (invoiceRepository *invoiceRepository) Update(ctx context.Context, updatedInvoice *model.Invoice) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedInvoice).Where("id = ?", updatedInvoice.Id).Exec(ctx)
	return err
}

func (invoiceRepository *invoiceRepository) DeleteById(ctx context.Context, id string) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.NewDelete().Model(&model.Invoice{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewDelete().Model(&model.InvoiceDetail{}).Where("invoice_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (invoiceRepository *invoiceRepository) GetAll(ctx context.Context) ([]*model.Invoice, map[string][]*model.InvoiceDetail, error) {
	var invoices []*model.Invoice
	var invoiceDetails []*model.InvoiceDetail

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoices).Scan(ctx); err != nil {
		return nil, nil, err
	}

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).Scan(ctx); err != nil {
		return nil, nil, err
	}

	invoiceDetailsMap := make(map[string][]*model.InvoiceDetail)
	for _, invoiceDetail := range invoiceDetails {
		invoiceDetailsMap[invoiceDetail.InvoiceId] = append(invoiceDetailsMap[invoiceDetail.InvoiceId], invoiceDetail)
	}

	return invoices, invoiceDetailsMap, nil
}
