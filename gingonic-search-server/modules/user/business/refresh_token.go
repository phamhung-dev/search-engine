package userbusiness

import (
	"context"
	"gingonic-search-server/common"
	"gingonic-search-server/component/tokenpvd"
	"gingonic-search-server/models"
	"gingonic-search-server/utils"
	"os"
	"strconv"
)

type RefreshTokenStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
}

type refreshTokenBusiness struct {
	storage       RefreshTokenStorage
	tokenProvider tokenpvd.TokenProvider
	expiredIn     int
}

func NewRefreshTokenBusiness(storage RefreshTokenStorage, tokenProvider tokenpvd.TokenProvider) *refreshTokenBusiness {
	expiredIn, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))

	if err != nil {
		expiredIn = 60 * 60 * 24 * 30
	}

	return &refreshTokenBusiness{
		storage:       storage,
		tokenProvider: tokenProvider,
		expiredIn:     expiredIn,
	}
}

func (business *refreshTokenBusiness) RefreshToken(context context.Context, data string) (*tokenpvd.Token, error) {
	payload, err := business.tokenProvider.ValidateRefreshToken(data)

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
	newPayload := tokenpvd.TokenPayload{
		UserId: user.FakeID,
	}

	token, err := business.tokenProvider.Generate(newPayload, business.expiredIn)

	if err != nil {
		return nil, err
	}

	return token, nil
}
