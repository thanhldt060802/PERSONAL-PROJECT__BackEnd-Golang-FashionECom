package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type invoiceDetailElasticsearchRepository struct {
}

type InvoiceDetailElasticsearchRepository interface {
	SyncAllAvailable(ctx context.Context, invoiceDetails []model.InvoiceDetail) error

	SyncCreating(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error
	SyncDeletingById(ctx context.Context, invoiceDetails []model.InvoiceDetail) error
}

func NewInvoiceDetailElasticsearchRepository() InvoiceDetailElasticsearchRepository {
	return &invoiceDetailElasticsearchRepository{}
}

func (invoiceDetailElasticsearchRepository *invoiceDetailElasticsearchRepository) SyncAllAvailable(ctx context.Context, invoiceDetails []model.InvoiceDetail) error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"invoice_details"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index does not exists on Elasticsearch
	if existsRes.StatusCode == 404 {
		// Create index on Elasticsearch using custom schema
		res, err := infrastructure.ElasticsearchClient.Indices.Create("invoice_details",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.InvoiceDetailSchemaMappingElasticsearch))))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("some thing wrong when creating invoice_details index")
		}

		hasFailure := false
		var closeBulkIndexer error

		// Create BulkIndexer for above index to index to Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "invoice_details",
		})
		if err != nil {
			return err
		}
		defer func() {
			if err := indexer.Close(ctx); err != nil {
				closeBulkIndexer = fmt.Errorf("close bulk indexer failed: %s", err.Error())
			}
		}()

		// Add all available data on PostgreSQL to BulkIndexer
		for _, invoiceDetail := range invoiceDetails {
			// Convert data to JSON data
			invoiceDetailJSON, err := json.Marshal(invoiceDetail)
			if err != nil {
				return err
			}

			// Add data to BulkIndexer
			err = indexer.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(invoiceDetail.Id, 10),
				Body:       bytes.NewReader(invoiceDetailJSON),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index invoice details with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
					hasFailure = true
				},
			})
			if err != nil {
				return err
			}
		}

		if hasFailure {
			return fmt.Errorf("some thing wrong when syncing all available invoice details")
		}
		if closeBulkIndexer != nil {
			return fmt.Errorf("some thing wrong when syncing all available invoice details (%s)", closeBulkIndexer.Error())
		}

		return nil
	}

	return fmt.Errorf("invoice_details index already exists after first sync all")
}

func (invoiceDetailElasticsearchRepository *invoiceDetailElasticsearchRepository) SyncCreating(ctx context.Context, newInvoiceDetails []model.InvoiceDetail) error {
	hasFailure := false
	var closeBulkIndexer error

	// Add many data to Elasticsearch by BulkIndexer
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: infrastructure.ElasticsearchClient,
		Index:  "invoice_details",
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := indexer.Close(ctx); err != nil {
			closeBulkIndexer = fmt.Errorf("close bulk indexer failed: %s", err.Error())
		}
	}()

	for _, newInvoiceDetail := range newInvoiceDetails {
		newInvoiceDetailJSON, err := json.Marshal(newInvoiceDetail)
		if err != nil {
			return err
		}

		err = indexer.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.FormatInt(newInvoiceDetail.Id, 10),
			Body:       bytes.NewReader(newInvoiceDetailJSON),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("Bulk index failed: %s", err.Error())
				} else {
					log.Printf("Index invoice detail with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
				}
				hasFailure = true
			},
		})
		if err != nil {
			return err
		}
	}

	if hasFailure {
		return fmt.Errorf("some thing wrong when syncing creating invoice details")
	}
	if closeBulkIndexer != nil {
		return fmt.Errorf("some thing wrong when syncing creating invoice details (%s)", closeBulkIndexer.Error())
	}

	return nil
}

func (invoiceDetailElasticsearchRepository *invoiceDetailElasticsearchRepository) SyncDeletingById(ctx context.Context, invoiceDetails []model.InvoiceDetail) error {
	hasFailure := false
	var closeBulkIndexer error

	// Delete many data from Elasticsearch by BulkIndexer
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: infrastructure.ElasticsearchClient,
		Index:  "invoice_details",
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := indexer.Close(ctx); err != nil {
			closeBulkIndexer = fmt.Errorf("close bulk indexer failed: %s", err.Error())
		}
	}()

	for _, invoiceDetail := range invoiceDetails {
		err = indexer.Add(ctx, esutil.BulkIndexerItem{
			Action:     "delete",
			DocumentID: strconv.FormatInt(invoiceDetail.Id, 10),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("Bulk index failed: %s", err.Error())
				} else {
					log.Printf("Unindex invoice detail with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
				}
				hasFailure = true
			},
		})
		if err != nil {
			return err
		}
	}

	if hasFailure {
		return fmt.Errorf("some thing wrong when syncing deleting invoice details")
	}
	if closeBulkIndexer != nil {
		return fmt.Errorf("some thing wrong when syncing deleting invoice details (%s)", closeBulkIndexer.Error())
	}

	return nil
}
