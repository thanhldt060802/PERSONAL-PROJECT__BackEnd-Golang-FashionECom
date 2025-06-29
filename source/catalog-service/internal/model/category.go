package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"tb_category"`

	Id        string     `bun:"id,pk"`
	Name      string     `bun:"name,notnull"`
	CreatedAt *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type CategoryView struct {
	bun.BaseModel `bun:"tb_category,alias:_category"`

	Id        string    `json:"id" bun:"id,pk"`
	Name      string    `json:"name" bun:"name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`
}
