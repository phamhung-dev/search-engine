package storage

import (
	"context"
	"windows-search-bot/models"
)

func (s *storage) Update(context context.Context, file *models.File) error {
	tx := s.db.Begin()

	if err := tx.Table(models.FileTableName).Save(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
