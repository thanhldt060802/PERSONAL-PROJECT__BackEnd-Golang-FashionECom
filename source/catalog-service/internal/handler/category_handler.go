package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CategoryHandler struct {
	categoryService   service.CategoryService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewCategoryHandler(api huma.API, categoryService service.CategoryService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *CategoryHandler {
	categoryHandler := &CategoryHandler{
		categoryService:   categoryService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	//
	//
	// For admin + customer
	// ######################################################################################

	// Get all categories
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/categories/all",
		Summary:     "/categories/all",
		Description: "Get all categories.",
		Tags:        []string{"Category"},
	}, categoryHandler.GetAllCategories)

	// Get category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Get category by id.",
		Tags:        []string{"Category"},
	}, categoryHandler.GetCategoryById)

	//
	//
	// For only admin
	// ######################################################################################

	// Create category
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/categories",
		Summary:     "/categories",
		Description: "Create category.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, categoryHandler.CreateCategory)

	// Update category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Update category by id.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, categoryHandler.UpdateCategoryById)

	// Delete category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Delete category by id.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, categoryHandler.DeleteCategoryById)

	return categoryHandler
}

func (categoryHandler *CategoryHandler) GetAllCategories(ctx context.Context, reqDTO *dto.GetCategoriesRequest) (*dto.PaginationBodyResponseList[dto.CategoryView], error) {
	categories, err := categoryHandler.categoryService.GetCategories(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all categories failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.CategoryView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get all categories successful"
	res.Body.Data = categories
	res.Body.Total = len(categories)
	return res, nil
}

func (categoryHandler *CategoryHandler) GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*dto.BodyResponse[dto.CategoryView], error) {
	foundCategory, err := categoryHandler.categoryService.GetCategoryById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get category by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.CategoryView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get category by id successful"
	res.Body.Data = *foundCategory
	return res, nil
}

func (categoryHandler *CategoryHandler) CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categoryService.CreateCategory(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create category failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create category successful"
	return res, nil
}

func (categoryHandler *CategoryHandler) UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categoryService.UpdateCategoryById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update category by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update category by id successful"
	return res, nil
}

func (categoryHandler *CategoryHandler) DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categoryService.DeleteCategoryById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete category by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete category by id successful"
	return res, nil
}
