package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Invoice struct {
	bun.BaseModel `bun:"table:invoices"`

	Id          string     `bun:"id,pk"`
	UserId      string     `bun:"user_id,notnull"`
	TotalAmount int64      `bun:"total_amount,notnull"`
	Status      string     `bun:"status,notnull"`
	CreatedAt   *time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   *time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
