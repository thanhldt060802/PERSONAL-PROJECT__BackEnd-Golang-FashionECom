package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type invoiceRepository struct {
}

type InvoiceRepository interface {
	// Integrate with Elasticsearch
	GetAll(ctx context.Context) ([]model.Invoice, error)

	// Main features
	GetByUserId(ctx context.Context, userId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error)
	GetById(ctx context.Context, id int64) (*model.Invoice, error)
	Create(ctx context.Context, newInvoice *model.Invoice) error
	Update(ctx context.Context, updatedInvoice *model.Invoice) error
	DeleteById(ctx context.Context, id int64) error
}

func NewInvoiceRepository() InvoiceRepository {
	return &invoiceRepository{}
}

//
//
// Integrate with Elasticsearch
// ######################################################################################

func (invoiceRepository *invoiceRepository) GetAll(ctx context.Context) ([]model.Invoice, error) {
	var invoices []model.Invoice

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoices).Scan(ctx); err != nil {
		return nil, err
	}

	return invoices, nil
}

//
//
// Main features
// ######################################################################################

func (invoiceRepository *invoiceRepository) GetByUserId(ctx context.Context, userId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error) {
	var invoices []model.Invoice

	query := infrastructure.PostgresDB.NewSelect().Model(&invoices).Where("user_id = ?", userId).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return invoices, nil
}

func (invoiceRepository *invoiceRepository) GetById(ctx context.Context, id int64) (*model.Invoice, error) {
	var invoice model.Invoice

	if err := infrastructure.PostgresDB.NewSelect().Model(&invoice).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (invoiceRepository *invoiceRepository) Create(ctx context.Context, newInvoice *model.Invoice) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newInvoice).Returning("*").Returning("*").Exec(ctx)
	return err
}

func (invoiceRepository *invoiceRepository) Update(ctx context.Context, updatedInvoice *model.Invoice) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedInvoice).Returning("*").Where("id = ?", updatedInvoice.Id).Exec(ctx)
	return err
}

func (invoiceRepository *invoiceRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Invoice{}).Where("id = ?", id).Exec(ctx)
	return err
}
