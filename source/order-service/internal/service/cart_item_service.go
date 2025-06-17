package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
)

type cartItemService struct {
	cartItemRepository repository.CartItemRepository
}

type CartItemService interface {
	// Main features
	GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) ([]dto.CartItemView, error)
	GetCartItemsByUserId(ctx context.Context, reqDTO *dto.GetCartItemsByUserIdRequest) ([]dto.CartItemView, error)
	CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error
	UpdateCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.UpdateCartItemByIdAndUserIdRequest) error
	DeleteCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.DeleteCartItemByIdAndUserIdRequest) error
}

func NewCartItemService(cartItemRepository repository.CartItemRepository) CartItemService {
	return &cartItemService{
		cartItemRepository: cartItemRepository,
	}
}

//
//
// Main features
// ######################################################################################

func (cartItemService *cartItemService) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) ([]dto.CartItemView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	cartItems, err := cartItemService.cartItemRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
	}

	// Missing->EnrichInformation

	return dto.ToListCartItemView(cartItems), nil
}

func (cartItemService *cartItemService) GetCartItemsByUserId(ctx context.Context, reqDTO *dto.GetCartItemsByUserIdRequest) ([]dto.CartItemView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	cartItems, err := cartItemService.cartItemRepository.GetByUserId(ctx, reqDTO.UserId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
	}

	// Missing->EnrichInformation

	return dto.ToListCartItemView(cartItems), nil
}

func (cartItemService *cartItemService) CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error {
	// Missing->ValidateProductId

	newCartItem := model.CartItem{
		UserId:    reqDTO.Body.UserId,
		ProductId: reqDTO.Body.ProductId,
		Quantity:  1,
	}
	if err := cartItemService.cartItemRepository.Create(ctx, &newCartItem); err != nil {
		return fmt.Errorf("insert cart item to postgresql failed: %s", err.Error())
	}

	return nil
}

func (cartItemService *cartItemService) UpdateCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.UpdateCartItemByIdAndUserIdRequest) error {
	foundCartItem, err := cartItemService.cartItemRepository.GetByIdAndUserId(ctx, reqDTO.Id, reqDTO.UserId)
	if err != nil {
		return fmt.Errorf("id or user id of cart item is not valid: %s", err.Error())
	}

	if reqDTO.Body.Quantity != nil {
		foundCartItem.Quantity = *reqDTO.Body.Quantity
	}

	if err := cartItemService.cartItemRepository.Update(ctx, foundCartItem); err != nil {
		return fmt.Errorf("update cart item on postgresql failed: %s", err.Error())
	}

	return nil
}

func (cartItemService *cartItemService) DeleteCartItemByIdAndUserId(ctx context.Context, reqDTO *dto.DeleteCartItemByIdAndUserIdRequest) error {
	if _, err := cartItemService.cartItemRepository.GetByIdAndUserId(ctx, reqDTO.Id, reqDTO.UserId); err != nil {
		return fmt.Errorf("id or user id of cart item is not valid")
	}

	if err := cartItemService.cartItemRepository.DeleteByIdAndUserId(ctx, reqDTO.Id, reqDTO.UserId); err != nil {
		return fmt.Errorf("delete cart item from postgresql failed: %s", err.Error())
	}

	return nil
}
