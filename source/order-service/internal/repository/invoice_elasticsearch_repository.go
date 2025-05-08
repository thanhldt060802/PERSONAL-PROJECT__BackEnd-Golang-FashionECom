package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type invoiceElasticsearchRepository struct {
}

type InvoiceElasticsearchRepository interface {
	SyncAllAvailable(ctx context.Context, invoices []model.Invoice) error

	SyncCreating(ctx context.Context, newInvoice *model.Invoice) error
	SyncUpdating(ctx context.Context, updatedInvoice *model.Invoice) error
	SyncDeletingById(ctx context.Context, id int64) error

	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField,
		status string,
		totalAmountGTL string,
		totalAmountLTE string,
		createdAtGTE string,
		createdAtLTE string,
	) ([]dto.InvoiceView, error)
}

func NewInvoiceElasticsearchRepository() InvoiceElasticsearchRepository {
	return &invoiceElasticsearchRepository{}
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SyncAllAvailable(ctx context.Context, invoices []model.Invoice) error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"invoices"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index does not exists on Elasticsearch
	if existsRes.StatusCode == 404 {
		// Create index on Elasticsearch using custom schema
		res, err := infrastructure.ElasticsearchClient.Indices.Create("invoices",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.InvoiceSchemaMappingElasticsearch))))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("some thing wrong when creating invoices index")
		}

		hasFailure := false
		var closeBulkIndexer error

		// Create BulkIndexer for above index to index to Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "invoices",
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
		for _, invoice := range invoices {
			// Convert data to JSON data
			invoiceJSON, err := json.Marshal(invoice)
			if err != nil {
				return err
			}

			// Add data to BulkIndexer
			err = indexer.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(invoice.Id, 10),
				Body:       bytes.NewReader(invoiceJSON),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index invoice with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
					hasFailure = true
				},
			})
			if err != nil {
				return err
			}
		}

		if hasFailure {
			return fmt.Errorf("some thing wrong when syncing all available invoices")
		}
		if closeBulkIndexer != nil {
			return fmt.Errorf("some thing wrong when syncing all available invoices (%s)", closeBulkIndexer.Error())
		}

		return nil
	}

	return fmt.Errorf("invoices index already exists after first sync all")
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SyncCreating(ctx context.Context, newInvoice *model.Invoice) error {
	// Add data to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"invoices",
		esutil.NewJSONReader(newInvoice),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(newInvoice.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("some thing wrong when syncing creating invoice")
	}

	return nil
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SyncUpdating(ctx context.Context, updatedInvoice *model.Invoice) error {
	// Update data on Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"invoices",
		esutil.NewJSONReader(updatedInvoice),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(updatedInvoice.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("some thing wrong when syncing updating invoice")
	}

	return nil
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SyncDeletingById(ctx context.Context, id int64) error {
	// Delete data from Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Delete(
		"invoices",
		strconv.FormatInt(id, 10),
		infrastructure.ElasticsearchClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("some thing wrong when syncing deleting invoice")
	}

	return nil
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField,
	status string,
	totalAmountGTL string,
	totalAmountLTE string,
	createdAtGTE string,
	createdAtLTE string,
) ([]dto.InvoiceView, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by full_name
	if status != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"status": status,
			},
		})
	}

	// If filtering by total_amount in range or partial range
	totalAmountRange := map[string]interface{}{}
	if totalAmountGTL != "" {
		totalAmountRange["gte"] = totalAmountGTL
	}
	if totalAmountLTE != "" {
		totalAmountRange["lte"] = totalAmountLTE
	}
	if len(totalAmountRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"total_amount": totalAmountRange,
			},
		})
	}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if createdAtGTE != "" {
		createdAtRange["gte"] = createdAtGTE
	}
	if createdAtLTE != "" {
		createdAtRange["lte"] = createdAtLTE
	}
	if len(createdAtRange) > 0 {
		createdAtRange["format"] = "strict_date_optional_time" // For format YYYY-MM-ddTHH:mm:ss
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": createdAtRange,
			},
		})
	}

	// If not filtering -> get all
	if len(mustConditions) == 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Setup query
	query := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// Apply sorting to query
	_sortFields := []map[string]interface{}{}
	for _, sortField := range sortFields {
		_sortFields = append(_sortFields, map[string]interface{}{
			model.InvoiceSchemaMappingElasticsearchSortFieldMap[sortField.Field]: strings.ToLower(sortField.Direction),
		})
	}
	query["sort"] = _sortFields

	// Convert query to JSON query
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("invoices"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse Elasticsearch response
	if res.IsError() {
		return nil, fmt.Errorf("some thing wrong when querying invoices")
	}

	// Declare Elasticsearch response
	var elasticsearchResponse struct {
		Hits struct {
			Hits []struct {
				Source dto.InvoiceView `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	// Unmarshal Elasticsearch response body to Elasticsearch response
	elasticsearchResponseBody := json.NewDecoder(res.Body)
	if err := elasticsearchResponseBody.Decode(&elasticsearchResponse); err != nil {
		return nil, err
	}

	// Extract data from Elasticsearch response
	invoices := make([]dto.InvoiceView, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		invoices[i] = hit.Source
	}

	return invoices, nil
}
