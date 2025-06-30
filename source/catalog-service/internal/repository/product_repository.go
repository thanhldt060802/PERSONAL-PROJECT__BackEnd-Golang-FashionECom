package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"

	"github.com/uptrace/bun"
)

type productRepository struct {
}

type ProductRepository interface {
	GetViewById(ctx context.Context, id string) (*model.ProductView, error)

	GetByListId(ctx context.Context, ids []string) ([]*model.Product, error)
	GetById(ctx context.Context, id string) (*model.Product, error)
	Create(ctx context.Context, newProduct *model.Product) error
	Update(ctx context.Context, updatedProduct *model.Product) error
	DeleteById(ctx context.Context, id string) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllViews(ctx context.Context) ([]*model.ProductView, error)

	// Order integration (extra features for order-service)
	UpdateStocks(ctx context.Context, updatedProducts []*model.Product) error
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (productRepository *productRepository) GetViewById(ctx context.Context, id string) (*model.ProductView, error) {
	product := new(model.ProductView)

	query := infrastructure.PostgresDB.NewSelect().Model(product).
		TableExpr("tb_product AS _product").
		Column("_product.*").
		ColumnExpr("_category.name AS category_name").
		ColumnExpr("_brand.name AS brand_name").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
		Where("_product.id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return product, nil
}

func (productRepository *productRepository) GetByListId(ctx context.Context, ids []string) ([]*model.Product, error) {
	var products []*model.Product

	query := infrastructure.PostgresDB.NewSelect().Model(products).Where("id IN (?)", bun.In(ids))

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

func (productRepository *productRepository) GetById(ctx context.Context, id string) (*model.Product, error) {
	product := new(model.Product)

	query := infrastructure.PostgresDB.NewSelect().Model(product).Where("id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return product, nil
}

func (productRepository *productRepository) Create(ctx context.Context, newProduct *model.Product) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newProduct).Returning("*").Exec(ctx)
	return err
}

func (productRepository *productRepository) Update(ctx context.Context, updatedProduct *model.Product) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedProduct).Where("id = ?", updatedProduct.Id).Exec(ctx)
	return err
}

func (productRepository *productRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Product{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (productRepository *productRepository) GetAllViews(ctx context.Context) ([]*model.ProductView, error) {
	var products []*model.ProductView

	query := infrastructure.PostgresDB.NewSelect().Model(products).
		TableExpr("tb_product AS _product").
		Column("_product.*").
		ColumnExpr("_category.name AS category_name").
		ColumnExpr("_brand.name AS brand_name").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id")

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

func (productRepository *productRepository) GetViewsByListId(ctx context.Context, ids []string) ([]*model.ProductView, error) {
	var products []*model.ProductView

	query := infrastructure.PostgresDB.NewSelect().Model(products).
		TableExpr("tb_product AS _product").
		Column("_product.*").
		ColumnExpr("_category.name AS category_name").
		ColumnExpr("_brand.name AS brand_name").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
		Where("_product.id IN (?)", bun.In(ids))

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return products, nil
}

func (productRepository *productRepository) UpdateStocks(ctx context.Context, updatedProducts []*model.Product) error {
	tx, err := infrastructure.PostgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, updatedProduct := range updatedProducts {
		if _, err := tx.NewUpdate().Model(updatedProduct).Where("id = ?", updatedProduct.Id).Exec(ctx); err != nil {
			return err
		}
	}

	return tx.Commit()
}
