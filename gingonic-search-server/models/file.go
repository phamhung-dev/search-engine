package models

import (
	"errors"
	"gingonic-search-server/common"
	"gingonic-search-server/utils"
	"time"
)

const (
	FileEntityName = "File"
	FileTableName  = "files"
	FileTableID    = 2
)

type File struct {
	common.BaseModel   `json:",inline"`
	Name               string    `json:"name" gorm:"column:name;type:text;not null"`
	Path               string    `json:"path" gorm:"column:path;type:text;not null"`
	Extension          string    `json:"extension" gorm:"column:extension;type:text;index:idx_extension;not null"`
	Size               int64     `json:"size" gorm:"column:size;type:bigint;index:idx_size;not null"`
	Content            string    `json:"content" gorm:"column:content;type:text;not null"`
	ReadOnly           bool      `json:"read_only" gorm:"column:read_only;default:false;not null"`
	Hidden             bool      `json:"hidden" gorm:"column:hidden;default:false;not null"`
	FileCreatedAt      time.Time `json:"file_created_at" gorm:"column:file_created_at;index:idx_file_created_at;not null"`
	FileLastModifiedAt time.Time `json:"file_last_modified_at" gorm:"column:file_last_modified_at;index:idx_file_last_modified_at;not null"`
	FileLastAccessedAt time.Time `json:"file_last_accessed_at" gorm:"column:file_last_accessed_at;index:idx_file_last_accessed_at;not null"`
}

func (File) TableName() string { return FileTableName }

func (f *File) Mask(isAdmin bool) {
	f.FakeID = utils.EncodeUID(f.ID, FileTableID)
}

type FileFilter struct {
	Name                        string     `json:"name,omitempty" form:"name"`
	Extension                   string     `json:"extension,omitempty" form:"extension"`
	Content                     string     `json:"content,omitempty" form:"content"`
	SizeMax                     uint64     `json:"size_max,omitempty" form:"size_max"`
	SizeMin                     uint64     `json:"size_min,omitempty" form:"size_min"`
	FileCreatedAtStartTime      *time.Time `json:"file_created_at_start_time,omitempty" form:"file_created_at_start_time"`
	FileCreatedAtEndTime        *time.Time `json:"file_created_at_end_time,omitempty" form:"file_created_at_end_time"`
	FileLastAccessedAtStartTime *time.Time `json:"file_last_accessed_at_start_time,omitempty" form:"file_last_accessed_at_start_time"`
	FileLastAccessedAtEndTime   *time.Time `json:"file_last_accessed_at_end_time,omitempty" form:"file_last_accessed_at_end_time"`
	FileLastModifiedAtStartTime *time.Time `json:"file_last_modified_at_start_time,omitempty" form:"file_last_modified_at_start_time"`
	FileLastModifiedAtEndTime   *time.Time `json:"file_last_modified_at_end_time,omitempty" form:"file_last_modified_at_end_time"`
}

func (f *FileFilter) Validate() error {
	if f.SizeMax > 0 && f.SizeMin > f.SizeMax {
		return ErrRangeSizeIsInvalid
	}

	if (f.FileCreatedAtStartTime != nil && f.FileCreatedAtEndTime != nil && f.FileCreatedAtStartTime.After(*f.FileCreatedAtEndTime)) ||
		(f.FileLastModifiedAtStartTime != nil && f.FileLastModifiedAtEndTime != nil && f.FileLastModifiedAtStartTime.After(*f.FileLastModifiedAtEndTime)) ||
		(f.FileLastAccessedAtStartTime != nil && f.FileLastAccessedAtEndTime != nil && f.FileLastAccessedAtStartTime.After(*f.FileLastAccessedAtEndTime)) {
		return ErrRangeTimeIsInvalid
	}

	return nil
}

var (
	ErrRangeSizeIsInvalid = errors.New("range size is invalid")
	ErrRangeTimeIsInvalid = errors.New("range time is invalid")
)
