package userstorage

import (
	"beego-search-server/common"
	"beego-search-server/models"
	"context"

	"github.com/beego/beego/v2/client/orm"
)

func (s *storage) Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error) {
	user.IsDeleted = false

	if cols == nil {
		cols = []string{"ID"}
	}
	cols = append(cols, "IsDeleted")
	if err := s.db.Read(user, cols...); err != nil {
		if err == orm.ErrNoRows {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}

	for i := range relmd {
		if _, err := s.db.LoadRelated(user, relmd[i]); err != nil {
			return nil, err
		}
	}
	return user, nil
}
