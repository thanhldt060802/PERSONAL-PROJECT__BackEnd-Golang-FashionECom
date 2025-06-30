package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type brandRepository struct {
}

type BrandRepository interface {
	GetAllViews(ctx context.Context, sortFields []*utils.SortField) ([]*model.BrandView, error)
	GetViewById(ctx context.Context, id string) (*model.BrandView, error)

	GetById(ctx context.Context, id string) (*model.Brand, error)
	GetByName(ctx context.Context, name string) (*model.Brand, error)
	Create(ctx context.Context, newBrand *model.Brand) error
	Update(ctx context.Context, updatedBrand *model.Brand) error
	DeleteById(ctx context.Context, id string) error
}

func NewBrandRepository() BrandRepository {
	return &brandRepository{}
}

func (brandRepository *brandRepository) GetAllViews(ctx context.Context, sortFields []*utils.SortField) ([]*model.BrandView, error) {
	var brands []*model.BrandView

	query := infrastructure.PostgresDB.NewSelect().Model(&brands)

	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("_brand.%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return brands, nil
}

func (brandRepository *brandRepository) GetViewById(ctx context.Context, id string) (*model.BrandView, error) {
	brand := new(model.BrandView)

	query := infrastructure.PostgresDB.NewSelect().Model(brand).Where("_brand.id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return brand, nil
}

func (brandRepository *brandRepository) GetById(ctx context.Context, id string) (*model.Brand, error) {
	brand := new(model.Brand)

	query := infrastructure.PostgresDB.NewSelect().Model(brand).Where("id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return brand, nil
}

func (brandRepository *brandRepository) GetByName(ctx context.Context, name string) (*model.Brand, error) {
	brand := new(model.Brand)

	query := infrastructure.PostgresDB.NewSelect().Model(brand).Where("name = ?", name)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return brand, nil
}

func (brandRepository *brandRepository) Create(ctx context.Context, newBrand *model.Brand) error {
	_, err := infrastructure.PostgresDB.NewInsert().Model(newBrand).Returning("*").Exec(ctx)
	return err
}

func (brandRepository *brandRepository) Update(ctx context.Context, updatedBrand *model.Brand) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedBrand).Exec(ctx)
	return err
}

func (brandRepository *brandRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Brand{}).Where("id = ?", id).Exec(ctx)
	return err
}
