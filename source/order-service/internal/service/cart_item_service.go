package service

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
)

type cartItemService struct {
	cartItemRepository repository.CartItemRepository
}

type CartItemService interface {
	GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) ([]dto.CartItemView, error)
	CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error
	UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemByIdRequest) error
	DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemByIdRequest) error
}

func NewCartItemService(cartItemRepository repository.CartItemRepository) CartItemService {
	return &cartItemService{
		cartItemRepository: cartItemRepository,
	}
}

func (cartItemService *cartItemService) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsRequest) ([]dto.CartItemView, error) {
	if infrastructure.CatalogServiceGRPCClient != nil {
		sortFields := utils.ParseSorter(reqDTO.SortBy)

		var foundCartItems []model.CartItem
		if reqDTO.UserId == "" {
			cartItems, err := cartItemService.cartItemRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
			if err != nil {
				return nil, fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
			}
			foundCartItems = cartItems
		} else {
			cartItems, err := cartItemService.cartItemRepository.GetByUserId(ctx, reqDTO.UserId, reqDTO.Offset, reqDTO.Limit, sortFields)
			if err != nil {
				return nil, fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
			}
			foundCartItems = cartItems
		}

		ids := make([]string, len(foundCartItems))
		for i, cartItem := range foundCartItems {
			ids[i] = cartItem.ProductId
		}

		convertReqDTO := &catalogservicepb.GetProductsByListIdRequest{}
		convertReqDTO.Ids = ids
		grpcRes, err := infrastructure.CatalogServiceGRPCClient.GetProductsByListId(ctx, convertReqDTO)
		if err != nil {
			return nil, fmt.Errorf("get products from catalog-service failed: %s", err.Error())
		}

		productProtos := grpcRes.Products
		foundCartItemExtraInfos := make([]dto.CartItemExtraInfo, len(productProtos))
		for i, productProto := range productProtos {
			foundCartItemExtraInfos[i].Name = productProto.Name
			foundCartItemExtraInfos[i].Sex = productProto.Sex
			foundCartItemExtraInfos[i].Price = productProto.Price
			foundCartItemExtraInfos[i].DiscountPercentage = productProto.DiscountPercentage
			foundCartItemExtraInfos[i].ImageURL = productProto.ImageUrl

			foundCartItemExtraInfos[i].CategoryId = productProto.CategoryId
			foundCartItemExtraInfos[i].CategoryName = productProto.CategoryName
			foundCartItemExtraInfos[i].BrandId = productProto.BrandId
			foundCartItemExtraInfos[i].BrandName = productProto.BrandName
		}

		return dto.ToListCartItemView(foundCartItems, foundCartItemExtraInfos), nil
	} else {
		return nil, fmt.Errorf("catalog-service is not running")
	}
}

func (cartItemService *cartItemService) CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error {
	if infrastructure.CatalogServiceGRPCClient != nil {
		// Missing->ValidateUser

		convertReqDTO := &catalogservicepb.GetProductsByListIdRequest{}
		convertReqDTO.Ids = []string{reqDTO.Body.ProductId}
		grpcRes, err := infrastructure.CatalogServiceGRPCClient.GetProductsByListId(ctx, convertReqDTO)
		if err != nil {
			return fmt.Errorf("get products from catalog-service failed: %s", err.Error())
		}

		if len(grpcRes.Products) == 0 {
			return fmt.Errorf("id of product is not valid")
		}

		newCartItem := model.CartItem{
			UserId:    reqDTO.Body.UserId,
			ProductId: reqDTO.Body.ProductId,
			Quantity:  1,
		}
		if err := cartItemService.cartItemRepository.Create(ctx, &newCartItem); err != nil {
			return fmt.Errorf("insert cart item to postgresql failed: %s", err.Error())
		}

		return nil
	} else {
		return fmt.Errorf("catalog-service is not running")
	}
}

func (cartItemService *cartItemService) UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemByIdRequest) error {
	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of cart item is not valid: %s", err.Error())
	}

	if reqDTO.UserId != "" && reqDTO.UserId != foundCartItem.UserId {
		return fmt.Errorf("id of cart item is not valid: no permission")
	}

	if reqDTO.Body.Quantity != nil {
		foundCartItem.Quantity = *reqDTO.Body.Quantity
	}

	if err := cartItemService.cartItemRepository.Update(ctx, foundCartItem); err != nil {
		return fmt.Errorf("update cart item on postgresql failed: %s", err.Error())
	}

	return nil
}

func (cartItemService *cartItemService) DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemByIdRequest) error {
	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of cart item is not valid: %s", err.Error())
	}

	if reqDTO.UserId != "" && reqDTO.UserId != foundCartItem.UserId {
		return fmt.Errorf("id of cart item is not valid: no permission")
	}

	if err := cartItemService.cartItemRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete cart item from postgresql failed: %s", err.Error())
	}

	return nil
}
