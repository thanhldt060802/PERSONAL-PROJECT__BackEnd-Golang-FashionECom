package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type InvoiceHandler struct {
	invoiceService    service.InvoiceService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewInvoiceHandler(api huma.API, invoiceService service.InvoiceService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *InvoiceHandler {
	invoiceHandler := &InvoiceHandler{
		invoiceService:    invoiceService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	// Get invoices
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices",
		Summary:     "/invoices",
		Description: "Get invoices.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoices)

	// Get invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Get invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoiceById)

	// Get invoices by user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/user-id/{user_id}",
		Summary:     "/invoices/user-id/{user_id}",
		Description: "Get invoices by user id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoicesByUserId)

	// Create invoice
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/invoices",
		Summary:     "/invoices",
		Description: "Create invoice.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.CreateInvoice)

	// Update invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Update invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.UpdateInvoiceById)

	// Delete invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Delete invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.DeleteInvoiceById)

	// Integrate with Elasticsearch

	// Sync all available invoices to Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/sync-to-elasticsearch",
		Summary:     "/invoices/sync-to-elasticsearch",
		Description: "Sync all available invoices to Elasticsearch.",
		Tags:        []string{"Sync Data"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.SyncAllAvailableInvoicesToElasticsearch)

	// Sync all available invoice details to Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoice-details/sync-to-elasticsearch",
		Summary:     "/invoice-details/sync-to-elasticsearch",
		Description: "Sync all available invoice details to Elasticsearch.",
		Tags:        []string{"Sync Data"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.SyncAllAvailableInvoiceDetailsToElasticsearch)

	// Get invoices with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/elasticsearch",
		Summary:     "/invoices/elasticsearch",
		Description: "Get invoice with Elasticsearch.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoicesWithElasticsearch)

	return invoiceHandler
}

func (invoiceHandler *InvoiceHandler) GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoices(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.BodyResponse[dto.InvoiceView], error) {
	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice by id successful"
	res.Body.Data = *foundInvoice
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoicesByUserId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices by user id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices by user id successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) CreateInvoice(ctx context.Context, reqDTO *dto.CreateInvoiceRequest) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.CreateInvoice(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create invoice failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create invoice successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceByIdRequest) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.UpdateInvoiceById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update invoice by id successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceByIdRequest) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.DeleteInvoiceById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete invoice by id successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) SyncAllAvailableInvoicesToElasticsearch(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.SyncAllAvailableInvoicesToElasticsearch(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all available invoices to Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all available invoices to Elasticsearch successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.SyncAllAvailableInvoiceDetailsToElasticsearch(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all available invoice details to Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all available invoice details to Elasticsearch successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoicesWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices with Elasticsearch successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}
