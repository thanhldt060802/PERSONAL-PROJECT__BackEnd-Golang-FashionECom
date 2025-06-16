package schema

var Product = `
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
      "name": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "description": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "sex": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "price": { "type": "long" },
      "discount_percentage": { "type": "integer" },
	    "stock": { "type": "integer" },
	    "image_url": { "type": "text" },
	    "category_id": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "category_name": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
	    "brand_id": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "brand_name": {
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

var ProductStandardizeSortFieldMap = map[string]string{
	"id":                  "id.keyword",
	"name":                "name.keyword",
	"description":         "description.keyword",
	"sex":                 "sex.keyword",
	"price":               "price",
	"discount_percentage": "discount_percentage",
	"stock":               "stock",
	"category_id":         "category_id.keyword",
	"category_name":       "category_name.keyword",
	"brand_id":            "brand_id.keyword",
	"brand_name":          "brand_name.keyword",
	"created_at":          "created_at",
	"updated_at":          "updated_at",
}
