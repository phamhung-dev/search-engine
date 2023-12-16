package storage

import (
	"context"
	"windows-search-bot/models"
)

func (s *storage) Find(context context.Context, conditions map[string]interface{}) (*models.File, error) {
	var file models.File

	if err := s.db.Table(models.FileTableName).Where(conditions).First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

func (s *storage) FindLastFileCreated(context context.Context) (*models.File, error) {
	var file models.File

	if err := s.db.Table(models.FileTableName).Order("created_at desc").First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}
