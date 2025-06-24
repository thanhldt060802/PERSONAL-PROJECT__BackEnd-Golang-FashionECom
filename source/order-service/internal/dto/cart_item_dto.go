package dto

import (
	"thanhldt060802/internal/model"
)

type CartItemView struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	CartItemExtraInfo
}

type CartItemExtraInfo struct {
	Name               string `json:"name"`
	Sex                string `json:"sex"`
	Price              int64  `json:"price"`
	DiscountPercentage int32  `json:"discount_percentage"`
	ImageURL           string `json:"image_url"`

	CategoryId   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	BrandId      string `json:"brand_id"`
	BrandName    string `json:"brand_name"`
}

func ToCartItemView(cartItem *model.CartItem, cartItemExtraInfo *CartItemExtraInfo) *CartItemView {
	return &CartItemView{
		Id:                cartItem.Id,
		UserId:            cartItem.UserId,
		ProductId:         cartItem.ProductId,
		Quantity:          cartItem.Quantity,
		CartItemExtraInfo: *cartItemExtraInfo,
	}
}

func ToListCartItemView(cartItems []model.CartItem, cartItemExtraInfos []CartItemExtraInfo) []CartItemView {
	cartItemViews := make([]CartItemView, len(cartItems))
	for i := range cartItems {
		cartItemViews[i] = *ToCartItemView(&cartItems[i], &cartItemExtraInfos[i])
	}

	return cartItemViews
}
