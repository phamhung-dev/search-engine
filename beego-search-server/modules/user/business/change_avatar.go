package userbusiness

import (
	"beego-search-server/common"
	"beego-search-server/component/objectstoragepvd"
	"beego-search-server/models"
	"context"
	"encoding/json"
)

type ChangeAvatarStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
	Update(context context.Context, user *models.User) (*models.User, error)
}

type changeAvatarBusiness struct {
	storage ModifyStorage
}

func NewChangeAvatarBusiness(storage ChangeAvatarStorage) *changeAvatarBusiness {
	return &changeAvatarBusiness{storage: storage}
}

func (business *changeAvatarBusiness) ChangeAvatar(context context.Context, objectStorageProvider objectstoragepvd.ObjectStorageProvider, id uint64, data *models.UserAvatar) (*models.User, error) {
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

	if err := data.UploadAvatar(context, objectStorageProvider); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonData, user); err != nil {
		return nil, err
	}

	user, err = business.storage.Update(context, user)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(models.UserEntityName, err)
	}

	user.Mask(false)

	return user, nil
}
