package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Brand struct {
	bun.BaseModel `bun:"tb_brand"`

	Id          string     `bun:"id,pk"`
	Name        string     `bun:"name,notnull"`
	Description string     `bun:"description,notnull"`
	CreatedAt   *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type BrandView struct {
	bun.BaseModel `bun:"tb_brand,alias:_brand"`

	Id          string    `json:"id" bun:"id,pk"`
	Name        string    `json:"name" bun:"name"`
	Description string    `json:"description" bun:"description"`
	CreatedAt   time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bun:"updated_at"`
}
