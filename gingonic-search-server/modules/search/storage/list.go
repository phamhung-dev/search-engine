package searchstorage

import (
	"context"
	"fmt"
	"gingonic-search-server/models"
	"gingonic-search-server/utils"
)

func (s *storage) List(context context.Context, filter *models.FileFilter, paging *utils.Paging, relmd ...string) ([]models.File, error) {
	var files []models.File

	db := s.db.Table(models.FileTableName)

	if filter.Extension != "" {
		db = db.Where("extension = ?", filter.Extension)
	}

	if filter.SizeMin > 0 {
		db = db.Where("size >= ?", filter.SizeMin)
	}

	if filter.SizeMax > 0 {
		db = db.Where("size <= ?", filter.SizeMax)
	}

	if filter.FileCreatedAtStartTime != nil {
		db = db.Where("file_created_at >= ?", *filter.FileCreatedAtStartTime)
	}

	if filter.FileCreatedAtEndTime != nil {
		db = db.Where("file_created_at <= ?", *filter.FileCreatedAtEndTime)
	}

	if filter.FileLastModifiedAtStartTime != nil {
		db = db.Where("file_last_modified_at >= ?", *filter.FileLastModifiedAtStartTime)
	}

	if filter.FileLastModifiedAtEndTime != nil {
		db = db.Where("file_last_modified_at <= ?", *filter.FileLastModifiedAtEndTime)
	}

	if filter.FileLastAccessedAtStartTime != nil {
		db = db.Where("file_last_accessed_at >= ?", *filter.FileLastAccessedAtStartTime)
	}

	if filter.FileLastAccessedAtEndTime != nil {
		db = db.Where("file_last_accessed_at <= ?", *filter.FileLastAccessedAtEndTime)
	}

	if filter.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", filter.Name))
	}

	if filter.Content != "" {
		db = db.Where("content like ?", fmt.Sprintf("%%%s%%", filter.Content))
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if paging.FakeCursor != "" {
		id, _, err := utils.DecodeUID(paging.FakeCursor)

		if err != nil {
			return nil, err
		}

		db = db.Where("created_at < ?", id)
	} else {
		offset := (paging.Page - 1) * paging.Limit

		db = db.Offset(offset)
	}

	if err := db.Limit(paging.Limit).Order("created_at desc").Find(&files).Error; err != nil {
		return nil, err
	}

	filesLength := len(files)

	if filesLength > 0 {
		last := files[filesLength-1]

		paging.NextCursor = utils.EncodeUID(last.ID, models.FileTableID)
	}

	return files, nil
}
