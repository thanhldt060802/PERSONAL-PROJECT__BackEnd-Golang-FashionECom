package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/model"
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

	// Get cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items",
		Summary:     "/cart-items",
		Description: "Get cart items.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItems)

	// Create cart item
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/cart-items",
		Summary:     "/cart-items",
		Description: "Create cart item.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.CreateCartItem)

	// Update cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/cart-items/id/{id}",
		Summary:     "/cart-items/id/{id}",
		Description: "Update cart item by id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.UpdateCartItemById)

	// Delete cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/cart-items/id/{id}",
		Summary:     "/cart-items/id/{id}",
		Description: "Delete cart item by id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, cartItemHandler.DeleteCartItemById)

	// Get my cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Get my cart items.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.GetMyCartItems)

	// Create my cart item
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Create my cart item.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.CreateMyCartItem)

	// Update my cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/cart-items/id/{id}",
		Summary:     "/cart-items/id/{id}",
		Description: "Update my cart item by id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.UpdateMyCartItemById)

	// Delete my cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/my-cart-items/id/{id}",
		Summary:     "/my-cart-items/id/{id}",
		Description: "Delete my cart item by id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, cartItemHandler.DeleteMyCartItemById)

	return cartItemHandler
}

func (cartItemHandler *CartItemHandler) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) (*dto.PaginationBodyResponseList[*model.CartItemView], error) {
	cartItems, err := cartItemHandler.cartItemService.GetCartItems(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*model.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items successful"
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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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

func (cartItemHandler *CartItemHandler) GetMyCartItems(ctx context.Context, reqDTO *dto.GetMyCartItemsRequest) (*dto.PaginationBodyResponseList[*model.CartItemView], error) {
	convertReqDTO := &dto.GetCartItemsRequest{}
	convertReqDTO.Offset = reqDTO.Offset
	convertReqDTO.Limit = reqDTO.Limit
	convertReqDTO.SortBy = reqDTO.SortBy
	convertReqDTO.UserId = ctx.Value("user_id").(string)

	cartItems, err := cartItemHandler.cartItemService.GetCartItems(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get my cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*model.CartItemView]{}
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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update my cart item by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	convertReqDTO := &dto.UpdateCartItemByIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.Body.Quantity = reqDTO.Body.Quantity
	convertReqDTO.UserId = ctx.Value("user_id").(string)

	if err := cartItemHandler.cartItemService.UpdateCartItemById(ctx, convertReqDTO); err != nil {
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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete my cart item by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	convertReqDTO := &dto.DeleteCartItemByIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.UserId = ctx.Value("user_id").(string)

	if err := cartItemHandler.cartItemService.DeleteCartItemById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete my cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete my cart item by id successful"
	return res, nil
}
