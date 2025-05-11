package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Brand struct {
	bun.BaseModel `bun:"table:brands"`

	Id          int64     `bun:"id,pk,autoincrement"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
