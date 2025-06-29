package model

import (
	"github.com/uptrace/bun"
)

type CartItem struct {
	bun.BaseModel `bun:"tb_cart_item"`

	Id        string `bun:"id,pk"`
	UserId    string `bun:"user_id,notnull"`
	ProductId string `bun:"product_id,notnull"`
	Quantity  int32  `bun:"quantity,notnull"`
}

type CartItemView struct {
	bun.BaseModel `bun:"tb_cart_item,alias:_cart_item"`

	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`

	ProductName               string `json:"product_name" bun:"product_name"`
	ProductSex                string `json:"product_sex" bun:"product_sex"`
	ProductPrice              int64  `json:"product_price" bun:"product_price"`
	ProductDiscountPercentage int32  `json:"product_discount_percentage" bun:"product_discount_percentage"`
	ProductImageURL           string `json:"product_image_url" bun:"product_image_url"`
	ProductCategoryId         string `json:"product_category_id" bun:"product_category_id"`
	ProductCategoryName       string `json:"product_category_name" bun:"product_category_name"`
	ProductBrandId            string `json:"product_brand_id" bun:"product_brand_id"`
	ProductBrandName          string `json:"product_brand_name" bun:"product_brand_name"`
}
