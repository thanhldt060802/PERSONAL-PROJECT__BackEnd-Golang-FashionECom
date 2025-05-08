package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/repository"
)

type invoiceService struct {
	invoiceElasticsearchRepository repository.InvoiceElasticsearchRepository
}

type InvoiceService interface {
	StatisticsRevenue(ctx context.Context, reqDTO *dto.StatisticsRevenueRequest) (*dto.RevenueReport, error)
}

func NewInvoiceService(invoiceElasticsearchRepository repository.InvoiceElasticsearchRepository) InvoiceService {
	return &invoiceService{
		invoiceElasticsearchRepository: invoiceElasticsearchRepository,
	}
}

func (invoiceService *invoiceService) StatisticsRevenue(ctx context.Context, reqDTO *dto.StatisticsRevenueRequest) (*dto.RevenueReport, error) {
	report, err := invoiceService.invoiceElasticsearchRepository.StatisticsRevenue(ctx,
		reqDTO.Status,
		reqDTO.TimeInterval,
		reqDTO.CreatedAtGTE,
		reqDTO.CreatedAtLTE,
	)
	if err != nil {
		return nil, fmt.Errorf("query invoices from elasticsearch failed: %s", err.Error())
	}

	return report, nil
}
