package model

import (
	"time"

	"github.com/uptrace/bun"
)

type CartItem struct {
	bun.BaseModel `bun:"table:cart_items"`

	Id        int64     `bun:"id,pk,autoincrement"`
	UserId    int64     `bun:"user_id,notnull"`
	ProductId int64     `bun:"product_id,notnull"`
	Quantity  int32     `bun:"quantity,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
