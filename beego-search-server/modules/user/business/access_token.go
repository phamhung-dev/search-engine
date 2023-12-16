package userbusiness

import (
	"beego-search-server/common"
	"beego-search-server/component/tokenpvd"
	"beego-search-server/models"
	"beego-search-server/utils"
	"context"
)

type AccessTokenStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
}

type accessTokenBusiness struct {
	storage       AccessTokenStorage
	tokenProvider tokenpvd.TokenProvider
}

func NewAccessTokenBusiness(storage AccessTokenStorage, tokenProvider tokenpvd.TokenProvider) *accessTokenBusiness {
	return &accessTokenBusiness{storage: storage, tokenProvider: tokenProvider}
}

func (business *accessTokenBusiness) AccessToken(context context.Context, data string) (*models.User, error) {
	payload, err := business.tokenProvider.ValidateAccessToken(data)

	if err != nil {
		return nil, err
	}

	id, _, err := utils.DecodeUID(payload.UserId)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	user.ID = id
	user, err = business.storage.Find(context, user, nil)

	if err == common.ErrRecordNotFound {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}

	if err != nil {
		return nil, common.ErrDB(err)
	}

	if user.IsLocked {
		return nil, common.ErrEntityIsLocked(models.UserEntityName, models.ErrUserIsLocked)
	}

	user.Mask(false)

	return user, nil
}
