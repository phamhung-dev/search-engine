package models

import (
	"linux-search-bot/common"
	"time"
)

const (
	FileEntityName = "File"
	FileTableName  = "files"
)

type File struct {
	common.BaseModel
	Name               string    `gorm:"column:name;type:text;not null"`
	Path               string    `gorm:"column:path;type:text;not null"`
	Extension          string    `gorm:"column:extension;type:text;index:idx_extension;not null"`
	Size               int64     `gorm:"column:size;type:bigint;index:idx_size;not null"`
	Content            string    `gorm:"column:content;type:text;not null"`
	ReadOnly           bool      `gorm:"column:read_only;default:false;not null"`
	Hidden             bool      `gorm:"column:hidden;default:false;not null"`
	FileCreatedAt      time.Time `gorm:"column:file_created_at;index:idx_file_created_at;not null"`
	FileLastModifiedAt time.Time `gorm:"column:file_last_modified_at;index:idx_file_last_modified_at;not null"`
	FileLastAccessedAt time.Time `gorm:"column:file_last_accessed_at;index:idx_file_last_accessed_at;not null"`
}

func (File) TableName() string { return FileTableName }
