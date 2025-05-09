package schema

var User = `
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "full_name": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "email": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "username": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "address": {
          "type": "text",
          "analyzer": "standard",
          "fields": {
            "keyword": { "type": "keyword" }
          }
        },
      "role_name": {
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

var UserStandardizeSortFieldMap = map[string]string{
	"id":         "id",
	"full_name":  "full_name.keyword",
	"email":      "email.keyword",
	"username":   "username.keyword",
	"address":    "address.keyword",
	"role_name":  "role_name.keyword",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
