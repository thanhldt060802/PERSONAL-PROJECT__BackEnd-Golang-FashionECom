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
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
	"thanhldt060802/internal/schema"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type catalogService struct {
}

type CatalogService interface {
	syncAllAvailableProducts() error

	GetProducts(ctx context.Context, reqDTO *elasticsearchservicepb.GetProductsRequest) ([]*elasticsearchservicepb.Product, error)
	syncCreatingProductLoop()
	syncUpdatingProductLoop()
	syncDeletingProductLoop()

	// StatisticsNumberOfProductsCreated(ctx context.Context, reqDTO *dto.StatisticsNumberOfProductsCreatedRequest) (*dto.NumberOfProductsCreatedReport, error)
}

func NewCatalogService() CatalogService {
	catalogService := &catalogService{}

	if err := catalogService.syncAllAvailableProducts(); err != nil {
		log.Printf("Sync all available products the first time failed: %s", err.Error())
	} else {
		log.Printf("Sync all available products the first time successful")
	}

	go catalogService.syncCreatingProductLoop()
	go catalogService.syncUpdatingProductLoop()
	go catalogService.syncDeletingProductLoop()

	return catalogService
}

func (catalogService *catalogService) syncAllAvailableProducts() error {
	// Check if index already exists on Elasticsearch
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"products"})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()

	// If index does not exists on Elasticsearch
	if existsRes.StatusCode == 404 {
		grpcRes, err := infrastructure.CatalogServiceGRPCClient.GetAllProducts(context.Background(), &catalogservicepb.GetAllProductsRequest{})
		if err != nil {
			return fmt.Errorf("get all product from catalog-service failed: %s", err.Error())
		}
		products := grpcRes.Products

		// Create index on Elasticsearch using custom schema
		res, err := infrastructure.ElasticsearchClient.Indices.Create("products",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(schema.Product))))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("create products index on elasticsearch failed: %s", res.String())
		}

		hasFailure := false
		var closeBulkIndexer error

		// Create BulkIndexer for above index to index to Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "products",
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
		for _, product := range products {
			// Convert data to JSON data
			productJSON, err := json.Marshal(dto.FromProductProtoToProductView(product))
			if err != nil {
				return err
			}

			// Add data to BulkIndexer
			err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(product.Id, 10),
				Body:       bytes.NewReader(productJSON),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index product with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
					hasFailure = true
				},
			})
			if err != nil {
				return err
			}
		}

		if hasFailure {
			return fmt.Errorf("sync all available products to elasticsearch failed: index product to bulk")
		}
		if closeBulkIndexer != nil {
			return fmt.Errorf("sync all available products to elasticsearch failed: %s", closeBulkIndexer.Error())
		}

		return nil
	}

	return fmt.Errorf("products index on elasticsearch already exists after first syncing all")
}

func (catalogService *catalogService) syncCreatingProductLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "catalog-service.created-product")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		var newProductView dto.ProductView
		if err := json.Unmarshal([]byte(msg.Payload), &newProductView); err != nil {
			log.Printf("Parse payload from event catalog-service.created-product failed: %s", err.Error())
			continue
		}

		func() {
			res, err := infrastructure.ElasticsearchClient.Index(
				"products",
				esutil.NewJSONReader(newProductView),
				infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(newProductView.Id, 10)),
				infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Insert product to Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync creating product failed: %s", res.String())
			} else {
				log.Printf("Sync creating product successful")
			}
		}()
	}
}

func (catalogService *catalogService) syncUpdatingProductLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "catalog-service.updated-product")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		var updatedProductView dto.ProductView
		if err := json.Unmarshal([]byte(msg.Payload), &updatedProductView); err != nil {
			log.Printf("Parse payload from event catalog-service.updated-product failed: %s", err.Error())
			continue
		}

		func() {
			res, err := infrastructure.ElasticsearchClient.Index(
				"products",
				esutil.NewJSONReader(updatedProductView),
				infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(updatedProductView.Id, 10)),
				infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Update product on Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync updating product failed: %s", res.String())
			} else {
				log.Printf("Sync updating product successful")
			}
		}()
	}
}

