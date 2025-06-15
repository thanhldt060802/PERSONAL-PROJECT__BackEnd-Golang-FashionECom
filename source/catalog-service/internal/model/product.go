package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	Id                 string     `bun:"id,pk"`
	Name               string     `bun:"name,notnull"`
	Description        string     `bun:"description,notnull"`
	Sex                string     `bun:"sex,notnull"`
	Price              int64      `bun:"price,notnull"`
	DiscountPercentage int32      `bun:"discount_percentage,notnull"`
	Stock              int32      `bun:"stock,notnull"`
	ImageURL           string     `bun:"image_url,notnull"`
	CategoryId         string     `bun:"category_id,notnull"`
	BrandId            string     `bun:"brand_id,notnull"`
	CreatedAt          *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
