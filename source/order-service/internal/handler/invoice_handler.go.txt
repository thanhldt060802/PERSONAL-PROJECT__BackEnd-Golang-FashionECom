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

	//
	//
	// Integrate with Elasticsearch
	// ######################################################################################

	// Get all users (integrate with Elasticsearch)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/all",
		Summary:     "/invoices/all",
		Description: "Get all users (integrate with Elasticsearch).",
		Tags:        []string{"For Sycing Data To Elasticsearch"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetAllInvoices)

	//
	//
	// Main features
	// ######################################################################################

	// Get invoices by user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/user-id/{user_id}",
		Summary:     "/invoices/user-id/{user_id}",
		Description: "Get invoices by user id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoicesByUserId)

	// Get invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Get invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoiceById)

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

	// Get account invoices
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices",
		Summary:     "/my-invoices",
		Description: "Get account invoices.",
		Tags:        []string{"Account Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.GetAccountInvoices)

	// Get invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices/id/{id}",
		Summary:     "/my-invoices/id/{id}",
		Description: "Get account invoice by id.",
		Tags:        []string{"Account Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.GetAccountInvoiceById)

	// Delete account invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/my-invoices/id/{id}",
		Summary:     "/my-invoices/id/{id}",
		Description: "Delete account invoice by id.",
		Tags:        []string{"Account Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.DeleteAccountInvoiceById)

	//
	//
	// Elasticsearch integration features
	// ######################################################################################

	// Get invoices

	return invoiceHandler
}

//
//
// Integrate with Elasticsearch
// ######################################################################################

func (invoiceHandler *InvoiceHandler) GetAllInvoices(ctx context.Context, _ *struct{}) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetAllInvoices(ctx)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get all invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get all invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

//
//
// Main features
// ######################################################################################

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

func (invoiceHandler *InvoiceHandler) GetAccountInvoices(ctx context.Context, reqDTO *dto.GetAccountInvoices) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	convertReqDTO := &dto.GetInvoicesByUserIdRequest{}
	convertReqDTO.UserId = ctx.Value("user_id").(int64)
	convertReqDTO.Offset = reqDTO.Offset
	convertReqDTO.Limit = reqDTO.Limit
	convertReqDTO.SortBy = reqDTO.SortBy

	invoices, err := invoiceHandler.invoiceService.GetInvoicesByUserId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get account invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get account invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetAccountInvoiceById(ctx context.Context, reqDTO *dto.GetAccountInvoiceByIdRequest) (*dto.BodyResponse[dto.InvoiceView], error) {
	convertReqDTO := &dto.GetInvoiceByIdRequest{}
	convertReqDTO.Id = reqDTO.Id

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get account invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	if foundInvoice.UserId != ctx.Value("user_id").(int64) {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusForbidden
		res.Code = "ERR_FORBIDDEN"
		res.Message = "Get account invoice by id failed"
		res.Details = []string{"id of invoice is not owned"}
		return nil, res
	}

	res := &dto.BodyResponse[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get account invoice by id successful"
	res.Body.Data = *foundInvoice
	return res, nil
}

func (invoiceHandler *InvoiceHandler) DeleteAccountInvoiceById(ctx context.Context, reqDTO *dto.DeleteAccountInvoiceByIdRequest) (*dto.SuccessResponse, error) {
	convertReqDTO1 := &dto.GetInvoiceByIdRequest{}
	convertReqDTO1.Id = reqDTO.Id

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, convertReqDTO1)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get account invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	if foundInvoice.UserId != ctx.Value("user_id").(int64) {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusForbidden
		res.Code = "ERR_FORBIDDEN"
		res.Message = "Delete account invoice by id failed"
		res.Details = []string{"id of invoice is not owned"}
		return nil, res
	}

	convertReqDTO2 := &dto.DeleteInvoiceByIdRequest{}
	convertReqDTO2.Id = reqDTO.Id

	if err := invoiceHandler.invoiceService.DeleteInvoiceById(ctx, convertReqDTO2); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete account invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete account invoice by id successful"
	return res, nil
}

// func (invoiceHandler *InvoiceHandler) SyncAllAvailableInvoicesToElasticsearch(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
// 	if err := invoiceHandler.invoiceService.SyncAllAvailableInvoicesToElasticsearch(ctx); err != nil {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusInternalServerError
// 		res.Code = "ERR_INTERNAL_SERVER"
// 		res.Message = "Sync all available invoices to Elasticsearch failed"
// 		res.Details = []string{err.Error()}
// 		return nil, res
// 	}

// 	res := &dto.SuccessResponse{}
// 	res.Body.Code = "OK"
// 	res.Body.Message = "Sync all available invoices to Elasticsearch successful"
// 	return res, nil
// }

// func (invoiceHandler *InvoiceHandler) SyncAllAvailableInvoiceDetailsToElasticsearch(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
// 	if err := invoiceHandler.invoiceService.SyncAllAvailableInvoiceDetailsToElasticsearch(ctx); err != nil {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusInternalServerError
// 		res.Code = "ERR_INTERNAL_SERVER"
// 		res.Message = "Sync all available invoice details to Elasticsearch failed"
// 		res.Details = []string{err.Error()}
// 		return nil, res
// 	}

// 	res := &dto.SuccessResponse{}
// 	res.Body.Code = "OK"
// 	res.Body.Message = "Sync all available invoice details to Elasticsearch successful"
// 	return res, nil
// }

// func (invoiceHandler *InvoiceHandler) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
// 	invoices, err := invoiceHandler.invoiceService.GetInvoicesWithElasticsearch(ctx, reqDTO)
// 	if err != nil {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusInternalServerError
// 		res.Code = "ERR_INTERNAL_SERVER"
// 		res.Message = "Get invoices with Elasticsearch failed"
// 		res.Details = []string{err.Error()}
// 		return nil, res
// 	}

// 	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
// 	res.Body.Code = "OK"
// 	res.Body.Message = "Get invoices with Elasticsearch successful"
// 	res.Body.Data = invoices
// 	res.Body.Total = len(invoices)
// 	return res, nil
// }
