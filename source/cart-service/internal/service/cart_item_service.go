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

type cartItemService struct {
	cartItemRepository repository.CartItemRepository
}

type CartItemService interface {
	// Main features
	GetAllCartItemsByUserId(ctx context.Context, reqDTO *dto.GetAllCartItemsByUserIdRequest) ([]dto.CartItemView, error)
	GetCartItemById(ctx context.Context, id int64) (*dto.CartItemView, error)
	CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error
	UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemByIdRequest) error
	DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemByIdRequest) error
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

func (cartItemService *cartItemService) GetAllCartItemsByUserId(ctx context.Context, reqDTO *dto.GetAllCartItemsByUserIdRequest) ([]dto.CartItemView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	cartItems, err := cartItemService.cartItemRepository.GetAllByUserId(ctx, reqDTO.UserId, &reqDTO.Offset, &reqDTO.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
	}

	return dto.ToListCartItemView(cartItems), nil
}

func (cartItemService *cartItemService) GetCartItemById(ctx context.Context, id int64) (*dto.CartItemView, error) {
	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("id of cart item is not valid: %s", err.Error())
	}

	return dto.ToCartItemView(foundCartItem), nil
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

func (cartItemService *cartItemService) UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemByIdRequest) error {
	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of cartItem is not valid: %s", err.Error())
	}

	if reqDTO.Body.Quantity != nil {
		foundCartItem.Quantity = *reqDTO.Body.Quantity
	}
	foundCartItem.UpdatedAt = time.Now().UTC()

	if err := cartItemService.cartItemRepository.Update(ctx, foundCartItem); err != nil {
		return fmt.Errorf("update cart item on postgresql failed: %s", err.Error())
	}

	return nil
}

func (cartItemService *cartItemService) DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemByIdRequest) error {
	if _, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of cart item is not valid")
	}

	if err := cartItemService.cartItemRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete cart item from postgresql failed: %s", err.Error())
	}

	return nil
}
