package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type productRepository struct {
}

type ProductRepository interface {
	// Main features
	GetById(ctx context.Context, id string) (*model.Product, error)
	Create(ctx context.Context, newProduct *model.Product) error
	Update(ctx context.Context, updatedProduct *model.Product) error
	DeleteById(ctx context.Context, id string) error

	// Elasticsearch integration features
	GetAll(ctx context.Context) ([]model.Product, error)

	// Extra integration with
	GetProductsByListId(ctx context.Context, listId []string) ([]model.Product, error)
	UpdateProductsByListInvoiceDetail(ctx context.Context, listInvoiceDetail []dto.InvoiceDetail) error
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

//
//
// Main features
// ######################################################################################

func (productRepository *productRepository) GetById(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product

	if err := infrastructure.PostgresDB.NewSelect().Model(&product).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &product, nil
}

func (productRepository *productRepository) GetByName(ctx context.Context, name string) (*model.Product, error) {
	var product model.Product

	if err := infrastructure.PostgresDB.NewSelect().Model(&product).Where("name = ?", name).Scan(ctx); err != nil {
		return nil, err
	}

	return &product, nil
}

func (productRepository *productRepository) Create(ctx context.Context, newProduct *model.Product) error {
	newProduct.Id = uuid.New().String()
	_, err := infrastructure.PostgresDB.NewInsert().Model(newProduct).Returning("*").Exec(ctx)
	return err
}

func (productRepository *productRepository) Update(ctx context.Context, updatedProduct *model.Product) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedProduct).Returning("*").Where("id = ?", updatedProduct.Id).Exec(ctx)
	return err
}

func (productRepository *productRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Product{}).Where("id = ?", id).Exec(ctx)
	return err
}

//
//
// Elasticsearch integration features
// ######################################################################################

func (productRepository *productRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product

	if err := infrastructure.PostgresDB.NewSelect().Model(&products).Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

//
//
// Extra GRPC integration features
// ######################################################################################

func (productRepository *productRepository) GetProductsByListId(ctx context.Context, ids []string) ([]model.Product, error) {
	var products []model.Product

	if err := infrastructure.PostgresDB.NewSelect().Model(&products).Where("id IN (?)", bun.In(ids)).Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

func (productRepository *productRepository) UpdateProductsByListInvoiceDetail(ctx context.Context, invoiceDetails []dto.InvoiceDetail) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, invoiceDetail := range invoiceDetails {

	}
}
