package model

import (
	"github.com/uptrace/bun"
)

type InvoiceDetail struct {
	bun.BaseModel `bun:"table:invoice_details"`

	Id                 string `bun:"id,pk"`
	InvoiceId          string `bun:"invoice_id,notnull"`
	ProductId          string `bun:"product_id,notnull"`
	Price              int64  `bun:"price,notnull"`
	DiscountPercentage int32  `bun:"discount_percentage,notnull"`
	Quantity           int32  `bun:"quantity,notnull"`
	TotalPrice         int64  `bun:"total_price,notnull"`
}
