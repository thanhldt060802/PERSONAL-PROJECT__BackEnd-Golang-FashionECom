package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"time"

	"github.com/google/uuid"
)

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
	brandRepository    repository.BrandRepository
}

type ProductService interface {
	GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*model.ProductView, error)
	CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error
	UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductByIdRequest) error
	DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductByIdRequest) error

	// Elasticsearch integration (init data for elasticsearch-service)
	GetAllProducts(ctx context.Context) ([]*model.ProductView, error)

	// Order integration (extra features for order-service)
	UpdateProductStocksByListInvoiceDetail(ctx context.Context, reqDTO *dto.UpdateProductStocksByListInvoiceDetailRequest) error

	// Elasticsearch integration features
	GetProducts(ctx context.Context, reqDTO *dto.GetProductsRequest) ([]*model.ProductView, error)
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository, brandRepository repository.BrandRepository) ProductService {
	return &productService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		brandRepository:    brandRepository,
	}
}

func (productService *productService) GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*model.ProductView, error) {
	foundProduct, err := productService.productRepository.GetViewById(ctx, reqDTO.Id)
	if err != nil {
		return nil, fmt.Errorf("id of product is not valid: %s", err.Error())
	}

	return foundProduct, nil
}

func (productService *productService) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error {
	newProduct := model.Product{
		Id:                 uuid.New().String(),
		Name:               reqDTO.Body.Name,
		Description:        reqDTO.Body.Description,
		Sex:                reqDTO.Body.Sex,
		Price:              reqDTO.Body.Price,
		DiscountPercentage: reqDTO.Body.DiscountPercentage,
		Stock:              reqDTO.Body.Stock,
		ImageURL:           reqDTO.Body.ImageURL,
		CategoryId:         reqDTO.Body.CategoryId,
		BrandId:            reqDTO.Body.BrandId,
	}
	if err := productService.productRepository.Create(ctx, &newProduct); err != nil {
		return fmt.Errorf("insert product to postgresql failed: %s", err.Error())
	}

	newProductView, _ := productService.productRepository.GetViewById(ctx, newProduct.Id)
	payload, _ := json.Marshal(newProductView)
	if err := infrastructure.RedisClient.Publish(ctx, "catalog-service.created-product", payload).Err(); err != nil {
		return fmt.Errorf("pulish event catalog-service.created-product failed: %s", err.Error())
	}

	return nil
}

func (productService *productService) UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductByIdRequest) error {
	foundProduct, err := productService.productRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of product is not valid: %s", err.Error())
	}

	if reqDTO.Body.Name != nil {
		foundProduct.Name = *reqDTO.Body.Name
	}
	if reqDTO.Body.Description != nil {
		foundProduct.Description = *reqDTO.Body.Description
	}
	if reqDTO.Body.Sex != nil {
		foundProduct.Sex = *reqDTO.Body.Sex
	}
	if reqDTO.Body.Price != nil {
		foundProduct.Price = *reqDTO.Body.Price
	}
	if reqDTO.Body.DiscountPercentage != nil {
		foundProduct.DiscountPercentage = *reqDTO.Body.DiscountPercentage
	}
	if reqDTO.Body.Stock != nil {
		foundProduct.Stock = *reqDTO.Body.Stock
	}
	if reqDTO.Body.ImageURL != nil {
		foundProduct.ImageURL = *reqDTO.Body.ImageURL
	}
	if reqDTO.Body.CategoryId != nil {
		if _, err = productService.categoryRepository.GetById(ctx, *reqDTO.Body.CategoryId); err != nil {
			return fmt.Errorf("id of category not found")
		}
		foundProduct.CategoryId = *reqDTO.Body.CategoryId
	}
	if reqDTO.Body.BrandId != nil {
		if _, err = productService.brandRepository.GetById(ctx, *reqDTO.Body.BrandId); err != nil {
			return fmt.Errorf("id of brand not found")
		}
		foundProduct.BrandId = *reqDTO.Body.BrandId
	}
	timeUpdate := time.Now().UTC()
	foundProduct.UpdatedAt = &timeUpdate

	if err := productService.productRepository.Update(ctx, foundProduct); err != nil {
		return fmt.Errorf("update product on postgresql failed: %s", err.Error())
	}

	updatedProductView, _ := productService.productRepository.GetViewById(ctx, foundProduct.Id)
	payload, _ := json.Marshal(updatedProductView)
	if err := infrastructure.RedisClient.Publish(ctx, "catalog-service.updated-product", payload).Err(); err != nil {
		return fmt.Errorf("pulish event catalog-service.updated-product failed: %s", err.Error())
	}

	return nil
}

