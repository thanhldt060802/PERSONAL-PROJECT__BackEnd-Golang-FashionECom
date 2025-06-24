package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

type CategoryService interface {
	// Main features
	GetAllCategories(ctx context.Context, reqDTO *dto.GetAllCategoriesRequest) ([]dto.CategoryView, error)
	GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*dto.CategoryView, error)
	CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error
	UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) error
	DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) error
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

func (categoryService *categoryService) GetAllCategories(ctx context.Context, reqDTO *dto.GetAllCategoriesRequest) ([]dto.CategoryView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	categories, err := categoryService.categoryRepository.GetAll(ctx, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query categories from postgresql failed: %s", err.Error())
	}

	return dto.ToListCategoryView(categories), nil
}

func (categoryService *categoryService) GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*dto.CategoryView, error) {
	foundCategory, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of category is not valid: %s", err.Error())
	}

	return dto.ToCategoryView(foundCategory), nil
}

func (categoryService *categoryService) CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error {
	if _, err := categoryService.categoryRepository.GetByName(ctx, reqDTO.Body.Name); err == nil {
		return fmt.Errorf("name of category is already exists")
	}

	newCategory := model.Category{
		Id:   uuid.New().String(),
		Name: reqDTO.Body.Name,
	}
	if err := categoryService.categoryRepository.Create(ctx, &newCategory); err != nil {
		return fmt.Errorf("insert category to postgresql failed: %s", err.Error())
	}

	return nil
}

func (categoryService *categoryService) UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) error {
	foundCategory, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of category is not valid: %s", err.Error())
	}

	if reqDTO.Body.Name != nil {
		if _, err := categoryService.categoryRepository.GetByName(ctx, *reqDTO.Body.Name); err == nil {
			return fmt.Errorf("name of category is already exists")
		}
		foundCategory.Name = *reqDTO.Body.Name
	}
	timeUpdate := time.Now().UTC()
	foundCategory.UpdatedAt = &timeUpdate

	if err := categoryService.categoryRepository.Update(ctx, foundCategory); err != nil {
		return fmt.Errorf("update category on postgresql failed: %s", err.Error())
	}

	return nil
}

func (categoryService *categoryService) DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) error {
	if _, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of category is not valid")
	}

	if err := categoryService.categoryRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete category from postgresql failed: %s", err.Error())
	}

	return nil
}
