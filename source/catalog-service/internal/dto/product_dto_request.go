package dto

//
//
// Main features
// ######################################################################################

type GetProductByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of broduct."`
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
		CategoryId         string `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
		BrandId            string `json:"brand_id" required:"true" minimum:"1" doc:"Brand id of product."`
	}
}

type UpdateProductByIdRequest struct {
	Id   string `path:"id" required:"true" doc:"Id of broduct."`
	Body struct {
		Name               *string `json:"name,omitempty" minLength:"1" doc:"Name of broduct."`
		Description        *string `json:"description,omitempty" minLength:"1" doc:"Description of broduct."`
		Sex                *string `json:"sex,omitempty" minLength:"1" enum:"MALE,FEMALE,UNISEX" doc:"Sex of product."`
		Price              *int64  `json:"price,omitempty" minimum:"0" doc:"Price of product."`
		DiscountPercentage *int32  `json:"discount_percentage,omitempty" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              *int32  `json:"stock,omitempty" minimun:"0" doc:"Stock of product."`
		ImageURL           *string `json:"image_url,omitempty" minLength:"1" doc:"Image URL of product."`
		CategoryId         *string `json:"category_id,omitempty" minimum:"1" doc:"Category id of product."`
		BrandId            *string `json:"brand_id,omitempty" minimum:"1" doc:"Brand id of product."`
	}
}

type DeleteProductByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of broduct."`
}

//
//
// Elasticsearch integration features
// ######################################################################################

type GetProductsRequest struct {
	Offset int32  `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int32  `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=name:desc,created_at will sort by name in descending order, then by created_at in ascending order."`
	// Filter
	CategoryId string `query:"category_id" example:"aaaaaaaa-bbbb-cccc-dddddddd" doc:"Filter by category id."`
	BrandId    string `query:"brand_id" example:"aaaaaaaa-bbbb-cccc-dddddddd" doc:"Filter by brand id."`
	// Search
	Name                  string `query:"name" example:"Quần A1" doc:"Search by name."`
	Description           string `query:"description" example:"Quá xịn" doc:"Search by description."`
	Sex                   string `query:"sex" example:"MALE" enum:"MALE,FEMALE,UNISEX" doc:"Search by sexx."`
	PriceGTE              string `query:"price_gte" pattern:"^[0-9]+$" example:"100000" doc:"Search by price greater than or equals."`
	PriceLTE              string `query:"price_lte" pattern:"^[0-9]+$" example:"200000" doc:"Search by price less than or equals."`
	DiscountPercentageGTE string `query:"discount_percentage_gte" pattern:"^[0-9]+$" example:"20" doc:"Search by discount percentage greater than or equals."`
	DiscountPercentageLTE string `query:"discount_percentage_lte" pattern:"^[0-9]+$" example:"30" doc:"Search by discount percentage less than or equals."`
	StockGTE              string `query:"stock_gte" pattern:"^[0-9]+$" example:"50" doc:"Search by stock greater than or equals."`
	StockLTE              string `query:"stock_lte" pattern:"^[0-9]+$" example:"100" doc:"Search by stock less than or equals."`
	CategoryName          string `query:"category_name" example:"Quần" doc:"Search by category name."`
	BrandName             string `query:"brand_name" example:"Gucci" doc:"Search by brand name."`
	CreatedAtGTE          string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Search by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE          string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Search by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}
