package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/google/uuid"
)

type brandRepository struct {
}

type BrandRepository interface {
	// Main features
	Get(ctx context.Context, offset *int, limit *int, sortFields []utils.SortField) ([]model.Brand, error)
	GetById(ctx context.Context, id string) (*model.Brand, error)
	GetByName(ctx context.Context, name string) (*model.Brand, error)
	Create(ctx context.Context, newBrand *model.Brand) error
	Update(ctx context.Context, updatedBrand *model.Brand) error
	DeleteById(ctx context.Context, id string) error
}

func NewBrandRepository() BrandRepository {
	return &brandRepository{}
}

//
//
// Main features
// ######################################################################################

func (brandRepository *brandRepository) Get(ctx context.Context, offset *int, limit *int, sortFields []utils.SortField) ([]model.Brand, error) {
	var brands []model.Brand

	query := infrastructure.PostgresDB.NewSelect().Model(&brands)

	if offset != nil {
		query = query.Offset(*offset)
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return brands, nil
}

func (brandRepository *brandRepository) GetById(ctx context.Context, id string) (*model.Brand, error) {
	var brand model.Brand

	if err := infrastructure.PostgresDB.NewSelect().Model(&brand).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return &brand, nil
}

func (brandRepository *brandRepository) GetByName(ctx context.Context, name string) (*model.Brand, error) {
	var brand model.Brand

	if err := infrastructure.PostgresDB.NewSelect().Model(&brand).Where("name = ?", name).Scan(ctx); err != nil {
		return nil, err
	}

	return &brand, nil
}

func (brandRepository *brandRepository) Create(ctx context.Context, newBrand *model.Brand) error {
	newBrand.Id = uuid.New().String()
	_, err := infrastructure.PostgresDB.NewInsert().Model(newBrand).Exec(ctx)
	return err
}

func (brandRepository *brandRepository) Update(ctx context.Context, updatedBrand *model.Brand) error {
	_, err := infrastructure.PostgresDB.NewUpdate().Model(updatedBrand).Where("id = ?", updatedBrand.Id).Exec(ctx)
	return err
}

func (brandRepository *brandRepository) DeleteById(ctx context.Context, id string) error {
	_, err := infrastructure.PostgresDB.NewDelete().Model(&model.Brand{}).Where("id = ?", id).Exec(ctx)
	return err
}
