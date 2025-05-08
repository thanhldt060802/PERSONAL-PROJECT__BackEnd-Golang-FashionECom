package model

import (
	"github.com/uptrace/bun"
)

type InvoiceDetail struct {
	bun.BaseModel `bun:"table:invoice_details"`

	Id                 int64 `bun:"id,pk,autoincrement" json:"id"`
	InvoiceId          int64 `bun:"invoice_id,notnull" json:"invoice_id"`
	ProductId          int64 `bun:"product_id,notnull" json:"product_id"`
	Price              int64 `bun:"price,notnull" json:"price"`
	DiscountPercentage int32 `bun:"discount_percentage,notnull" json:"discount_percentage"`
	Quantity           int32 `bun:"quantity,notnull" json:"quantity"`
	TotalPrice         int64 `bun:"total_price,notnull" json:"total_price"`
}

// Integrate with Elasticsearch

var InvoiceDetailSchemaMappingElasticsearch = `
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
	  "invoice_id": { "type": "long" },
	  "product_id": { "type": "long" },
	  "price": { "type": "long" },
	  "discount_percentage": { "type": "integer" },
	  "quantity": { "type": "integer" },
	  "total_price": { "type": "long" }
    }
  }
}`
