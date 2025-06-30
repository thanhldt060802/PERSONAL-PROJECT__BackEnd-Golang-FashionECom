package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"time"

	"github.com/google/uuid"
)

type invoiceService struct {
	invoiceRepository  repository.InvoiceRepository
	cartItemRepository repository.CartItemRepository
}

type InvoiceService interface {
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*model.InvoiceView, error)
	CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllInvoices(ctx context.Context) ([]*model.InvoiceView, error)

	// Elasticsearch integration features
	GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) ([]*model.InvoiceView, error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, cartItemRepository repository.CartItemRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository:  invoiceRepository,
		cartItemRepository: cartItemRepository,
	}
}

func (invoiceService *invoiceService) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*model.InvoiceView, error) {
	foundInvoice, err := invoiceService.invoiceRepository.GetViewById(ctx, reqDTO.Id, true)
	if err != nil {
		return nil, fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	return foundInvoice, nil
}

func (invoiceService *invoiceService) CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error {
	newInvoice := &model.Invoice{
		Id:          uuid.New().String(),
		UserId:      reqDTO.Body.UserId,
		TotalAmount: 0,
		Status:      "CREATED",
	}

	if len(reqDTO.Body.InvoiceDetails) != 0 {
		newInvoiceDetails := []*model.InvoiceDetail{}
		for _, invoiceDetail := range reqDTO.Body.InvoiceDetails {
			newInvoiceDetails = append(newInvoiceDetails, &model.InvoiceDetail{
				Id:                 uuid.New().String(),
				InvoiceId:          newInvoice.Id,
				ProductId:          invoiceDetail.ProductId,
				Price:              invoiceDetail.Price,
				DiscountPercentage: invoiceDetail.DiscountPercentage,
				Quantity:           invoiceDetail.Quantity,
				TotalPrice:         invoiceDetail.TotalPrice,
			})

			newInvoice.TotalAmount += invoiceDetail.TotalPrice
		}

		if err := invoiceService.invoiceRepository.Create(ctx, newInvoice, newInvoiceDetails); err != nil {
			return fmt.Errorf("insert invoice to postgresql failed: %s", err.Error())
		}
	} else {
		cartItems, err := invoiceService.cartItemRepository.GetAllViewsByUserId(ctx, reqDTO.Body.UserId)
		if err != nil {
			return fmt.Errorf("query cart items from postgresql failed: %s", err.Error())
		}
		if len(cartItems) == 0 {
			return fmt.Errorf("cart items is empty")
		}
		for _, cartItem := range cartItems {
			totalPrice := int64(float64(cartItem.ProductPrice) * float64(100-cartItem.ProductDiscountPercentage) / 100 * float64(cartItem.Quantity))
			reqDTO.Body.InvoiceDetails = append(reqDTO.Body.InvoiceDetails, dto.InvoiceDetail{
				ProductId:          cartItem.ProductId,
				Price:              cartItem.ProductPrice,
				DiscountPercentage: cartItem.ProductDiscountPercentage,
				Quantity:           cartItem.Quantity,
				TotalPrice:         totalPrice,
			})
		}

		newInvoiceDetails := []*model.InvoiceDetail{}
		for _, invoiceDetail := range reqDTO.Body.InvoiceDetails {
			newInvoiceDetails = append(newInvoiceDetails, &model.InvoiceDetail{
				Id:                 uuid.New().String(),
				InvoiceId:          newInvoice.Id,
				ProductId:          invoiceDetail.ProductId,
				Price:              invoiceDetail.Price,
				DiscountPercentage: invoiceDetail.DiscountPercentage,
				Quantity:           invoiceDetail.Quantity,
				TotalPrice:         invoiceDetail.TotalPrice,
			})

			newInvoice.TotalAmount += invoiceDetail.TotalPrice
		}

		if err := invoiceService.invoiceRepository.Create(ctx, newInvoice, newInvoiceDetails); err != nil {
			return fmt.Errorf("insert invoice to postgresql failed: %s", err.Error())
		}

		if err := invoiceService.cartItemRepository.DeleteByUserId(ctx, reqDTO.Body.UserId); err != nil {
			return fmt.Errorf("delete cart items from postgresql failed: %s", err.Error())
		}
	}

	{
		convertReqDTO := &catalogservicepb.UpdateProductStocksByListInvoiceDetailRequest{}
		convertReqDTO.InvoiceDetails = make([]*catalogservicepb.InvoiceDetail, len(reqDTO.Body.InvoiceDetails))
		for i := range reqDTO.Body.InvoiceDetails {
			convertReqDTO.InvoiceDetails[i].ProductId = reqDTO.Body.InvoiceDetails[i].ProductId
			convertReqDTO.InvoiceDetails[i].Quantity = reqDTO.Body.InvoiceDetails[i].Quantity
		}
		_, err := infrastructure.CatalogServiceGRPCClient.UpdateProductStocksByListInvoiceDetail(ctx, convertReqDTO)
		if err != nil {
			return fmt.Errorf("update products from catalog-service failed: %s", err.Error())
		}
	}

	newInvoiceView, _ := invoiceService.invoiceRepository.GetViewById(ctx, newInvoice.Id, false)
	payload, _ := json.Marshal(newInvoiceView)
	if err := infrastructure.RedisClient.Publish(ctx, "order-service.created-invoice", payload).Err(); err != nil {
		return fmt.Errorf("pulish event order-service.created-invoice failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	if reqDTO.UserId != "" && reqDTO.UserId != foundInvoice.UserId {
		return fmt.Errorf("id of invoice is not valid: no permission")
	}

	if reqDTO.Body.Status != nil {
		foundInvoice.Status = *reqDTO.Body.Status
	}
	timeUpdate := time.Now().UTC()
	foundInvoice.UpdatedAt = &timeUpdate

	if err := invoiceService.invoiceRepository.Update(ctx, foundInvoice); err != nil {
		return fmt.Errorf("update invoice on postgresql failed: %s", err.Error())
	}

	updatedInvoiceView, _ := invoiceService.invoiceRepository.GetViewById(ctx, foundInvoice.Id, false)
	payload, _ := json.Marshal(updatedInvoiceView)
	if err := infrastructure.RedisClient.Publish(ctx, "order-service.updated-invoice", payload).Err(); err != nil {
		return fmt.Errorf("pulish event order-service.updated-invoice failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	if reqDTO.UserId != "" && reqDTO.UserId != foundInvoice.UserId {
		return fmt.Errorf("id of invoice is not valid: no permission")
	}

	if err := invoiceService.invoiceRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete invoice from postgresql failed: %s", err.Error())
	}

	if err := infrastructure.RedisClient.Publish(ctx, "order-service.deleted-invoice", reqDTO.Id).Err(); err != nil {
		return fmt.Errorf("pulish event order-service.deleted-invoice failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) GetAllInvoices(ctx context.Context) ([]*model.InvoiceView, error) {
	foundInvoices, err := invoiceService.invoiceRepository.GetAllViews(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("query invoices from postgresql failed: %s", err.Error())
	}

	return foundInvoices, nil
}

func (invoiceService *invoiceService) GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) ([]*model.InvoiceView, error) {
	if infrastructure.ElasticsearchServiceGRPCClient != nil {
		convertReqDTO := &elasticsearchservicepb.GetInvoicesRequest{}
		convertReqDTO.Offset = reqDTO.Offset
		convertReqDTO.Limit = reqDTO.Limit
		convertReqDTO.SortBy = reqDTO.SortBy
		convertReqDTO.UserId = reqDTO.UserId
		convertReqDTO.TotalAmountGte = reqDTO.TotalAmountGTE
		convertReqDTO.TotalAmountLte = reqDTO.TotalAmountLTE
		convertReqDTO.Status = reqDTO.Status
		convertReqDTO.CreatedAtGte = reqDTO.CreatedAtGTE
		convertReqDTO.CreatedAtLte = reqDTO.CreatedAtLTE

		grpcRes, err := infrastructure.ElasticsearchServiceGRPCClient.GetInvoices(ctx, convertReqDTO)
		if err != nil {
			return nil, fmt.Errorf("get invoices from elasticsearch-service failed: %s", err.Error())
		}

		return model.FromListInvoiceProtoToListInvoiceView(grpcRes.Invoices), nil
	} else {
		return nil, fmt.Errorf("elasticsearch-service is not running")
	}
}
