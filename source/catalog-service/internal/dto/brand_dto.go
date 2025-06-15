package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type BrandView struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToBrandView(brand *model.Brand) *BrandView {
	return &BrandView{
		Id:          brand.Id,
		Name:        brand.Name,
		Description: brand.Description,
		CreatedAt:   *brand.CreatedAt,
		UpdatedAt:   *brand.UpdatedAt,
	}
}

func ToListBrandView(brands []model.Brand) []BrandView {
	brandViews := make([]BrandView, len(brands))
	for i, brand := range brands {
		brandViews[i] = *ToBrandView(&brand)
	}

	return brandViews
}
