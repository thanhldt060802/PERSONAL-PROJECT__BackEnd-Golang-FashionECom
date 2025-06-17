package model

import (
	"github.com/uptrace/bun"
)

type CartItem struct {
	bun.BaseModel `bun:"table:cart_items"`

	Id        string `bun:"id,pk"`
	UserId    string `bun:"user_id,notnull"`
	ProductId string `bun:"product_id,notnull"`
	Quantity  int32  `bun:"quantity,notnull"`
}
