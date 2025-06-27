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
	cartItemService   service.CartItemService
	jwtAuthMiddleware *middleware.JWTAuthMiddleware
}

func NewInvoiceHandler(api huma.API, invoiceService service.InvoiceService, cartItemService service.CartItemService, jwtAuthMiddleware *middleware.JWTAuthMiddleware) *InvoiceHandler {
	invoiceHandler := &InvoiceHandler{
		invoiceService:    invoiceService,
		cartItemService:   cartItemService,
		jwtAuthMiddleware: jwtAuthMiddleware,
	}

	return invoiceHandler
}

func (invoiceHandler *InvoiceHandler) GetInvocies(ctx context.Context, reqDTO *dto.GetInvoicesRequest) (*dto.PaginationBodyResponseList[*dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoices(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.PaginationBodyResponseList[*dto.InvoiceView]{}
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

func (invoiceHandler *InvoiceHandler) GetMyInvoices(ctx context.Context, reqDTO *dto.GetMyInvoicesRequest) (*dto.PaginationBodyResponseList[*dto.InvoiceView], error) {
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

	res := &dto.PaginationBodyResponseList[*dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get my invoices successful"
	res.Body.Data = invoices
	res.Body.Total = len(invoices)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) CreateMyInvoice(ctx context.Context, _ *struct{}) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateInvoiceRequest{}
	convertReqDTO.Body.UserId = ctx.Value("user_id").(string)

	cartItems, err := invoiceHandler.cartItemService.GetCartItems(ctx, convertReqDT)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get my cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	convertReqDTO.Body.ProductId = reqDTO.Body.ProductId

	if err := cartItemHandler.cartItemService.CreateCartItem(ctx, convertReqDTO); err != nil {
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
