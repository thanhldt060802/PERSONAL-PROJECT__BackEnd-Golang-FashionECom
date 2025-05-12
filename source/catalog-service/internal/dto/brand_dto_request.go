package dto

type GetBrandsRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"created_at:desc,name" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,name will sort by created_at in descending order, then by name in ascending order."`
}

type GetBrandByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of brand."`
}

type CreateBrandRequest struct {
	Body struct {
		Name        string `json:"name" required:"true" minLength:"1" doc:"Name of brand (unique)."`
		Description string `json:"description" required:"true" minLength:"1" doc:"Description of brand."`
	}
}

type UpdateBrandByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of brand."`
	Body struct {
		Name        *string `json:"name,omitempty" minLength:"1" doc:"Name of brand (unique)."`
		Description *string `json:"description,omitempty" minLength:"1" doc:"Description of brand."`
	}
}

type DeleteBrandByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of brand."`
}
