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

	// Get cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items",
		Summary:     "/cart-items",
		Description: "Get cart items.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItems)

	// Get cart items by user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items/user-id/{user_id}",
		Summary:     "/cart-items/user-id/{user_id}",
		Description: "Get cart items by user id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItemsByUserId)

	// Create cart item
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/cart-items",
		Summary:     "/cart-items",
		Description: "Create cart item.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.CreateCartItem)

	// Update cart item by id and user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/cart-items/id/{id}/user-id/{user_id}",
		Summary:     "/cart-items/id/{id}/user-id/{user_id}",
		Description: "Update cart item by id and user id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.UpdateCartItemByIdAndUserId)

	// Delte cart item by id and user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/cart-items/id/{id}/user-id/{user_id}",
		Summary:     "/cart-items/id/{id}/user-id/{user_id}",
		Description: "Update cart item by id and user id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.DeleteCartItemByIdAndUserId)

	// Get my cart items

	// Create my cart item

	// Update my cart item

	// Delete my cart item

	//
	//
	// Main features
	// ######################################################################################

	return cartItemHandler
}

//
//
// Main features
// ######################################################################################

func (cartItemHandler *CartItemHandler) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	cartItems, err := cartItemHandler.cartItemService.GetCartItems(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items successful"
	res.Body.Data = cartItems
	res.Body.Total = len(cartItems)
	return res, nil
}

func (cartItemHandler *CartItemHandler) GetCartItemsByUserId(ctx context.Context, reqDTO *dto.GetCartItemsByUserIdRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	cartItems, err := cartItemHandler.cartItemService.GetCartItemsByUserId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items by user id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items by user id successful"
	res.Body.Data = cartItems
	res.Body.Total = len(cartItems)
	return res, nil
}

func (cartItemHandler *CartItemHandler) CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) (*dto.SuccessResponse, error) {
	if err := cartItemHandler.cartItemService.CreateCartItem(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create cart item failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create cart item successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) UpdateCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.UpdateCartItemByIdAndUserIdRequest) (*dto.SuccessResponse, error) {
	if err := cartItemHandler.cartItemService.UpdateCartItemByIdAndUserId(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item by id and user id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update cart item by id and user id successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) DeleteCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.DeleteCartItemByIdAndUserIdRequest) (*dto.SuccessResponse, error) {
	if err := cartItemHandler.cartItemService.DeleteCartItemByIdAndUserId(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item by id and user id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete cart item by id and user id successful"
	return res, nil
}

//
//
// Extra features
// ######################################################################################

func (cartItemHandler *CartItemHandler) GetMyCartItems(ctx context.Context, reqDTO *dto.GetMyCartItemsRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	convertReqDTO := &dto.GetCartItemsByUserIdRequest{}
	convertReqDTO.UserId = ctx.Value("user_id").(string)
	convertReqDTO.Offset = reqDTO.Offset
	convertReqDTO.Limit = reqDTO.Limit
	convertReqDTO.SortBy = reqDTO.SortBy

	cartItems, err := cartItemHandler.cartItemService.GetCartItemsByUserId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get my cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get my cart items successful"
	res.Body.Data = cartItems
	res.Body.Total = len(cartItems)
	return res, nil
}

func (cartItemHandler *CartItemHandler) CreateMyCartItem(ctx context.Context, reqDTO *dto.CreateMyCartItemRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateCartItemRequest{}
	convertReqDTO.Body.UserId = ctx.Value("user_id").(string)
	convertReqDTO.Body.ProductId = reqDTO.Body.ProductId

	if err := cartItemHandler.cartItemService.CreateCartItem(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create my cart item failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create my cart item successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) UpdateMyCartItemById(ctx context.Context, reqDTO *dto.UpdateMyCartItemByIdRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.UpdateCartItemByIdAndUserIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.UserId = ctx.Value("user_id").(string)
	convertReqDTO.Body.Quantity = reqDTO.Body.Quantity

	if err := cartItemHandler.cartItemService.UpdateCartItemByIdAndUserId(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update my cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update my cart item by id successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) DeleteMyCartItemById(ctx context.Context, reqDTO *dto.DeleteMyCartItemByIdRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.DeleteCartItemByIdAndUserIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.UserId = ctx.Value("user_id").(string)

	if err := cartItemHandler.cartItemService.DeleteCartItemByIdAndUserId(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete my cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete cart item by id successful"
	return res, nil
}
