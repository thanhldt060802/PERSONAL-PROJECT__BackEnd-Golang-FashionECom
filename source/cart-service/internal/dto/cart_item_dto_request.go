package dto

type GetAllCartItemsByUserIdRequest struct {
	UserId int64  `path:"user_id" required:"true" doc:"User id of cart item."`
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type CreateCartItemRequest struct {
	Body struct {
		UserId    int64 `json:"user_id" required:"true" minimun:"1" doc:"User id of cart item."`
		ProductId int64 `json:"product_id" required:"true" minimun:"1" doc:"Product id of cart item."`
	}
}

type UpdateCartItemByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of cart item."`
	Body struct {
		Quantity *int32 `json:"quantity,omitempty" minimun:"1" doc:"Quantiy of cart item."`
	}
}

type DeleteCartItemByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of cart item."`
}

type CreateAccountCartItemRequest struct {
	Body struct {
		ProductId int64 `json:"product_id" required:"true" minimun:"1" doc:"Product id of account cart item."`
	}
}

type UpdateAccountCartItemByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of account cart item."`
	Body struct {
		Quantity *int32 `json:"quantity,omitempty" minimun:"1" doc:"Quantiy of account cart item."`
	}
}

type DeleteAccountCartItemByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of account cart item."`
}
