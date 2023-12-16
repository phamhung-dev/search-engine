package models

import (
	"beego-search-server/common"
	"beego-search-server/utils"
	"errors"
	"time"
)

const (
	FileEntityName = "File"
	FileTableName  = "files"
	FileTableID    = 2
)

type File struct {
	common.BaseModel   `json:",inline"`
	Name               string    `json:"name" orm:"column(name)"`
	Path               string    `json:"path" orm:"column(path)"`
	Extension          string    `json:"extension" orm:"column(extension);index(idx_extension)"`
	Size               int64     `json:"size" orm:"column(size);index(idx_size)"`
	Content            string    `json:"content" orm:"column(content)"`
	ReadOnly           bool      `json:"read_only" orm:"column(read_only);default(false)"`
	Hidden             bool      `json:"hidden" orm:"column(hidden);default(false)"`
	FileCreatedAt      time.Time `json:"file_created_at" orm:"column(file_created_at);index(idx_file_created_at)"`
	FileLastModifiedAt time.Time `json:"file_last_modified_at" orm:"column(file_last_modified_at);index(idx_file_last_modified_at)"`
	FileLastAccessedAt time.Time `json:"file_last_accessed_at" orm:"column(file_last_accessed_at);index(idx_file_last_accessed_at)"`
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
