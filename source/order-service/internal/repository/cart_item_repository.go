package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/google/uuid"
)

type cartItemRepository struct {
}

type CartItemRepository interface {
	GetViews(ctx context.Context, offset int, limit int, sortFields []*utils.SortField) ([]*model.CartItemView, error)
	GetAllViewsByUserId(ctx context.Context, userId string) ([]*model.CartItemView, error)
	GetViewsByUserId(ctx context.Context, userId string, offset int, limit int, sortFields []*utils.SortField) ([]*model.CartItemView, error)

	GetById(ctx context.Context, id string) (*model.CartItem, error)
	Create(ctx context.Context, newCartItem *model.CartItem) error
	Update(ctx context.Context, updatedCartItem *model.CartItem) error
	DeleteByUserId(ctx context.Context, userId string) error
	DeleteById(ctx context.Context, id string) error
}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

func (cartItemRepository *cartItemRepository) GetViews(ctx context.Context, offset int, limit int, sortFields []*utils.SortField) ([]*model.CartItemView, error) {
	var cartItems []*model.CartItemView

	query := infrastructure.PostgresDB.NewSelect().Model(&cartItems).
		TableExpr("tb_cart_item AS _cart_item").
		Column("_cart_item.*").
		ColumnExpr("_product.name AS product_name").
		ColumnExpr("_product.sex AS product_sex").
		ColumnExpr("_product.price AS product_price").
		ColumnExpr("_product.discount_percentage AS product_discount_percentage").
		ColumnExpr("_product.image_url AS product_image_url").
		ColumnExpr("_product.category_id AS product_category_id").
		ColumnExpr("_product.brand_id AS product_brand_id").
		ColumnExpr("_category.name AS product_category_name").
		ColumnExpr("_brand.name AS product_brand_name").
		Join("JOIN tb_product AS _product ON _product.id = _cart_item.product_id").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
		Offset(offset).
		Limit(limit)

	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("_cart_item.%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) GetAllViewsByUserId(ctx context.Context, userId string) ([]*model.CartItemView, error) {
	var cartItems []*model.CartItemView

	err := infrastructure.PostgresDB.NewSelect().Model(&cartItems).
		TableExpr("tb_cart_item AS _cart_item").
		Column("_cart_item.*").
		ColumnExpr("_product.name AS product_name").
		ColumnExpr("_product.sex AS product_sex").
		ColumnExpr("_product.price AS product_price").
		ColumnExpr("_product.discount_percentage AS product_discount_percentage").
		ColumnExpr("_product.image_url AS product_image_url").
		ColumnExpr("_product.category_id AS product_category_id").
		ColumnExpr("_product.brand_id AS product_brand_id").
		ColumnExpr("_category.name AS product_category_name").
		ColumnExpr("_brand.name AS product_brand_name").
		Join("JOIN tb_product AS _product ON _product.id = _cart_item.product_id").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
		Where("_cart_item.user_id = ?", userId).Scan(ctx)

	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) GetViewsByUserId(ctx context.Context, userId string, offset int, limit int, sortFields []*utils.SortField) ([]*model.CartItemView, error) {
	var cartItems []*model.CartItemView

	query := infrastructure.PostgresDB.NewSelect().Model(&cartItems).
		TableExpr("tb_cart_item AS _cart_item").
		Column("_cart_item.*").
		ColumnExpr("_product.name AS product_name").
		ColumnExpr("_product.sex AS product_sex").
		ColumnExpr("_product.price AS product_price").
		ColumnExpr("_product.discount_percentage AS product_discount_percentage").
		ColumnExpr("_product.image_url AS product_image_url").
		ColumnExpr("_product.category_id AS product_category_id").
		ColumnExpr("_product.brand_id AS product_brand_id").
		ColumnExpr("_category.name AS product_category_name").
		ColumnExpr("_brand.name AS product_brand_name").
		Join("JOIN tb_product AS _product ON _product.id = _cart_item.product_id").
		Join("JOIN tb_category AS _category ON _category.id = _product.category_id").
		Join("JOIN tb_brand AS _brand ON _brand.id = _product.brand_id").
		Where("_cart_item.user_id = ?", userId).
		Offset(offset).
		Limit(limit)

	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("_cart_item.%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) GetById(ctx context.Context, id string) (*model.CartItem, error) {
	cartItem := new(model.CartItem)

	if err := infrastructure.PostgresDB.NewSelect().Model(cartItem).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (cartItemRepository *cartItemRepository) Create(ctx context.Context, newCartItem *model.CartItem) error {
	newCartItem.Id = uuid.New().String()
	_, err := infrastructure.PostgresDB.NewInsert().Model(newCartItem).Returning("*").Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) Update(ctx context.Context, updatedCartItem *model.CartItem) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedCartItem).Where("id = ?", updatedCartItem.Id).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) DeleteByUserId(ctx context.Context, userId string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.CartItem{}).Where("user_id = ?", userId).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.CartItem{}).Where("id = ?", id).Exec(ctx)
	return err
}
