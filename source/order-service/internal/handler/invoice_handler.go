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
	}, invoiceHandler.GetInvocies)

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

	// Get my invoices
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices",
		Summary:     "/my-invoices",
		Description: "Get my invoices.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.GetMyInvoices)

	// Get my invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices/id/{id}",
		Summary:     "/my-invoices/id/{id}",
		Description: "Get my invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.GetMyInvoiceById)

	// Create my invoice
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/my-invoices",
		Summary:     "/my-invoices",
		Description: "Create my invoice.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication},
	}, invoiceHandler.CreateMyInvoice)

	return invoiceHandler
}

func (invoiceHandler *InvoiceHandler) GetInvocies(ctx context.Context, reqDTO *dto.GetInvoicesRequest) (*dto.PaginationBodyResponseList[*model.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoices(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*model.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.BodyResponse[*model.InvoiceView], error) {
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[*model.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice by id successful"
	res.Body.Data = foundInvoice
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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update invoice by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

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

func (invoiceHandler *InvoiceHandler) GetMyInvoices(ctx context.Context, reqDTO *dto.GetMyInvoicesRequest) (*dto.PaginationBodyResponseList[*model.InvoiceView], error) {
	convertReqDTO := &dto.GetInvoicesRequest{}
	convertReqDTO.Offset = reqDTO.Offset
	convertReqDTO.Limit = reqDTO.Limit
	convertReqDTO.SortBy = reqDTO.SortBy
	convertReqDTO.UserId = ctx.Value("user_id").(string)
	convertReqDTO.TotalAmountGTE = reqDTO.TotalAmountGTE
	convertReqDTO.TotalAmountLTE = reqDTO.TotalAmountLTE
	convertReqDTO.Status = reqDTO.Status
	convertReqDTO.CreatedAtGTE = reqDTO.CreatedAtGTE
	convertReqDTO.CreatedAtLTE = reqDTO.CreatedAtLTE

	invoices, err := invoiceHandler.invoiceService.GetInvoices(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get my invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*model.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get my invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetMyInvoiceById(ctx context.Context, reqDTO *dto.GetMyInvoiceByIdRequest) (*dto.BodyResponse[*model.InvoiceView], error) {
	if reqDTO.Id == "{id}" {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get my invoice by id failed"
		res.Details = []string{"missing path parameters: id"}
		return nil, res
	}

	convertReqDTO := &dto.GetInvoiceByIdRequest{}
	convertReqDTO.Id = reqDTO.Id
	convertReqDTO.UserId = ctx.Value("user_id").(string)

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get my invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[*model.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get my invoice by id successful"
	res.Body.Data = foundInvoice
	return res, nil
}

func (invoiceHandler *InvoiceHandler) CreateMyInvoice(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateInvoiceRequest{}
	convertReqDTO.Body.UserId = ctx.Value("user_id").(string)
	convertReqDTO.Body.InvoiceDetails = []dto.InvoiceDetail{}

	if err := invoiceHandler.invoiceService.CreateInvoice(ctx, convertReqDTO); err != nil {
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
