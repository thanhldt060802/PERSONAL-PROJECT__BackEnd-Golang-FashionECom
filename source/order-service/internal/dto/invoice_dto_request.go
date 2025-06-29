package dto

type GetInvoicesRequest struct {
	Offset int32  `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int32  `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"name:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=name:desc,created_at will sort by name in descending order, then by created_at in ascending order."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
	// Search
	TotalAmountGTE string `query:"total_amount_gte" pattern:"^[0-9]+$" example:"100000" doc:"Search by total amount greater than or equals."`
	TotalAmountLTE string `query:"total_amount_lte" pattern:"^[0-9]+$" example:"200000" doc:"Search by total amount less than or equals."`
	Status         string `query:"status" example:"CREATED" enum:"CREATED,PAID,CANCEL" doc:"Search by status."`
	CreatedAtGTE   string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Search by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE   string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Search by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}

type GetInvoiceByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of invoice item."`
	// Filter
	UserId string `query:"user_id" doc:"Filter by user id."`
}

type CreateInvoiceRequest struct {
	Body struct {
		UserId         string          `json:"user_id" required:"true" minimum:"1" doc:"User id of invoice."`
		InvoiceDetails []InvoiceDetail `json:"invoice_details" required:"true" doc:"Invoice details."`
	}
}
type InvoiceDetail struct {
	ProductId          string `json:"product_id" required:"true" minimum:"1" doc:"Product id of invoice detail."`
	Price              int64  `json:"product_price" required:"true" minimum:"0" doc:"Price of product of invoice detail."`
	DiscountPercentage int32  `json:"discount_percentage" required:"true" minimum:"0" doc:"Discount percentage of product of invoice detail."`
	Quantity           int32  `json:"quantity" required:"true" minimum:"1" doc:"Quantity of product of invoice detail."`
	TotalPrice         int64  `json:"total_price" required:"true" minimum:"0" doc:"Total price of product of invoice detail."`
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
	Offset int32  `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int32  `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"total_amount:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=total_amount:desc,created_at will sort by total_amount in descending order, then by created_at in ascending order."`
	// Search
	TotalAmountGTE string `query:"total_amount_gte" pattern:"^[0-9]+$" example:"100000" doc:"Search by total amount greater than or equals."`
	TotalAmountLTE string `query:"total_amount_lte" pattern:"^[0-9]+$" example:"200000" doc:"Search by total amount less than or equals."`
	Status         string `query:"status" example:"CREATED" enum:"CREATED,PAID,CANCEL" doc:"Search by status."`
	CreatedAtGTE   string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Search by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE   string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Search by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}

type GetMyInvoiceByIdRequest struct {
	Id string `path:"id" required:"true" doc:"Id of invoice item."`
}
