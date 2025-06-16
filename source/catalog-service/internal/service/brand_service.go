package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"
)

type brandService struct {
	brandRepository repository.BrandRepository
}

type BrandService interface {
	// Main features
	GetBrands(ctx context.Context, reqDTO *dto.GetBrandsRequest) ([]dto.BrandView, error)
	GetBrandById(ctx context.Context, reqDTO *dto.GetBrandByIdRequest) (*dto.BrandView, error)
	CreateBrand(ctx context.Context, reqDTO *dto.CreateBrandRequest) error
	UpdateBrandById(ctx context.Context, reqDTO *dto.UpdateBrandByIdRequest) error
	DeleteBrandById(ctx context.Context, reqDTO *dto.DeleteBrandByIdRequest) error
}

func NewBrandService(brandRepository repository.BrandRepository) BrandService {
	return &brandService{
		brandRepository: brandRepository,
	}
}

//
//
// Main features
// ######################################################################################

func (brandService *brandService) GetBrands(ctx context.Context, redDTO *dto.GetBrandsRequest) ([]dto.BrandView, error) {
	sortFields := utils.ParseSorter(redDTO.SortBy)

	brands, err := brandService.brandRepository.Get(ctx, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query brands from postgresql failed: %s", err.Error())
	}

	return dto.ToListBrandView(brands), nil
}

func (brandService *brandService) GetBrandById(ctx context.Context, reqDTO *dto.GetBrandByIdRequest) (*dto.BrandView, error) {
	foundBrand, err := brandService.brandRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of brand is not valid: %s", err.Error())
	}

	return dto.ToBrandView(foundBrand), nil
}

func (brandService *brandService) CreateBrand(ctx context.Context, reqDTO *dto.CreateBrandRequest) error {
	if _, err := brandService.brandRepository.GetByName(ctx, reqDTO.Body.Name); err == nil {
		return fmt.Errorf("name of brand is already exists")
	}

	newBrand := model.Brand{
		Name:        reqDTO.Body.Name,
		Description: reqDTO.Body.Description,
	}
	if err := brandService.brandRepository.Create(ctx, &newBrand); err != nil {
		return fmt.Errorf("insert brand to postgresql failed: %s", err.Error())
	}

	return nil
}

func (brandService *brandService) UpdateBrandById(ctx context.Context, reqDTO *dto.UpdateBrandByIdRequest) error {
	foundBrand, err := brandService.brandRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of brand is not valid: %s", err.Error())
	}

	if reqDTO.Body.Name != nil {
		if _, err := brandService.brandRepository.GetByName(ctx, *reqDTO.Body.Name); err == nil {
			return fmt.Errorf("name of brand is already exists")
		}
		foundBrand.Name = *reqDTO.Body.Name
	}
	if reqDTO.Body.Description != nil {
		foundBrand.Description = *reqDTO.Body.Description
	}
	timeUpdate := time.Now().UTC()
	foundBrand.UpdatedAt = &timeUpdate

	if err := brandService.brandRepository.Update(ctx, foundBrand); err != nil {
		return fmt.Errorf("update brand on postgresql failed: %s", err.Error())
	}

	return nil
}

func (brandService *brandService) DeleteBrandById(ctx context.Context, reqDTO *dto.DeleteBrandByIdRequest) error {
	if _, err := brandService.brandRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of brand is not valid")
	}

	if err := brandService.brandRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete brand from postgresql failed: %s", err.Error())
	}

	return nil
}
