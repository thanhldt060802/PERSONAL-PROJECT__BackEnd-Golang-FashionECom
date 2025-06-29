package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type categoryRepository struct {
}

type CategoryRepository interface {
	GetAllViews(ctx context.Context, sortFields []*utils.SortField) ([]*model.CategoryView, error)
	GetViewById(ctx context.Context, id string) (*model.CategoryView, error)

	GetById(ctx context.Context, id string) (*model.Category, error)
	GetByName(ctx context.Context, name string) (*model.Category, error)
	Create(ctx context.Context, newCategory *model.Category) error
	Update(ctx context.Context, updatedCategory *model.Category) error
	DeleteById(ctx context.Context, id string) error
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (categoryRepository *categoryRepository) GetAllViews(ctx context.Context, sortFields []*utils.SortField) ([]*model.CategoryView, error) {
	var categories []*model.CategoryView

	query := infrastructure.PostgresDB.NewSelect().Model(&categories)

	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("_category.%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return categories, nil
}

func (categoryRepository *categoryRepository) GetViewById(ctx context.Context, id string) (*model.CategoryView, error) {
	category := new(model.CategoryView)

	if err := infrastructure.PostgresDB.NewSelect().Model(category).Where("_category.id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (categoryRepository *categoryRepository) GetById(ctx context.Context, id string) (*model.Category, error) {
	category := new(model.Category)

	if err := infrastructure.PostgresDB.NewSelect().Model(category).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (categoryRepository *categoryRepository) GetByName(ctx context.Context, name string) (*model.Category, error) {
	category := new(model.Category)

	if err := infrastructure.PostgresDB.NewSelect().Model(category).Where("name = ?", name).Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (categoryRepository *categoryRepository) Create(ctx context.Context, newCategory *model.Category) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newCategory).Returning("*").Exec(ctx)
	return err
}

func (categoryRepository *categoryRepository) Update(ctx context.Context, updatedCategory *model.Category) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedCategory).Where("id = ?", updatedCategory.Id).Exec(ctx)
	return err
}

func (categoryRepository *categoryRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Category{}).Where("id = ?", id).Exec(ctx)
	return err
}
