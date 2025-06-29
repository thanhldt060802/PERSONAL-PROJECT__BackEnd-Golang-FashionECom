package schema

var Invoice = `
{
  "mappings": {
    "properties": {
      "id": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "user_id": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
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

var InvoiceStandardizeSortFieldMap = map[string]string{
	"id":           "id.keyword",
	"user_id":      "user_id.keyword",
	"total_amount": "total_amount",
	"status":       "status.keyword",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}
