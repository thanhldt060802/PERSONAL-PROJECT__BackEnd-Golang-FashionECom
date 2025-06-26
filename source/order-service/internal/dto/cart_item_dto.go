package dto

import (
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/model"
)

type CartItemView struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`

	ProductName              string `json:"product_name"`
	ProducSex                string `json:"product_sex"`
	ProducPrice              int64  `json:"product_price"`
	ProducDiscountPercentage int32  `json:"product_discount_percentage"`
	ProducImageURL           string `json:"product_image_url"`
	ProducCategoryId         string `json:"product_category_id"`
	ProducCategoryName       string `json:"product_category_name"`
	ProducBrandId            string `json:"product_brand_id"`
	ProducBrandName          string `json:"product_brand_name"`
}

type CartItemExtraInfo struct {
}

func ToCartItemView(cartItem *model.CartItem, productProto *catalogservicepb.Product) *CartItemView {
	return &CartItemView{
		Id:        cartItem.Id,
		UserId:    cartItem.UserId,
		ProductId: cartItem.ProductId,
		Quantity:  cartItem.Quantity,

		ProductName:              productProto.Name,
		ProducSex:                productProto.Sex,
		ProducPrice:              productProto.Price,
		ProducDiscountPercentage: productProto.DiscountPercentage,
		ProducImageURL:           productProto.ImageUrl,
		ProducCategoryId:         productProto.CategoryId,
		ProducCategoryName:       productProto.CategoryName,
		ProducBrandId:            productProto.BrandId,
		ProducBrandName:          productProto.BrandName,
	}
}

func ToListCartItemView(cartItems []*model.CartItem, productProtos []*catalogservicepb.Product) []*CartItemView {
	cartItemViews := make([]*CartItemView, len(cartItems))
	for i := range cartItems {
		cartItemViews[i] = ToCartItemView(cartItems[i], productProtos[i])
	}

	return cartItemViews
}
