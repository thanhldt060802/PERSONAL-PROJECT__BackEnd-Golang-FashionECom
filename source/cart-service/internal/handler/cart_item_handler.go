package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CartItemHandler struct {
	cartItemService   service.CartItemService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewCartItemHandler(api huma.API, cartItemService service.CartItemService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *CartItemHandler {
	cartItemHandler := &CartItemHandler{
		cartItemService:   cartItemService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	//
	//
	// Main features
	// ######################################################################################

	// Get all account cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Get all account cart items.",
		Tags:        []string{"Account Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.GetAllAccountCartItems)

	// Create account cart item
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Create account cart item.",
		Tags:        []string{"Account Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.CreateAccountCartItem)

	// Update account cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/my-cart-items/id/{id}",
		Summary:     "/my-cart-items/id/{id}",
		Description: "Update account cart item by id.",
		Tags:        []string{"Account Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.UpdateAccountCartItemById)

	// Delete account cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/my-cart-items/id/{id}",
		Summary:     "/my-cart-items/id/{id}",
		Description: "Delete account cart item by id.",
		Tags:        []string{"Account Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.DeleteCartItemById)

	return cartItemHandler
}

//
//
// Main features
// ######################################################################################

func (cartItemHandler *CartItemHandler) GetAllAccountCartItems(ctx context.Context, _ *struct{}) (*dto.BodyResponse[[]dto.CartItemView], error) {
	convertReqDTO := &dto.GetAllCartItemsByUserIdRequest{}
	convertReqDTO.UserId = ctx.Value("user_id").(int64)

	cartItems, err := cartItemHandler.cartItemService.GetAllCartItemsByUserId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[[]dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get all cart items successful"
	res.Body.Data = cartItems
	return res, nil
}

func (cartItemHandler *CartItemHandler) CreateAccountCartItem(ctx context.Context, reqDTO *dto.CreateAccountCartItemRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateCartItemRequest{}
	convertReqDTO.Body.UserId = ctx.Value("user_id").(int64)
	convertReqDTO.Body.ProductId = reqDTO.Body.ProductId

	if err := cartItemHandler.cartItemService.CreateCartItem(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create account cart item failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create account cart item successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) UpdateAccountCartItemById(ctx context.Context, reqDTO *dto.UpdateAccountCartItemByIdRequest) (*dto.SuccessResponse, error) {
	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, reqDTO.Id)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	if foundCartItem.UserId != ctx.Value("user_id").(int64) {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusForbidden
		res.Code = "ERR_FORBIDDEN"
		res.Message = "Update cart item by id failed"
		res.Details = []string{"id of cart item is not owned"}
		return nil, res
	}

	convertReqDTO := &dto.UpdateCartItemByIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.Body.Quantity = reqDTO.Body.Quantity

	if err := cartItemHandler.cartItemService.UpdateCartItemById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update account cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update account cart item by id successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteAccountCartItemByIdRequest) (*dto.SuccessResponse, error) {
	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, reqDTO.Id)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	if foundCartItem.UserId != ctx.Value("user_id").(int64) {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusForbidden
		res.Code = "ERR_FORBIDDEN"
		res.Message = "Delete cart item by id failed"
		res.Details = []string{"id of cart item is not owned"}
		return nil, res
	}

	convertReqDTO := &dto.DeleteCartItemByIdRequest{}
	convertReqDTO.Id = reqDTO.Id

	if err := cartItemHandler.cartItemService.DeleteCartItemById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete cart item by id successful"
	return res, nil
}
