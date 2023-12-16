package storage

import (
	"context"
	"linux-search-bot/common"
	"linux-search-bot/models"

	"gorm.io/gorm"
)

func (s *storage) Create(context context.Context, file *models.File) error {
	if file == nil {
		return common.ErrDataNotFound
	}

	tx := s.db.Begin()

	if err := tx.Table(models.FileTableName).Create(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *storage) BatchCreate(context context.Context, files []models.File) error {
	if files == nil {
		return common.ErrDataNotFound
	}

	tx := s.db.Begin()

	if err := tx.Session(&gorm.Session{SkipDefaultTransaction: true}).Table(models.FileTableName).Create(&files).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
