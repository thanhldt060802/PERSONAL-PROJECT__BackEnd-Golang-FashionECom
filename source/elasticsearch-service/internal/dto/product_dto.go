package dto

import (
	"thanhldt060802/internal/grpc/client/catalogservicepb"
	"thanhldt060802/internal/grpc/service/elasticsearchservicepb"
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

// Receive

func FromProductProtoToProductView(productProto *catalogservicepb.Product) *ProductView {
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

// Send

func FromProductViewToProductProto(productView *ProductView) *elasticsearchservicepb.Product {
	return &elasticsearchservicepb.Product{
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

func FromListProductViewToListProductProto(productViews []ProductView) []*elasticsearchservicepb.Product {
	userProtos := make([]*elasticsearchservicepb.Product, len(productViews))
	for i := range userProtos {
		userProtos[i] = FromProductViewToProductProto(&productViews[i])
	}

	return userProtos
}
