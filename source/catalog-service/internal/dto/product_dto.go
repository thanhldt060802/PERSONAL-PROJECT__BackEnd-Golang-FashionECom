package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type ProductView struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Sex                string    `json:"sex"`
	Price              int64     `json:"price"`
	DiscountPercentage int32     `json:"discount_percentage"`
	Stock              int32     `json:"stock"`
	ImageURL           string    `json:"image_url"`
	CategoryId         int64     `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	BrandId            int64     `json:"brand_id"`
	BrandName          string    `json:"brand_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func ToProductView(product *model.Product, category *model.Category, brand *model.Brand) *ProductView {
	return &ProductView{
		Id:                 product.Id,
		Name:               product.Name,
		Description:        product.Description,
		Sex:                product.Sex,
		Price:              product.Price,
		DiscountPercentage: product.DiscountPercentage,
		Stock:              product.Stock,
		ImageURL:           product.ImageURL,
		CategoryId:         product.CategoryId,
		CategoryName:       category.Name,
		BrandId:            product.BrandId,
		BrandName:          brand.Name,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
	}
}

func ToListProductView(products []model.Product, categories []model.Category, brands []model.Brand) []ProductView {
	categoryMap := map[int64]*model.Category{}
	for _, category := range categories {
		categoryMap[category.Id] = &category
	}

	brandMap := map[int64]*model.Brand{}
	for _, brand := range brands {
		brandMap[brand.Id] = &brand
	}

	productViews := make([]ProductView, len(products))
	for i, product := range products {
		productViews[i] = *ToProductView(&product, categoryMap[product.CategoryId], brandMap[product.BrandId])
	}

	return productViews
}
