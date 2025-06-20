package repository

import (
	"context"
	"thanhldt060802/infrastructure"
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

	// Elasticsearch integration features (init data for elasticsearch-service)
	GetAll(ctx context.Context) ([]model.Product, error)

	// Order integration features (extra features for order-service)
	GetByListId(ctx context.Context, ids []string) ([]model.Product, error)
	UpdateStocks(ctx context.Context, updatedProducts []model.Product) error
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
// Elasticsearch integration features (init data for elasticsearch-service)
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
// Order integration features (extra features for order-service)
// ######################################################################################

func (productRepository *productRepository) GetByListId(ctx context.Context, ids []string) ([]model.Product, error) {
	var products []model.Product

	if err := infrastructure.PostgresDB.NewSelect().Model(&products).Where("id IN (?)", bun.In(ids)).Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

func (productRepository *productRepository) UpdateStocks(ctx context.Context, updatedProducts []model.Product) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, updatedProduct := range updatedProducts {
		if _, err := tx.NewUpdate().Model(updatedProduct).Returning("*").Where("id = ?", updatedProduct.Id).Exec(ctx); err != nil {
			return err
		}
	}

	return tx.Commit()
}
