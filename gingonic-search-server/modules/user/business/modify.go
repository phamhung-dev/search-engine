package userbusiness

import (
	"context"
	"encoding/json"
	"gingonic-search-server/common"
	"gingonic-search-server/models"
)

type ModifyStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
	Update(context context.Context, user *models.User) (*models.User, error)
}

type modifyBusiness struct {
	storage ModifyStorage
}

func NewModifyBusiness(storage ModifyStorage) *modifyBusiness {
	return &modifyBusiness{storage: storage}
}

func (business *modifyBusiness) Modify(context context.Context, id uint64, data *models.UserUpdate) (*models.User, error) {
	if err := data.Validate(); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

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

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonData, user); err != nil {
		return nil, err
	}

	if data.Password != "" {
		if err := user.HashPassword(); err != nil {
			return nil, err
		}
	}

	user, err = business.storage.Update(context, user)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(models.UserEntityName, err)
	}

	user.Mask(false)

	return user, nil
}
