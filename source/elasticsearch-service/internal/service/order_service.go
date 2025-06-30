package service

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
	"thanhldt060802/internal/grpc/client/orderservicepb"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
	"thanhldt060802/internal/schema"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type orderService struct {
}

type OrderService interface {
	SyncAllAvailableInvoices() error

	GetInvoices(ctx context.Context, reqDTO *elasticsearchservicepb.GetInvoicesRequest) ([]*elasticsearchservicepb.Invoice, error)
	syncCreatingInvoiceLoop()
	syncUpdatingInvoiceLoop()
	syncDeletingInvoiceLoop()
}

func NewOrderService(sync string) OrderService {
	orderService := &orderService{}

	go func() {
		if sync == "true" {
			for range infrastructure.OrderServiceGRPCClientConnectionEvent {
				close(infrastructure.OrderServiceGRPCClientConnectionEvent)
				break
			}

			if err := orderService.SyncAllAvailableInvoices(); err != nil {
				log.Printf("Sync all available users the first time failed: %s", err.Error())
			} else {
				log.Printf("Sync all available users the first time successful")
			}

			infrastructure.OrderServiceGRPCConnection.Close()
		}

		go orderService.syncCreatingInvoiceLoop()
		go orderService.syncUpdatingInvoiceLoop()
		go orderService.syncDeletingInvoiceLoop()
	}()

	return orderService
}

func (orderService *orderService) SyncAllAvailableInvoices() error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"invoices"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index exists on Elasticsearch
	if existsRes.StatusCode == 200 {
		// Delete index on Elasticsearch
		res, err := infrastructure.ElasticsearchClient.Indices.Delete([]string{"invoices"})
		if err != nil {
			log.Fatalf("Cannot delete index invoices: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("delete invoices index on elasticsearch failed: %s", res.String())
		}
	}

	// Create index on Elasticsearch using custom schema
	res, err := infrastructure.ElasticsearchClient.Indices.Create("invoices",
		infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(schema.Invoice))))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create invoices index on elasticsearch failed: %s", res.String())
	}

	grpcRes, err := infrastructure.OrderServiceGRPCClient.GetAllInvoices(context.Background(), &orderservicepb.GetAllInvoicesRequest{})
	if err != nil {
		return fmt.Errorf("get all invoices from order-service failed: %s", err.Error())
	}
	invoices := grpcRes.Invoices

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
		if err := indexer.Close(context.Background()); err != nil {
			closeBulkIndexer = fmt.Errorf("close bulk indexer failed: %s", err.Error())
		}
	}()

	// Add all available data on PostgreSQL to BulkIndexer
	for _, invoice := range invoices {
		// Convert data to JSON data
		invoiceJSON, err := json.Marshal(dto.FromInvoiceProtoToInvoiceView(invoice))
		if err != nil {
			return err
		}

		// Add data to BulkIndexer
		err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: invoice.Id,
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
		return fmt.Errorf("sync all available invoices to elasticsearch failed: index invoice to bulk")
	}
	if closeBulkIndexer != nil {
		return fmt.Errorf("sync all available invoices to elasticsearch failed: %s", closeBulkIndexer.Error())
	}

	return nil
}

func (orderService *orderService) syncCreatingInvoiceLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "order-service.created-invoice")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		var newInvoiceView dto.InvoiceView
		if err := json.Unmarshal([]byte(msg.Payload), &newInvoiceView); err != nil {
			log.Printf("Parse payload from event order-service.created-invoice failed: %s", err.Error())
			continue
		}

		func() {
			res, err := infrastructure.ElasticsearchClient.Index(
				"invoices",
				esutil.NewJSONReader(newInvoiceView),
				infrastructure.ElasticsearchClient.Index.WithDocumentID(newInvoiceView.Id),
				infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Insert invoice to Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync creating invoice failed: %s", res.String())
			} else {
				log.Printf("Sync creating invoice successful")
			}
		}()
	}
}

func (orderService *orderService) syncUpdatingInvoiceLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "order-service.updated-invoice")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		var updatedInvoiceView dto.InvoiceView
		if err := json.Unmarshal([]byte(msg.Payload), &updatedInvoiceView); err != nil {
			log.Printf("Parse payload from event order-service.updated-invoice failed: %s", err.Error())
			continue
		}

		func() {
			res, err := infrastructure.ElasticsearchClient.Index(
				"invoices",
				esutil.NewJSONReader(updatedInvoiceView),
				infrastructure.ElasticsearchClient.Index.WithDocumentID(updatedInvoiceView.Id),
				infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Update invoice on Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync updating invoice failed: %s", res.String())
			} else {
				log.Printf("Sync updating invoice successful")
			}
		}()
	}
}

func (orderService *orderService) syncDeletingInvoiceLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "order-service.deleted-invoice")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		invoiceIdStr := msg.Payload

		func() {
			res, err := infrastructure.ElasticsearchClient.Delete(
				"invoices",
				invoiceIdStr,
				infrastructure.ElasticsearchClient.Delete.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Delete invoice from Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync deleting invoice failed: %s", res.String())
			} else {
				log.Printf("Sync deleting invoice successful")
			}
		}()
	}
}

func (orderService *orderService) GetInvoices(ctx context.Context, reqDTO *elasticsearchservicepb.GetInvoicesRequest) ([]*elasticsearchservicepb.Invoice, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by user_id
	if reqDTO.UserId != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"term": map[string]interface{}{
				"user_id.keyword": reqDTO.UserId,
			},
		})
	}

	// If searching by total_amount in range or partial range
	totalAmountRange := map[string]interface{}{}
	if reqDTO.TotalAmountGte != "" {
		value, _ := strconv.ParseInt(reqDTO.TotalAmountGte, 10, 64)
		totalAmountRange["gte"] = value
	}
	if reqDTO.TotalAmountLte != "" {
		value, _ := strconv.ParseInt(reqDTO.TotalAmountLte, 10, 64)
		totalAmountRange["lte"] = value
	}
	if len(totalAmountRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"total_amount_range": totalAmountRange,
			},
		})
	}

	// If searching by status
	if reqDTO.Status != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"status": reqDTO.Status,
			},
		})
	}

	// If searching by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if reqDTO.CreatedAtGte != "" {
		createdAtRange["gte"] = reqDTO.CreatedAtGte
	}
	if reqDTO.CreatedAtLte != "" {
		createdAtRange["lte"] = reqDTO.CreatedAtLte
	}
	if len(createdAtRange) > 0 {
		createdAtRange["format"] = "strict_date_optional_time" // For format YYYY-MM-ddTHH:mm:ss
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": createdAtRange,
			},
		})
	}

	// If not searching -> get all
	if len(mustConditions) == 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Setup query
	query := map[string]interface{}{
		"from": reqDTO.Offset,
		"size": reqDTO.Limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// Apply sorting to query
	sortFields := utils.ParseSorter(reqDTO.SortBy)
	_sortFields := []map[string]interface{}{}
	for _, sortField := range sortFields {
		_sortFields = append(_sortFields, map[string]interface{}{
			schema.UserStandardizeSortFieldMap[sortField.Field]: strings.ToLower(sortField.Direction),
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
		return nil, fmt.Errorf("some thing wrong when querying invoices on elasticsearch")
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

	return dto.FromListInvoiceViewToListInvoiceProto(invoices), nil
}
