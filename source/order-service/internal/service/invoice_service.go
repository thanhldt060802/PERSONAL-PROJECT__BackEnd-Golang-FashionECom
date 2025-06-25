package service

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"time"

	"github.com/google/uuid"
)

type invoiceService struct {
	invoiceRepository repository.InvoiceRepository
}

type InvoiceService interface {
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error)
	CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllInvoices(ctx context.Context) ([]dto.InvoiceView, error)

	// Elasticsearch integration features
	GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) ([]dto.InvoiceView, error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, invoiceDetailRepository repository.InvoiceDetailRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository: invoiceRepository,
	}
}

func (invoiceService *invoiceService) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.InvoiceView, error) {
	foundInvoice, foundInvoiceDetails, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	ids := make([]string, len(foundInvoiceDetails))
	for i, invoiceDetail := range foundInvoiceDetails {
		ids[i] = invoiceDetail.ProductId
	}

	convertReqDTO := &catalogservicepb.GetProductsByListIdRequest{}
	convertReqDTO.Ids = ids
	grpcRes, err := infrastructure.CatalogServiceGRPCClient.GetProductsByListId(ctx, convertReqDTO)
	if err != nil {
		return nil, fmt.Errorf("get products from catalog-service failed: %s", err.Error())
	}

	productProtos := grpcRes.Products
	foundInvoiceDetailExtraInfos := make([]dto.InvoiceDetailExtraInfo, len(productProtos))
	for i, productProto := range productProtos {
		foundInvoiceDetailExtraInfos[i].Name = productProto.Name
		foundInvoiceDetailExtraInfos[i].Sex = productProto.Sex
		foundInvoiceDetailExtraInfos[i].ImageURL = productProto.ImageUrl

		foundInvoiceDetailExtraInfos[i].CategoryId = productProto.CategoryId
		foundInvoiceDetailExtraInfos[i].CategoryName = productProto.CategoryName
		foundInvoiceDetailExtraInfos[i].BrandId = productProto.BrandId
		foundInvoiceDetailExtraInfos[i].BrandName = productProto.BrandName
	}

	return dto.ToInvoiceView(foundInvoice, dto.ToListInvoiceDetailView(foundInvoiceDetails, foundInvoiceDetailExtraInfos)), nil
}

func (invoiceService *invoiceService) CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) error {
	newInvoice := &model.Invoice{
		Id:          uuid.New().String(),
		UserId:      reqDTO.Body.UserId,
		TotalAmount: 0,
		Status:      "CREATED",
	}

	newInvoiceDetails := []model.InvoiceDetail{}
	for _, invoiceDetail := range reqDTO.Body.Details {
		newInvoiceDetails = append(newInvoiceDetails, model.InvoiceDetail{
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

	// Missing->SyncCreatingToElasticsearch

	return nil
}

func (invoiceService *invoiceService) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) error {
	foundInvoice, _, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
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

	// Missing->SyncUpdatingToElasticsearch

	return nil
}

func (invoiceService *invoiceService) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) error {
	foundInvoice, _, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of invoice is not valid: %s", err.Error())
	}

	if reqDTO.UserId != "" && reqDTO.UserId != foundInvoice.UserId {
		return fmt.Errorf("id of invoice is not valid: no permission")
	}

	if err := invoiceService.invoiceRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete invoice from postgresql failed: %s", err.Error())
	}

	// Missing->SyncDeletingToElasticsearch

	return nil
}

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
