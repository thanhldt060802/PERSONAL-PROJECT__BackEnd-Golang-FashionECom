package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type CartItemView struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	ProductId int64     `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCartItemView(cartItem *model.CartItem) *CartItemView {
	return &CartItemView{
		Id:        cartItem.Id,
		UserId:    cartItem.UserId,
		ProductId: cartItem.ProductId,
		Quantity:  cartItem.Quantity,
		CreatedAt: cartItem.CreatedAt,
		UpdatedAt: cartItem.UpdatedAt,
	}
}

func ToListCartItemView(cartItems []model.CartItem) []CartItemView {
	cartItemViews := make([]CartItemView, len(cartItems))
	for i, cartItem := range cartItems {
		cartItemViews[i] = *ToCartItemView(&cartItem)
	}

	return cartItemViews
}
