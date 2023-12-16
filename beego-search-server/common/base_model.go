package common

import (
	"time"
)

type BaseModel struct {
	ID        uint64    `json:"-" orm:"column(id);pk;auto"`
	FakeID    string    `json:"id" orm:"-"`
	CreatedAt time.Time `json:"created_at" orm:"column(created_at);index(idx_created_at);auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);auto_now"`
}
