package userstorage

import (
	"context"
	"gingonic-search-server/models"
)

func (s *storage) Create(context context.Context, user *models.User) (*models.User, error) {
	tx := s.db.Begin()

	if err := tx.Table(models.UserTableName).Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}
