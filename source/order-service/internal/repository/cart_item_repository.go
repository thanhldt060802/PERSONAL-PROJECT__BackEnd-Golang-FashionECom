package repository

import (
	"context"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/google/uuid"
)

type cartItemRepository struct {
}

type CartItemRepository interface {
	// Main features
	Get(ctx context.Context, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error)
	GetByUserId(ctx context.Context, userId string, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error)
	Create(ctx context.Context, newCartItem *model.CartItem) error
	Update(ctx context.Context, updatedCartItem *model.CartItem) error
	DeleteById(ctx context.Context, id string) error
}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

//
//
// Main features
// ######################################################################################

func (cartItemRepository *cartItemRepository) Get(ctx context.Context, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error) {
	var cartItems []model.CartItem

	query := infrastructure.PostgresDB.NewSelect().Model(&cartItems)

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

func (cartItemRepository *cartItemRepository) GetByUserId(ctx context.Context, userId string, offset *int, limit *int, sortFields []utils.SortField) ([]model.CartItem, error) {
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

func (cartItemRepository *cartItemRepository) Create(ctx context.Context, newCartItem *model.CartItem) error {
	newCartItem.Id = uuid.New().String()
	_, err := infrastructure.PostgresDB.NewInsert().Model(newCartItem).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) Update(ctx context.Context, updatedCartItem *model.CartItem) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedCartItem).Where("id = ?", updatedCartItem.Id).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.CartItem{}).Where("id = ?", id).Exec(ctx)
	return err
}
