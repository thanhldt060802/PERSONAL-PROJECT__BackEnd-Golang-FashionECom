package model

import (
	"thanhldt060802/internal/grpc/client/elasticsearchservicepb"
	"thanhldt060802/internal/grpc/service/catalogservicepb"
	"time"

	"github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Product struct {
	bun.BaseModel `bun:"tb_product"`

	Id                 string     `bun:"id,pk"`
	Name               string     `bun:"name,notnull"`
	Description        string     `bun:"description,notnull"`
	Sex                string     `bun:"sex,notnull"`
	Price              int64      `bun:"price,notnull"`
	DiscountPercentage int32      `bun:"discount_percentage,notnull"`
	Stock              int32      `bun:"stock,notnull"`
	ImageURL           string     `bun:"image_url,notnull"`
	CategoryId         string     `bun:"category_id,notnull"`
	BrandId            string     `bun:"brand_id,notnull"`
	CreatedAt          *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type ProductView struct {
	bun.BaseModel `bun:"tb_product,alias:_product"`

	Id                 string    `json:"id" bun:"id,pk"`
	Name               string    `json:"name" bun:"name"`
	Description        string    `json:"description" bun:"description"`
	Sex                string    `json:"sex" bun:"sex"`
	Price              int64     `json:"price" bun:"price"`
	DiscountPercentage int32     `json:"discount_percentage" bun:"discount_percentage"`
	Stock              int32     `json:"stock" bun:"stock"`
	ImageURL           string    `json:"image_url" bun:"image_url"`
	CategoryId         string    `json:"category_id" bun:"category_id"`
	CategoryName       string    `json:"category_name" bun:"category_name"`
	BrandId            string    `json:"brand_id" bun:"brand_id"`
	BrandName          string    `json:"brand_name" bun:"brand_name"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at"`
}

// View -> Proto

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

func FromListProductViewToListProductProto(productViews []*ProductView) []*catalogservicepb.Product {
	productProtos := make([]*catalogservicepb.Product, len(productViews))
	for i, productView := range productViews {
		productProtos[i] = FromProductViewToProductProto(productView)
	}

	return productProtos
}

// Proto -> View

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

func FromListProductProtoToListProductView(productProtos []*elasticsearchservicepb.Product) []*ProductView {
	productViews := make([]*ProductView, len(productProtos))
	for i, productProto := range productProtos {
		productViews[i] = FromProductProtoToProductView(productProto)
	}

	return productViews
}
