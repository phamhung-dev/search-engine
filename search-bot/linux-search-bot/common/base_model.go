package common

import (
	"time"
)

type BaseModel struct {
	ID        uint64    `gorm:"column:id;type:serial;primary_key"`
	CreatedAt time.Time `gorm:"column:created_at;index:idx_created_at;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}