func (catalogService *catalogService) syncDeletingProductLoop() {
	subscribe := infrastructure.RedisClient.Subscribe(context.Background(), "catalog-service.deleted-product")
	defer subscribe.Close()

	ch := subscribe.Channel()

	for msg := range ch {
		productIdStr := msg.Payload

		func() {
			res, err := infrastructure.ElasticsearchClient.Delete(
				"products",
				productIdStr,
				infrastructure.ElasticsearchClient.Delete.WithRefresh("true"),
			)
			if err != nil {
				log.Printf("Delete product from Elasticsearch failed: %s", err.Error())
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Sync deleting product failed: %s", res.String())
			} else {
				log.Printf("Sync deleting product successful")
			}
		}()
	}
}

func (catalogService *catalogService) GetProducts(ctx context.Context, reqDTO *elasticsearchservicepb.GetProductsRequest) ([]*elasticsearchservicepb.Product, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by name
	if reqDTO.Name != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"name": reqDTO.Name,
			},
		})
	}

	// If filtering by description
	if reqDTO.Description != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"description": reqDTO.Description,
			},
		})
	}

	// If filtering by sex
	if reqDTO.Sex != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"sex": reqDTO.Sex,
			},
		})
	}

	// If filtering by price in range or partial range
	priceRange := map[string]interface{}{}
	if reqDTO.PriceGte != "" {
		value, _ := strconv.ParseInt(reqDTO.PriceGte, 10, 64)
		priceRange["gte"] = value
	}
	if reqDTO.PriceLte != "" {
		value, _ := strconv.ParseInt(reqDTO.PriceLte, 10, 64)
		priceRange["lte"] = value
	}
	if len(priceRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"price_range": priceRange,
			},
		})
	}

	// If filtering by discount_percentage in range or partial range
	discountPercentageRange := map[string]interface{}{}
	if reqDTO.DiscountPercentageGte != "" {
		value, _ := strconv.ParseInt(reqDTO.DiscountPercentageGte, 10, 64)
		discountPercentageRange["gte"] = value
	}
	if reqDTO.DiscountPercentageLte != "" {
		value, _ := strconv.ParseInt(reqDTO.DiscountPercentageLte, 10, 64)
		discountPercentageRange["lte"] = value
	}
	if len(discountPercentageRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"discount_percentage_range": discountPercentageRange,
			},
		})
	}

	// If filtering by stock in range or partial range
	stockRange := map[string]interface{}{}
	if reqDTO.StockGte != "" {
		value, _ := strconv.ParseInt(reqDTO.StockGte, 10, 64)
		stockRange["gte"] = value
	}
	if reqDTO.StockLte != "" {
		value, _ := strconv.ParseInt(reqDTO.StockLte, 10, 64)
		stockRange["lte"] = value
	}
	if len(stockRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"stock_range": stockRange,
			},
		})
	}

	// If filtering by category_name
	if reqDTO.CategoryName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"category_name": reqDTO.CategoryName,
			},
		})
	}

	// If filtering by brand_name
	if reqDTO.BrandName != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"brand_name": reqDTO.BrandName,
			},
		})
	}

	// If filtering by created_at in range or partial range
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

	// If not filtering -> get all
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
			schema.ProductStandardizeSortFieldMap[sortField.Field]: strings.ToLower(sortField.Direction),
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
		infrastructure.ElasticsearchClient.Search.WithIndex("products"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse Elasticsearch response
	if res.IsError() {
		return nil, fmt.Errorf("some thing wrong when querying products on elasticsearch")
	}

	// Declare Elasticsearch response
	var elasticsearchResponse struct {
		Hits struct {
			Hits []struct {
				Source dto.ProductView `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	// Unmarshal Elasticsearch response body to Elasticsearch response
	elasticsearchResponseBody := json.NewDecoder(res.Body)
	if err := elasticsearchResponseBody.Decode(&elasticsearchResponse); err != nil {
		return nil, err
	}

	// Extract data from Elasticsearch response
	products := make([]dto.ProductView, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		products[i] = hit.Source
	}

	return dto.FromListProductViewToListProductProto(products), nil
}
