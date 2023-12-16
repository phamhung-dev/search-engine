package searchstorage

import (
	"beego-search-server/models"
	"beego-search-server/utils"
	"context"
)

func (s *storage) List(context context.Context, filter *models.FileFilter, paging *utils.Paging, relmd ...string) ([]models.File, error) {
	var files []models.File

	db := s.db.QueryTable(models.FileTableName)

	if filter.Extension != "" {
		db = db.Filter("extension", filter.Extension)
	}

	if filter.SizeMin > 0 {
		db = db.Filter("size__gte", filter.SizeMin)
	}

	if filter.SizeMax > 0 {
		db = db.Filter("size__lte", filter.SizeMax)
	}

	if filter.FileCreatedAtStartTime != nil {
		db = db.Filter("file_created_at__gte", *filter.FileCreatedAtStartTime)
	}

	if filter.FileCreatedAtEndTime != nil {
		db = db.Filter("file_created_at__lte", *filter.FileCreatedAtEndTime)
	}

	if filter.FileLastModifiedAtStartTime != nil {
		db = db.Filter("file_last_modified_at__gte", *filter.FileLastModifiedAtStartTime)
	}

	if filter.FileLastModifiedAtEndTime != nil {
		db = db.Filter("file_last_modified_at__lte", *filter.FileLastModifiedAtEndTime)
	}

	if filter.FileLastAccessedAtStartTime != nil {
		db = db.Filter("file_last_accessed_at__gte", *filter.FileLastAccessedAtStartTime)
	}

	if filter.FileLastAccessedAtEndTime != nil {
		db = db.Filter("file_last_accessed_at__lte", *filter.FileLastAccessedAtEndTime)
	}

	if filter.Name != "" {
		db = db.Filter("name__icontains", filter.Name)
	}

	if filter.Content != "" {
		db = db.Filter("content__icontains", filter.Content)
	}

	var err error

	if paging.Total, err = db.Count(); err != nil {
		return nil, err
	}

	if paging.FakeCursor != "" {
		id, _, err := utils.DecodeUID(paging.FakeCursor)

		if err != nil {
			return nil, err
		}

		db = db.Filter("created_at__lt", id)
	} else {
		offset := (paging.Page - 1) * paging.Limit

		db = db.Offset(offset)
	}

	if _, err := db.Limit(paging.Limit).OrderBy("-created_at").All(&files); err != nil {
		return nil, err
	}

	filesLength := len(files)

	if filesLength > 0 {
		last := files[filesLength-1]

		paging.NextCursor = utils.EncodeUID(last.ID, models.FileTableID)
	}

	return files, nil
}
