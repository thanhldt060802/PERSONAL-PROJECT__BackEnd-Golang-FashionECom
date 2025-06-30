package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
)

type invoiceRepository struct {
}

type InvoiceRepository interface {
	GetViewById(ctx context.Context, id string, dataExpansion bool) (*model.InvoiceView, error)

	GetById(ctx context.Context, id string) (*model.Invoice, error)
	Create(ctx context.Context, newInvoice *model.Invoice, newInvoiceDetails []*model.InvoiceDetail) error
	Update(ctx context.Context, updatedInvoice *model.Invoice) error
	DeleteById(ctx context.Context, id string) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllViews(ctx context.Context, dataExpansion bool) ([]*model.InvoiceView, error)
}

func NewInvoiceRepository() InvoiceRepository {
	return &invoiceRepository{}
}

func (invoiceRepository *invoiceRepository) GetViewById(ctx context.Context, id string, dataExpansion bool) (*model.InvoiceView, error) {
	invoice := new(model.InvoiceView)

	query := infrastructure.PostgresDB.NewSelect().Model(invoice).Where("_invoice.id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if dataExpansion {
		var invoiceDetails []*model.InvoiceDetailView

		query := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).
			TableExpr("tb_invoice_detail AS _invoice_detail").
			Column("_invoice_detail.*").
			ColumnExpr("_product.name AS product_name").
			ColumnExpr("_product.sex AS product_sex").
			ColumnExpr("_product.image_url AS product_image_url").
			ColumnExpr("_product.category_id AS product_category_id").
			ColumnExpr("_product.brand_id AS product_brand_id").
			ColumnExpr("_category.name AS product_category_name").
			ColumnExpr("_brand.name AS product_brand_name").
			Join("JOIN tb_product AS _product ON _product.id = _invoice_detail.product_id").
			Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
			Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
			Where("_invoice_detail.invoice_id = ?", id)

		if err := query.Scan(ctx); err != nil {
			return nil, err
		}

		invoice.InvoiceDetails = invoiceDetails
	}

	return invoice, nil
}

func (invoiceRepository *invoiceRepository) GetById(ctx context.Context, id string) (*model.Invoice, error) {
	invoice := new(model.Invoice)

	query := infrastructure.PostgresDB.NewSelect().Model(&invoice).Where("id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (invoiceRepository *invoiceRepository) Create(ctx context.Context, newInvoice *model.Invoice, newInvoiceDetails []*model.InvoiceDetail) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.NewInsert().Model(newInvoice).Returning("*").Exec(ctx); err != nil {
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

func (invoiceRepository *invoiceRepository) GetAllViews(ctx context.Context, dataExpansion bool) ([]*model.InvoiceView, error) {
	var invoices []*model.InvoiceView

	query := infrastructure.PostgresDB.NewSelect().Model(&invoices)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if dataExpansion {
		for i := range invoices {
			var invoiceDetails []*model.InvoiceDetailView

			query := infrastructure.PostgresDB.NewSelect().Model(&invoiceDetails).
				TableExpr("tb_invoice_detail AS _invoice_detail").
				Column("_invoice_detail.*").
				ColumnExpr("_product.name AS product_name").
				ColumnExpr("_product.sex AS product_sex").
				ColumnExpr("_product.image_url AS product_image_url").
				ColumnExpr("_product.category_id AS product_category_id").
				ColumnExpr("_product.brand_id AS product_brand_id").
				ColumnExpr("_category.name AS product_category_name").
				ColumnExpr("_brand.name AS product_brand_name").
				Join("JOIN tb_product AS _product ON _product.id = _invoice_detail.product_id").
				Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
				Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
				Where("_invoice_detail.invoice_id = ?", invoices[i].Id)

			if err := query.Scan(ctx); err != nil {
				return nil, err
			}

			invoices[i].InvoiceDetails = invoiceDetails
		}
	}

	return invoices, nil
}
