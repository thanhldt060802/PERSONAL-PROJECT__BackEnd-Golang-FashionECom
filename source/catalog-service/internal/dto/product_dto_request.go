package dto

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
