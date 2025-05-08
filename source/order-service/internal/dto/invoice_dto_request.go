package dto

type GetInvoicesRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"total_amount:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=total_amount:desc,created_at will sort by total_amount in descending order, then by created_at in ascending order."`
}

type GetInvoiceByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice item."`
}

type GetInvoicesByUserIdRequest struct {
	UserId int64  `path:"user_id" required:"true" doc:"In of user."`
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"total_amount:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=total_amount:desc,created_at will sort by total_amount in descending order, then by created_at in ascending order."`
}

type CreateInvoiceRequest struct {
	Body struct {
		UserId      int64 `json:"user_id" required:"true" minimum:"1" doc:"User id of invoice."`
		TotalAmount int64 `json:"total_amount" required:"true" minimum:"0" doc:"Total amount of invoice."`
		Details     []struct {
			ProductId          int64 `json:"product_id" required:"true" minimum:"1" doc:"Product id of invoice detail."`
			Price              int64 `json:"product_price" required:"true" minimum:"0" doc:"Product price of invoice detail."`
			DiscountPercentage int32 `json:"discount_percentage" required:"true" minimum:"0" doc:"Product discount percentage of invoice detail."`
			Quantity           int32 `json:"quantity" required:"true" minimum:"1" doc:"Product Quantity of invoice detail."`
			TotalPrice         int64 `json:"total_price" required:"true" minimum:"0" doc:"Product total price id of invoice detail."`
		} `json:"details" required:"true" doc:"Details of invoice."`
	}
}

type UpdateInvoiceByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of invoice."`
	Body struct {
		Status *string `json:"status,omitempty" minimum:"1" enum:"CREATED,PENDING,CANCEL,DONE" doc:"Status of invoice."`
	}
}

type DeleteInvoiceByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice."`
}

// Integrate with Elasticsearch

type GetInvoicesWithElasticsearchRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"created_at:asc" example:"total_amount:desc,created_at" doc:"Sort by one or more fields separated by commas. For example: sort_by=total_amount:desc,created_at will sort by total_amount in descending order, then by created_at in ascending order."`
	// Filter
	Status         string `query:"status" enum:"CREATED,PENDING,CANCEL,DONE" example:"CREATED" doc:"Filter by status."`
	TotalAmountGTE string `query:"total_amount_gte" pattern:"^[0-9]+$" example:"250000" doc:"Filter by total amount greater than or equal."`
	TotalAmountLTE string `query:"total_amount_lte" pattern:"^[0-9]+$" example:"500000" doc:"Filter by total amount less than or equal."`
	CreatedAtGTE   string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE   string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}
