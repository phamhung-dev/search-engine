package storage

import (
	"context"
	"linux-search-bot/models"
)

func (s *storage) Delete(context context.Context, file *models.File) error {
	tx := s.db.Begin()

	if err := tx.Table(models.FileTableName).Delete(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
