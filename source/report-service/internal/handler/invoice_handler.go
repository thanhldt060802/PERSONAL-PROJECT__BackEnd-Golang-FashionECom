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

	// Statistics revenue
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/statistics-revenue",
		Summary:     "/invoices/statistics-revenue",
		Description: "Statistics revenue.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{jwtAuthMiddleware.Authentication, jwtAuthMiddleware.RequireAdmin},
	}, invoiceHandler.StatisticsRevenue)

	return invoiceHandler
}

func (invoiceHandler *InvoiceHandler) StatisticsRevenue(ctx context.Context, reqDTO *dto.StatisticsRevenueRequest) (*dto.BodyResponse[dto.RevenueReport], error) {
	report, err := invoiceHandler.invoiceService.StatisticsRevenue(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Statistics revenue failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[dto.RevenueReport]{}
	res.Body.Code = "OK"
	res.Body.Message = "Statistics revenue successful"
	res.Body.Data = *report
	return res, nil
}
