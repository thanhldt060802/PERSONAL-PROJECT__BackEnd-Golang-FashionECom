package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type cartItemRepository struct {
}

type CartItemRepository interface {
	// Main features
	GetByUserId(ctx context.Context, userId int64, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error)
	GetById(ctx context.Context, id int64) (*model.CartItem, error)
	Create(ctx context.Context, newCartItem *model.CartItem) error
	Update(ctx context.Context, updatedCartItem *model.CartItem) error
	DeleteById(ctx context.Context, id int64) error
}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

//
//
// Main features
// ######################################################################################

func (cartItemRepository *cartItemRepository) GetByUserId(ctx context.Context, userId int64, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error) {
	var cartItems []model.CartItem

	query := infrastructure.PostgresDB.NewSelect().Model(&cartItems).Where("user_id = ?", userId)

	if offset != nil {
		query = query.Offset(*offset)
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) GetById(ctx context.Context, id int64) (*model.CartItem, error) {
	var cartItem model.CartItem

	if err := infrastructure.PostgresDB.NewSelect().Model(&cartItem).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func (cartItemRepository *cartItemRepository) Create(ctx context.Context, newCartItem *model.CartItem) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newCartItem).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) Update(ctx context.Context, updatedCartItem *model.CartItem) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedCartItem).Where("id = ?", updatedCartItem.Id).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.CartItem{}).Where("id = ?", id).Exec(ctx)
	return err
}
