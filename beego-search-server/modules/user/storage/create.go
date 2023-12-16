package userstorage

import (
	"beego-search-server/models"
	"context"
)

func (s *storage) Create(context context.Context, user *models.User) (*models.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if _, errInsert := tx.Insert(user); errInsert != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errRollback
		}
		return nil, errInsert
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errRollback
		}
		return nil, errCommit
	}

	return user, nil
}
