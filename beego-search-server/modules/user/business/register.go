package userbusiness

import (
	"beego-search-server/common"
	"beego-search-server/models"
	"context"
	"encoding/json"
)

type RegisterStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
	Create(context context.Context, user *models.User) (*models.User, error)
}

type registerBusiness struct {
	storage RegisterStorage
}

func NewRegisterBusiness(storage RegisterStorage) *registerBusiness {
	return &registerBusiness{storage: storage}
}

func (business *registerBusiness) Register(context context.Context, data *models.UserCreate) (*models.User, error) {
	if err := data.Validate(); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	user := new(models.User)
	user.Email = data.Email
	cols := []string{"Email"}
	_, err := business.storage.Find(context, user, cols)

	if err == nil {
		return nil, models.ErrEmailExisted
	}

	if err != common.ErrRecordNotFound {
		return nil, common.ErrDB(err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonData, user); err != nil {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	user, err = business.storage.Create(context, user)

	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.UserEntityName, err)
	}

	user.Mask(false)

	return user, nil
}
