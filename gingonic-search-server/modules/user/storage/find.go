package userstorage

import (
	"context"
	"gingonic-search-server/common"
	"gingonic-search-server/models"

	"gorm.io/gorm"
)

func (s *storage) Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error) {
	user.IsDeleted = false

	if cols == nil {
		cols = []string{"ID"}
	}

	db := s.db.Table(models.UserTableName)
	for i := range relmd {
		db = db.Preload(relmd[i])
	}

	if err := db.Where("is_deleted", user.IsDeleted).Where(user, cols).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}
