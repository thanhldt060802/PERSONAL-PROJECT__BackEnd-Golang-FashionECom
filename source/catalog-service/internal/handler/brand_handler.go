package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type BrandHandler struct {
	brandService      service.BrandService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewBrandHandler(api huma.API, brandService service.BrandService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *BrandHandler {
	brandHandler := &BrandHandler{
		brandService:      brandService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	//
	//
	// For admin + customer
	// ######################################################################################

	// Get all brands
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/brands/all",
		Summary:     "/brands/all",
		Description: "Get all brands.",
		Tags:        []string{"Brand"},
	}, brandHandler.GetBrands)

	// Get brand by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/brands/id/{id}",
		Summary:     "/brands/id/{id}",
		Description: "Get brand by id.",
		Tags:        []string{"Brand"},
	}, brandHandler.GetBrandById)

	//
	//
	// For only admin
	// ######################################################################################

	// Create brand
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/brands",
		Summary:     "/brands",
		Description: "Create brand.",
		Tags:        []string{"Brand"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, brandHandler.CreateBrand)

	// Update brand by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/brands/id/{id}",
		Summary:     "/brands/id/{id}",
		Description: "Update brand by id.",
		Tags:        []string{"Brand"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, brandHandler.UpdateBrandById)

	// Delete brand by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/brands/id/{id}",
		Summary:     "/brands/id/{id}",
		Description: "Delete brand by id.",
		Tags:        []string{"Brand"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, brandHandler.DeleteBrandById)

	return brandHandler
}

func (brandHandler *BrandHandler) GetBrands(ctx context.Context, reqDTO *dto.GetBrandsRequest) (*dto.PaginationBodyResponseList[dto.BrandView], error) {
	brands, err := brandHandler.brandService.GetBrands(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all brands failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.BrandView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get all brands successful"
	res.Body.Data = brands
	res.Body.Total = len(brands)
	return res, nil
}

func (brandHandler *BrandHandler) GetBrandById(ctx context.Context, reqDTO *dto.GetBrandByIdRequest) (*dto.BodyResponse[dto.BrandView], error) {
	foundBrand, err := brandHandler.brandService.GetBrandById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get brand by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.BrandView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get brand by id successful"
	res.Body.Data = *foundBrand
	return res, nil
}

func (brandHandler *BrandHandler) CreateBrand(ctx context.Context, reqDTO *dto.CreateBrandRequest) (*dto.SuccessResponse, error) {
	if err := brandHandler.brandService.CreateBrand(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create brand failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create brand successful"
	return res, nil
}

func (brandHandler *BrandHandler) UpdateBrandById(ctx context.Context, reqDTO *dto.UpdateBrandByIdRequest) (*dto.SuccessResponse, error) {
	if err := brandHandler.brandService.UpdateBrandById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update brand by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update brand by id successful"
	return res, nil
}

func (brandHandler *BrandHandler) DeleteBrandById(ctx context.Context, reqDTO *dto.DeleteBrandByIdRequest) (*dto.SuccessResponse, error) {
	if err := brandHandler.brandService.DeleteBrandById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete brand by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete brand by id successful"
	return res, nil
}
