package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type CategoryView struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCategoryView(category *model.Category) *CategoryView {
	return &CategoryView{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: *category.CreatedAt,
		UpdatedAt: *category.UpdatedAt,
	}
}

func ToListCategoryView(categorys []*model.Category) []*CategoryView {
	categoryViews := make([]*CategoryView, len(categorys))
	for i, category := range categorys {
		categoryViews[i] = ToCategoryView(category)
	}

	return categoryViews
}
