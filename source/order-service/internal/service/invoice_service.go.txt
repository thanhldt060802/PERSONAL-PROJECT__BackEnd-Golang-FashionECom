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

type invoiceService struct {
	invoiceRepository       repository.InvoiceRepository
	invoiceDetailRepository repository.InvoiceDetailRepository
}

type InvoiceService interface {
	// Integrate with Elasticsearch
	GetAllInvoices(ctx context.Context) ([]dto.InvoiceView, error)

	// Main features
	GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdRequest) ([]dto.InvoiceView, error)
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error)
	CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error

	// Integrate with Elasticsearch
	// GetInvoices()
	// SyncAllAvailableInvoicesToElasticsearch(ctx context.Context) error
	// SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context) error
	// GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]dto.InvoiceView, error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, invoiceDetailRepository repository.InvoiceDetailRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository:       invoiceRepository,
		invoiceDetailRepository: invoiceDetailRepository,
	}
}

//
//
// Integrate with Elasticsearch
// ######################################################################################

func (invoiceService *invoiceService) GetAllInvoices(ctx context.Context) ([]dto.InvoiceView, error) {
	invoices, err := invoiceService.invoiceRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query invoices from postgresql failed: %s", err.Error())
	}

	invoiceViews := dto.ToListInvoiceView(invoices)
	for i, _ := range invoiceViews {
		foundInvoiceDetail, err := invoiceService.invoiceDetailRepository.GetAllByInvoiceId(ctx, invoiceViews[i].Id)
		if err != nil {
			return nil, fmt.Errorf("query invoice details from postgresql failed: %s", err.Error())
		}

		invoiceViews[i].Details = dto.ToListInvoiceDetailView(foundInvoiceDetail)
	}

	return invoiceViews, nil
}

//
//
// Main features
// ######################################################################################

func (invoiceService *invoiceService) GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdRequest) ([]dto.InvoiceView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceRepository.GetByUserId(ctx, reqDTO.UserId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return dto.ToListInvoiceView(invoices), nil
}

func (invoiceService *invoiceService) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error) {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	foundInvoiceDetails, _ := invoiceService.invoiceDetailRepository.GetAllByInvoiceId(ctx, reqDTO.Id)

	return dto.ToInvoiceView(foundInvoice, dto.ToListInvoiceDetailView(foundInvoiceDetails)), nil
}

func (invoiceService *invoiceService) CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error {
	newInvoice := &model.Invoice{
		UserId:      reqDTO.Body.UserId,
		TotalAmount: reqDTO.Body.TotalAmount,
		Status:      "CREATED",
	}

	newInvoiceDetails := []model.InvoiceDetail{}
	for _, invoiceDetail := range reqDTO.Body.Details {
		newInvoiceDetails = append(newInvoiceDetails, model.InvoiceDetail{
			InvoiceId:          newInvoice.Id,
			ProductId:          invoiceDetail.ProductId,
			Price:              invoiceDetail.Price,
			DiscountPercentage: invoiceDetail.DiscountPercentage,
			Quantity:           invoiceDetail.Quantity,
			TotalPrice:         invoiceDetail.TotalPrice,
		})
	}

	if err := invoiceService.invoiceRepository.Create(ctx, newInvoice, newInvoiceDetails); err != nil {
		return fmt.Errorf("insert invoice to postgresql failed: %s", err.Error())
	}

	// Missing->SyncCreatingToElasticsearch

	return nil
}

func (invoiceService *invoiceService) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	if reqDTO.Body.Status != nil {
		foundInvoice.Status = *reqDTO.Body.Status
	}
	foundInvoice.UpdatedAt = time.Now().UTC()

	if err := invoiceService.invoiceRepository.Update(ctx, foundInvoice); err != nil {
		return fmt.Errorf("update invoice on postgresql failed: %s", err.Error())
	}

	// Missing->SyncUpdatingToElasticsearch

	return nil
}

func (invoiceService *invoiceService) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error {
	if _, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of invoice is not valid")
	}

	if err := invoiceService.invoiceRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete invoice from postgresql failed: %s", err.Error())
	}

	// Missing->SyncDeletingToElasticsearch

	return nil
}

// Integrate with Elasticsearch

// func (invoiceService *invoiceService) SyncAllAvailableInvoicesToElasticsearch(ctx context.Context) error {
// 	invoices, err := invoiceService.invoiceRepository.GetAll(ctx)
// 	if err != nil {
// 		return fmt.Errorf("query invoices from postgresql failed: %s", err.Error())
// 	}

// 	if err := invoiceService.invoiceElasticsearchRepository.SyncAllAvailable(ctx, invoices); err != nil {
// 		return fmt.Errorf("sync all available invoices to elasticsearch failed: %s", err.Error())
// 	}

// 	return nil
// }

// func (invoiceService *invoiceService) SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context) error {
// 	invoiceDetails, err := invoiceService.invoiceDetailRepository.GetAll(ctx)
// 	if err != nil {
// 		return fmt.Errorf("query invoice details from postgresql failed: %s", err.Error())
// 	}

// 	if err := invoiceService.invoiceDetailElasticsearchRepository.SyncAllAvailable(ctx, invoiceDetails); err != nil {
// 		return fmt.Errorf("sync all available invoice details to elasticsearch failed: %s", err.Error())
// 	}

// 	return nil
// }

// func (invoiceService *invoiceService) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]dto.InvoiceView, error) {
// 	sortFields := utils.ParseSorter(reqDTO.SortBy)

// 	invoices, err := invoiceService.invoiceElasticsearchRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields,
// 		reqDTO.Status,
// 		reqDTO.TotalAmountGTE,
// 		reqDTO.TotalAmountLTE,
// 		reqDTO.CreatedAtGTE,
// 		reqDTO.CreatedAtLTE,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("query invoices from elasticsearch failed: %s", err.Error())
// 	}

// 	return invoices, nil
// }
