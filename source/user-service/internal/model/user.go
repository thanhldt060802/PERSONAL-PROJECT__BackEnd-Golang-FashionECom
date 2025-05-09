package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id             int64     `bun:"id,pk,autoincrement" json:"id"`
	FullName       string    `bun:"full_name,notnull" json:"full_name"`
	Email          string    `bun:"email,notnull" json:"email"`
	Username       string    `bun:"username,notnull" json:"username"`
	HashedPassword string    `bun:"hashed_password,notnull"`
	Address        string    `bun:"address,notnull" json:"address"`
	RoleName       string    `bun:"role_name,notnull" json:"role_name"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}
