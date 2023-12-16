package userstorage

import (
	"beego-search-server/models"
	"context"
)

func (s *storage) Update(context context.Context, user *models.User) (*models.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	if _, errUpdate := tx.Update(user); errUpdate != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errRollback
		}
		return nil, errUpdate
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errRollback
		}
		return nil, errCommit
	}

	return user, nil
}
