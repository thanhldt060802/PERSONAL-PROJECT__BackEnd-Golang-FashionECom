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
	invoiceRepository                    repository.InvoiceRepository
	invoiceElasticsearchRepository       repository.InvoiceElasticsearchRepository
	invoiceDetailRepository              repository.InvoiceDetailRepository
	invoiceDetailElasticsearchRepository repository.InvoiceDetailElasticsearchRepository
}

type InvoiceService interface {
	GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) ([]dto.InvoiceView, error)
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error)
	GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdRequest) ([]dto.InvoiceView, error)
	CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error

	// Integrate with Elasticsearch

	SyncAllAvailableInvoicesToElasticsearch(ctx context.Context) error
	SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context) error
	GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]dto.InvoiceView, error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, invoiceElasticsearchRepository repository.InvoiceElasticsearchRepository, invoiceDetailRepository repository.InvoiceDetailRepository, invoiceDetailElasticsearchRepository repository.InvoiceDetailElasticsearchRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository:                    invoiceRepository,
		invoiceElasticsearchRepository:       invoiceElasticsearchRepository,
		invoiceDetailRepository:              invoiceDetailRepository,
		invoiceDetailElasticsearchRepository: invoiceDetailElasticsearchRepository,
	}
}

func (invoiceService *invoiceService) GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) ([]dto.InvoiceView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("query invoices from postgresql failed: %s", err.Error())
	}

	return dto.ToListInvoiceView(invoices), nil
}

func (invoiceService *invoiceService) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error) {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	foundInvoiceDetails, err := invoiceService.invoiceDetailRepository.GetAllByInvoiceId(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("invoice has no invoice details: %s", err.Error())
	}

	return dto.ToInvoiceView(foundInvoice, dto.ToListInvoiceDetailView(foundInvoiceDetails)), nil
}

func (invoiceService *invoiceService) GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdRequest) ([]dto.InvoiceView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceRepository.GetByUserId(ctx, reqDTO.UserId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return dto.ToListInvoiceView(invoices), nil
}

func (invoiceService *invoiceService) CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error {
	newInvoice := &model.Invoice{
		UserId:      ctx.Value("user_id").(int64),
		TotalAmount: reqDTO.Body.TotalAmount,
		Status:      "CREATED",
	}
	if err := invoiceService.invoiceRepository.Create(ctx, newInvoice); err != nil {
		return fmt.Errorf("insert invoice to postgresql failed: %s", err.Error())
	}

	if err := invoiceService.invoiceElasticsearchRepository.SyncCreating(ctx, newInvoice); err != nil {
		return fmt.Errorf("sync creating invoice to elasticsearch failed: %s", err.Error())
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
	if err := invoiceService.invoiceDetailRepository.Create(ctx, newInvoiceDetails); err != nil {
		return fmt.Errorf("insert invoice details to postgresql failed: %s", err.Error())
	}

	if err := invoiceService.invoiceDetailElasticsearchRepository.SyncCreating(ctx, newInvoiceDetails); err != nil {
		return fmt.Errorf("sync creating invoice details to elasticsearch failed: %s", err.Error())
	}

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

	if err := invoiceService.invoiceElasticsearchRepository.SyncUpdating(ctx, foundInvoice); err != nil {
		return fmt.Errorf("sync updating invoice on elasticsearch failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error {
	if _, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of invoice is not valid")
	}

	foundInvoiceDetails, _ := invoiceService.invoiceDetailRepository.GetAllByInvoiceId(ctx, reqDTO.Id)

	if err := invoiceService.invoiceDetailRepository.DeleteByInvoiceId(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete invoice details from postgresql failed: %s", err.Error())
	}
	if err := invoiceService.invoiceRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete invoice from postgresql failed: %s", err.Error())
	}

	if err := invoiceService.invoiceDetailElasticsearchRepository.SyncDeletingById(ctx, foundInvoiceDetails); err != nil {
		return fmt.Errorf("sync deleting invoice details from elasticsearch failed: %s", err.Error())
	}
	if err := invoiceService.invoiceElasticsearchRepository.SyncDeletingById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("sync deleting invoice from elasticsearch failed: %s", err.Error())
	}

	return nil
}

// Integrate with Elasticsearch

func (invoiceService *invoiceService) SyncAllAvailableInvoicesToElasticsearch(ctx context.Context) error {
	invoices, err := invoiceService.invoiceRepository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("query invoices from postgresql failed: %s", err.Error())
	}

	if err := invoiceService.invoiceElasticsearchRepository.SyncAllAvailable(ctx, invoices); err != nil {
		return fmt.Errorf("sync all available invoices to elasticsearch failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context) error {
	invoiceDetails, err := invoiceService.invoiceDetailRepository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("query invoice details from postgresql failed: %s", err.Error())
	}

	if err := invoiceService.invoiceDetailElasticsearchRepository.SyncAllAvailable(ctx, invoiceDetails); err != nil {
		return fmt.Errorf("sync all available invoice details to elasticsearch failed: %s", err.Error())
	}

	return nil
}

func (invoiceService *invoiceService) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]dto.InvoiceView, error) {
	sortFields := utils.ParseSorter(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceElasticsearchRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields,
		reqDTO.Status,
		reqDTO.TotalAmountGTE,
		reqDTO.TotalAmountLTE,
		reqDTO.CreatedAtGTE,
		reqDTO.CreatedAtLTE,
	)
	if err != nil {
		return nil, fmt.Errorf("query invoices from elasticsearch failed: %s", err.Error())
	}

	return invoices, nil
}
