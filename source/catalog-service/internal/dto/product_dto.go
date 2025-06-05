package dto

import (
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/grpc/service/catalogservicepb"
	"thanhldt060802/internal/model"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductView struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Sex                string    `json:"sex"`
	Price              int64     `json:"price"`
	DiscountPercentage int32     `json:"discount_percentage"`
	Stock              int32     `json:"stock"`
	ImageURL           string    `json:"image_url"`
	CategoryId         int64     `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	BrandId            int64     `json:"brand_id"`
	BrandName          string    `json:"brand_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func ToProductView(product *model.Product, category *model.Category, brand *model.Brand) *ProductView {
	return &ProductView{
		Id:                 product.Id,
		Name:               product.Name,
		Description:        product.Description,
		Sex:                product.Sex,
		Price:              product.Price,
		DiscountPercentage: product.DiscountPercentage,
		Stock:              product.Stock,
		ImageURL:           product.ImageURL,
		CategoryId:         product.CategoryId,
		CategoryName:       category.Name,
		BrandId:            product.BrandId,
		BrandName:          brand.Name,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
	}
}

func ToListProductView(products []model.Product, categories []model.Category, brands []model.Brand) []ProductView {
	categoryMap := map[int64]*model.Category{}
	for _, category := range categories {
		categoryMap[category.Id] = &category
	}

	brandMap := map[int64]*model.Brand{}
	for _, brand := range brands {
		brandMap[brand.Id] = &brand
	}

	productViews := make([]ProductView, len(products))
	for i, product := range products {
		productViews[i] = *ToProductView(&product, categoryMap[product.CategoryId], brandMap[product.BrandId])
	}

	return productViews
}

// Send

func FromProductViewToProductProto(productView *ProductView) *catalogservicepb.Product {
	return &catalogservicepb.Product{
		Id:                 productView.Id,
		Name:               productView.Name,
		Description:        productView.Description,
		Sex:                productView.Sex,
		Price:              productView.Price,
		DiscountPercentage: productView.DiscountPercentage,
		Stock:              productView.Stock,
		ImageUrl:           productView.ImageURL,
		CategoryId:         productView.CategoryId,
		CategoryName:       productView.CategoryName,
		BrandId:            productView.BrandId,
		BrandName:          productView.BrandName,
		CreatedAt:          timestamppb.New(productView.CreatedAt),
		UpdatedAt:          timestamppb.New(productView.UpdatedAt),
	}
}

func FromListProductViewToListProductProto(productViews []ProductView) []*catalogservicepb.Product {
	productProtos := make([]*catalogservicepb.Product, len(productViews))
	for i, prproductView := range productViews {
		productProtos[i] = FromProductViewToProductProto(&prproductView)
	}

	return productProtos
}

// Receive

func FromProductProtoToProductView(productProto *elasticsearchservicepb.Product) *ProductView {
	return &ProductView{
		Id:                 productProto.Id,
		Name:               productProto.Name,
		Description:        productProto.Description,
		Sex:                productProto.Sex,
		Price:              productProto.Price,
		DiscountPercentage: productProto.DiscountPercentage,
		Stock:              productProto.Stock,
		ImageURL:           productProto.ImageUrl,
		CategoryId:         productProto.CategoryId,
		CategoryName:       productProto.CategoryName,
		BrandId:            productProto.BrandId,
		BrandName:          productProto.BrandName,
		CreatedAt:          productProto.CreatedAt.AsTime(),
		UpdatedAt:          productProto.UpdatedAt.AsTime(),
	}
}

func FromListProductProtoToListProductView(productProtos []*elasticsearchservicepb.Product) []ProductView {
	productViews := make([]ProductView, len(productProtos))
	for i, productProto := range productProtos {
		productViews[i] = *FromProductProtoToProductView(productProto)
	}

	return productViews
}
