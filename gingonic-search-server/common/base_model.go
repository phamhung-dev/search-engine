package common

import (
	"time"
)

type BaseModel struct {
	ID        uint64    `json:"-" gorm:"column:id;type:serial;primary_key"`
	FakeID    string    `json:"id" gorm:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index:idx_created_at;autoCreateTime;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`
}
