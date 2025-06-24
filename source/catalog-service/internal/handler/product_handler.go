package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type ProductHandler struct {
	productService    service.ProductService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewProductHandler(api huma.API, productService service.ProductService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *ProductHandler {
	productHandler := &ProductHandler{
		productService:    productService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	//
	//
	// For admin + customer
	// ######################################################################################

	// Get product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Get product by id.",
		Tags:        []string{"Product"},
	}, productHandler.GetProductById)

	//
	//
	// For only admin
	// ######################################################################################

	// Create product
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/products",
		Summary:     "/products",
		Description: "Create product.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, productHandler.CreateProduct)

	// Update product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Update product by id.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, productHandler.UpdateProductById)

	// Delete product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Delete product by id.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, productHandler.DeleteProductById)

	//
	//
	// For admin + customer
	// ######################################################################################

	// Get products
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products",
		Summary:     "/products",
		Description: "Get products.",
		Tags:        []string{"Product"},
	}, productHandler.GetProducts)

	return productHandler
}

func (productHandler *ProductHandler) GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*dto.BodyResponse[dto.ProductView], error) {
	foundProduct, err := productHandler.productService.GetProductById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get product by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get product by id successful"
	res.Body.Data = *foundProduct
	return res, nil
}

func (productHandler *ProductHandler) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.CreateProduct(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create product failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create product successful"
	return res, nil
}

func (productHandler *ProductHandler) UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductByIdRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.UpdateProductById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update product by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update product by id successful"
	return res, nil
}

func (productHandler *ProductHandler) DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductByIdRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.DeleteProductById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete product by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete product by id successful"
	return res, nil
}

func (productHandler *ProductHandler) GetProducts(ctx context.Context, reqDTO *dto.GetProductsRequest) (*dto.PaginationBodyResponseList[dto.ProductView], error) {
	products, err := productHandler.productService.GetProducts(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get products failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get products successful"
	res.Body.Data = products
	res.Body.Total = len(products)
	return res, nil
}
