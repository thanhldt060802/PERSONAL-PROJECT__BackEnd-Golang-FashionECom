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

	ProductName               string `json:"product_name"`
	ProductSex                string `json:"product_sex"`
	ProductPrice              int64  `json:"product_price"`
	ProductDiscountPercentage int32  `json:"product_discount_percentage"`
	ProductImageURL           string `json:"product_image_url"`
	ProductCategoryId         string `json:"product_category_id"`
	ProductCategoryName       string `json:"product_category_name"`
	ProductBrandId            string `json:"product_brand_id"`
	ProductBrandName          string `json:"product_brand_name"`
}

func ToCartItemView(cartItem *model.CartItem, productProto *catalogservicepb.Product) *CartItemView {
	return &CartItemView{
		Id:        cartItem.Id,
		UserId:    cartItem.UserId,
		ProductId: cartItem.ProductId,
		Quantity:  cartItem.Quantity,

		ProductName:               productProto.Name,
		ProductSex:                productProto.Sex,
		ProductPrice:              productProto.Price,
		ProductDiscountPercentage: productProto.DiscountPercentage,
		ProductImageURL:           productProto.ImageUrl,
		ProductCategoryId:         productProto.CategoryId,
		ProductCategoryName:       productProto.CategoryName,
		ProductBrandId:            productProto.BrandId,
		ProductBrandName:          productProto.BrandName,
	}
}

func ToListCartItemView(cartItems []*model.CartItem, productProtos []*catalogservicepb.Product) []*CartItemView {
	productProtoMap := make(map[string]*catalogservicepb.Product)
	for _, productProto := range productProtos {
		productProtoMap[productProto.Id] = productProto
	}

	cartItemViews := make([]*CartItemView, len(cartItems))
	for i := range cartItems {
		cartItemViews[i] = ToCartItemView(cartItems[i], productProtoMap[cartItems[i].ProductId])
	}

	return cartItemViews
}
