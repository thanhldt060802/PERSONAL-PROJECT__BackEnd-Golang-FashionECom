package dto

//
//
// Main features
// ######################################################################################

type GetProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of broduct."`
}

type CreateProductRequest struct {
	Body struct {
		Name               string `json:"name" required:"true" minLength:"1" doc:"Name of product."`
		Description        string `json:"description" required:"true" minLength:"1" doc:"Description of product."`
		Sex                string `json:"sex" required:"true" minLength:"1" enum:"MALE,FEMALE,UNISEX" doc:"Sex of product."`
		Price              int64  `json:"price" required:"true" minimum:"0" doc:"Price of product."`
		DiscountPercentage int32  `json:"discount_percentage" required:"true" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              int32  `json:"stock" required:"true" minimun:"0" doc:"Stock of product."`
		ImageURL           string `json:"image_url" required:"true" minLength:"1" doc:"Image URL of product."`
		CategoryId         int64  `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
		BrandId            int64  `json:"brand_id" required:"true" minimum:"1" doc:"Brand id of product."`
	}
}

type UpdateProductByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of broduct."`
	Body struct {
		Name               *string `json:"name,omitempty" minLength:"1" doc:"Name of broduct."`
		Description        *string `json:"description,omitempty" minLength:"1" doc:"Description of broduct."`
		Sex                *string `json:"sex,omitempty" minLength:"1" enum:"MALE,FEMALE,UNISEX" doc:"Sex of product."`
		Price              *int64  `json:"price,omitempty" minimum:"0" doc:"Price of product."`
		DiscountPercentage *int32  `json:"discount_percentage,omitempty" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              *int32  `json:"stock,omitempty" minimun:"0" doc:"Stock of product."`
		ImageURL           *string `json:"image_url,omitempty" minLength:"1" doc:"Image URL of product."`
		CategoryId         *int64  `json:"category_id,omitempty" minimum:"1" doc:"Category id of product."`
		BrandId            *int64  `json:"brand_id,omitempty" minimum:"1" doc:"Brand id of product."`
	}
}

type DeleteProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of broduct."`
}

//
//
// Elasticsearch integration features
// ######################################################################################

type GetProductsRequest struct {
	Offset int32  `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int32  `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"full_name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=full_name:desc,created_at will sort by full_name in descending order, then by created_at in ascending order."`
	// Filter
	Name                  *string `query:"name" example:"Quần A1" doc:"Filter by name."`
	Description           *string `query:"description" example:"Quá xịn" doc:"Filter by description."`
	Sex                   *string `query:"sex" example:"MALE" enum:"MALE,FEMALE,UNISEX" doc:"Filter by sexx."`
	PriceGTE              *int64  `query:"price_gte" example:"100000" doc:"Filter by price greater than or equals."`
	PriceLTE              *int64  `query:"price_lte" example:"200000" doc:"Filter by price less than or equals."`
	DiscountPercentageGTE *int32  `query:"discount_percentage_gte" example:"20" doc:"Filter by discount percentage greater than or equals."`
	DiscountPercentageLTE *int32  `query:"discount_percentage_lte" example:"30" doc:"Filter by discount percentage less than or equals."`
	StockGTE              *int32  `query:"stock_gte" example:"50" doc:"Filter by stock greater than or equals."`
	StockLTE              *int32  `query:"stock_lte" example:"100" doc:"Filter by stock less than or equals."`
	CategoryName          *string `query:"category_name" example:"Quần" doc:"Filter by category name."`
	BrandName             *string `query:"brand_name" example:"Gucci" doc:"Filter by brand name."`
	CreatedAtGTE          *string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE          *string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}
