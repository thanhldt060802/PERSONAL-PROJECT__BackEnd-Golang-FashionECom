package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Invoice struct {
	bun.BaseModel `bun:"table:invoices"`

	Id          int64     `bun:"id,pk,autoincrement"`
	UserId      int64     `bun:"user_id,notnull"`
	TotalAmount int64     `bun:"total_amount,notnull"`
	Status      string    `bun:"status,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

// Integrate with Elasticsearch

var InvoiceSchemaMappingElasticsearch = `
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
	  "user_id": { "type": "long" },
	  "total_amount": { "type": "long" },
	  "status": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
	  "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`

var InvoiceSchemaMappingElasticsearchSortFieldMap = map[string]string{
	"id":           "id",
	"user_id":      "user_id",
	"total_amount": "total_amount",
	"status":       "status.keyword",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}
