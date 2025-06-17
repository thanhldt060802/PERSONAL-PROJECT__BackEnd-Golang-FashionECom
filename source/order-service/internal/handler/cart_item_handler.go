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

	//
	//
	// Main features
	// ######################################################################################

	// Get account cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Get account cart items.",
		Tags:        []string{"Account Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.GetAccountCartItems)

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

func (cartItemHandler *CartItemHandler) UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemByIdRequest) (*dto.SuccessResponse, error) {
	if err := cartItemHandler.cartItemService.UpdateCartItemById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update cart item by id successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemByIdRequest) (*dto.SuccessResponse, error) {
	if err := cartItemHandler.cartItemService.DeleteCartItemById(ctx, reqDTO); err != nil {
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

func (cartItemHandler *CartItemHandler) UpdateAccountCartItemById(ctx context.Context, reqDTO *dto.UpdateAccountCartItemByIdRequest) (*dto.SuccessResponse, error) {
	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, reqDTO.Id)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update account cart item by id failed"
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