func (productService *productService) DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductByIdRequest) error {
	if _, err := productService.productRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of product is not valid")
	}

	if err := productService.productRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("delete product from postgresql failed: %s", err.Error())
	}

	if err := infrastructure.RedisClient.Publish(ctx, "catalog-service.deleted-product", reqDTO.Id).Err(); err != nil {
		return fmt.Errorf("pulish event product-service.deleted-product failed: %s", err.Error())
	}

	return nil
}

func (productService *productService) GetAllProducts(ctx context.Context) ([]*model.ProductView, error) {
	products, err := productService.productRepository.GetAllViews(ctx)
	if err != nil {
		return nil, fmt.Errorf("query products from postgresql failed: %s", err.Error())
	}

	return products, nil
}

func (productService *productService) UpdateProductStocksByListInvoiceDetail(ctx context.Context, reqDTO *dto.UpdateProductStocksByListInvoiceDetailRequest) error {
	ids := make([]string, len(reqDTO.InvoiceDetails))
	quantityMap := map[string]int32{}
	for i, invoiceDetail := range reqDTO.InvoiceDetails {
		ids[i] = invoiceDetail.ProductId
		quantityMap[invoiceDetail.ProductId] = invoiceDetail.Quantity
	}

	foundProducts, err := productService.productRepository.GetByListId(ctx, ids)
	if err != nil {
		return fmt.Errorf("query products from postgresql failed: %s", err.Error())
	}

	for i := range foundProducts {
		if foundProducts[i].Stock < quantityMap[foundProducts[i].Id] {
			return fmt.Errorf("not enough stock for product id: %s", foundProducts[i].Id)
		}
		foundProducts[i].Stock = foundProducts[i].Stock - quantityMap[foundProducts[i].Id]
		timeUpdate := time.Now().UTC()
		foundProducts[i].UpdatedAt = &timeUpdate
	}

	if err := productService.productRepository.UpdateStocks(ctx, foundProducts); err != nil {
		return fmt.Errorf("update stock of products from postgresql failed: %s", err.Error())
	}

	for _, foundProduct := range foundProducts {
		updatedProductView, _ := productService.productRepository.GetViewById(ctx, foundProduct.Id)
		payload, _ := json.Marshal(updatedProductView)
		if err := infrastructure.RedisClient.Publish(ctx, "catalog-service.updated-product", payload).Err(); err != nil {
			return fmt.Errorf("pulish event catalog-service.updated-product failed: %s", err.Error())
		}
	}

	return nil
}

func (productService *productService) GetProducts(ctx context.Context, reqDTO *dto.GetProductsRequest) ([]*model.ProductView, error) {
	if infrastructure.ElasticsearchServiceGRPCClient != nil {
		convertReqDTO := &elasticsearchservicepb.GetProductsRequest{}
		convertReqDTO.Offset = reqDTO.Offset
		convertReqDTO.Limit = reqDTO.Limit
		convertReqDTO.SortBy = reqDTO.SortBy
		convertReqDTO.CategoryId = reqDTO.CategoryId
		convertReqDTO.BrandId = reqDTO.BrandId
		convertReqDTO.Name = reqDTO.Name
		convertReqDTO.Description = reqDTO.Description
		convertReqDTO.Sex = reqDTO.Sex
		convertReqDTO.PriceGte = reqDTO.PriceGTE
		convertReqDTO.PriceLte = reqDTO.PriceLTE
		convertReqDTO.DiscountPercentageGte = reqDTO.DiscountPercentageGTE
		convertReqDTO.DiscountPercentageLte = reqDTO.DiscountPercentageLTE
		convertReqDTO.StockGte = reqDTO.StockGTE
		convertReqDTO.StockLte = reqDTO.StockLTE
		convertReqDTO.CategoryName = reqDTO.CategoryName
		convertReqDTO.BrandName = reqDTO.BrandName
		convertReqDTO.CreatedAtGte = reqDTO.CreatedAtGTE
		convertReqDTO.CreatedAtLte = reqDTO.CreatedAtLTE

		grpcRes, err := infrastructure.ElasticsearchServiceGRPCClient.GetProducts(ctx, convertReqDTO)
		if err != nil {
			return nil, fmt.Errorf("get products from elasticsearch-service failed: %s", err.Error())
		}

		return model.FromListProductProtoToListProductView(grpcRes.Products), nil
	} else {
		return nil, fmt.Errorf("elasticsearch-service is not running")
	}
}
