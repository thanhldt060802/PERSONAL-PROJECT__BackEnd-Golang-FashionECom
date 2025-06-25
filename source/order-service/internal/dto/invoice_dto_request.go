package dto

type GetInvoicesRequest struct {
}
type GetInvoiceByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of invoice item."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type CreateInvoiceRequest struct {
	Body struct {
		UserId  string `json:"user_id" required:"true" minimum:"1" doc:"User id of invoice."`
		Details []struct {
			ProductId          string `json:"product_id" required:"true" minimum:"1" doc:"Product id of invoice detail."`
			Price              int64  `json:"product_price" required:"true" minimum:"0" doc:"Price of product of invoice detail."`
			DiscountPercentage int32  `json:"discount_percentage" required:"true" minimum:"0" doc:"Discount percentage of product of invoice detail."`
			Quantity           int32  `json:"quantity" required:"true" minimum:"1" doc:"Quantity of product of invoice detail."`
			TotalPrice         int64  `json:"total_price" required:"true" minimum:"0" doc:"Total price of product of invoice detail."`
		} `json:"details" required:"true" doc:"Details of invoice."`
	}
}

type UpdateInvoiceByIdRequest struct {
	Id   string `path:"id" required:"true" doc:"Id of invoice."`
	Body struct {
		Status *string `json:"status,omitempty" minLength:"1" enum:"CREATED,PENDING,CANCEL,DONE" doc:"Status of invoice."`
	}
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type DeleteInvoiceByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of invoice."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}
type GetMyInvoicesRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"total_amount:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=total_amount:desc,created_at will sort by total_amount in descending order, then by created_at in ascending order."`
}

type GetMyInvoiceByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of invoice item."`
}
