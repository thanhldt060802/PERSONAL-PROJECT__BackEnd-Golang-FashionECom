package dto

type GetCartItemsRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type CreateCartItemRequest struct {
	Body struct {
		UserId    string `json:"user_id" required:"true" minLength:"1" doc:"User id of cart item."`
		ProductId string `json:"product_id" required:"true" minLength:"1" doc:"Product id of cart item."`
	}
}

type UpdateCartItemByIdRequest struct {
	Id   string `path:"id" doc:"Id of cart item."`
	Body struct {
		Quantity *int32 `json:"quantity,omitempty" minimum:"1" doc:"Quantiy of cart item."`
	}
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type DeleteCartItemByIdRequest struct {
	Id string `path:"id" doc:"Id of cart item."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type GetMyCartItemsRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type CreateMyCartItemRequest struct {
	Body struct {
		ProductId string `json:"product_id" required:"true" minLength:"1" doc:"Product id of cart item."`
	}
}

type UpdateMyCartItemByIdRequest struct {
	Id   string `path:"id" doc:"Id of cart item."`
	Body struct {
		Quantity *int32 `json:"quantity,omitempty" minimum:"1" doc:"Quantiy of cart item."`
	}
}

type DeleteMyCartItemByIdRequest struct {
	Id string `path:"id" doc:"Id of cart item."`
}
