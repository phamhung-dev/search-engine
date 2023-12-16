package storage

import (
	"context"
	"windows-search-bot/models"
)

func (s *storage) ListFilesInFolder(context context.Context, folderName string) ([]models.File, error) {
	var files []models.File

	if err := s.db.Table(models.FileTableName).Where(`regexp_replace(path, '(.*)[\\|/].*', '\1') like ? escape(':')`, folderName).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}
