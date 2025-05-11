package model

import (
	"github.com/uptrace/bun"
)

type InvoiceDetail struct {
	bun.BaseModel `bun:"table:invoice_details"`

	Id                 int64 `bun:"id,pk,autoincrement"`
	InvoiceId          int64 `bun:"invoice_id,notnull"`
	ProductId          int64 `bun:"product_id,notnull"`
	Price              int64 `bun:"price,notnull"`
	DiscountPercentage int32 `bun:"discount_percentage,notnull"`
	Quantity           int32 `bun:"quantity,notnull"`
	TotalPrice         int64 `bun:"total_price,notnull"`
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
