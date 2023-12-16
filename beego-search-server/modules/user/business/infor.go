package userbusiness

import (
	"beego-search-server/common"
	"beego-search-server/models"
	"context"
)

type InforStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
}

type inforBusiness struct {
	storage InforStorage
}

func NewInforBusiness(storage InforStorage) *inforBusiness {
	return &inforBusiness{storage: storage}
}

func (business *inforBusiness) Infor(context context.Context, id uint64) (*models.User, error) {
	user := new(models.User)
	user.ID = id
	user, err := business.storage.Find(context, user, nil)

	if err == common.ErrRecordNotFound {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}

	if err != nil {
		return nil, common.ErrDB(err)
	}

	if user.IsLocked {
		return nil, models.ErrUserIsLocked
	}

	user.Mask(false)

	return user, nil
}
