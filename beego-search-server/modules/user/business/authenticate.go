package userbusiness

import (
	"beego-search-server/common"
	"beego-search-server/component/tokenpvd"
	"beego-search-server/models"
	"context"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticateStorage interface {
	Find(context context.Context, user *models.User, cols []string, relmd ...string) (*models.User, error)
}

type authenticateBusiness struct {
	storage       AuthenticateStorage
	tokenProvider tokenpvd.TokenProvider
	expiredIn     int
}

func NewAuthenticateBusiness(storage AuthenticateStorage, tokenProvider tokenpvd.TokenProvider) *authenticateBusiness {
	expiredIn, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))

	if err != nil {
		expiredIn = 60 * 60 * 24 * 30
	}

	return &authenticateBusiness{
		storage:       storage,
		tokenProvider: tokenProvider,
		expiredIn:     expiredIn,
	}
}

func (business *authenticateBusiness) Authenticate(context context.Context, data *models.UserLogin) (*tokenpvd.Token, error) {
	if err := data.Validate(); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	user := new(models.User)
	user.Email = data.Email
	cols := []string{"Email"}
	user, err := business.storage.Find(context, user, cols)

	if err != nil {
		return nil, models.ErrEmailOrPasswordIsIncorrect
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, models.ErrEmailOrPasswordIsIncorrect
	}

	if user.IsLocked {
		return nil, common.ErrEntityIsLocked(models.UserEntityName, models.ErrUserIsLocked)
	}

	user.Mask(false)
	payload := tokenpvd.TokenPayload{
		UserId: user.FakeID,
	}

	token, err := business.tokenProvider.Generate(payload, business.expiredIn)

	if err != nil {
		return nil, err
	}

	return token, nil
}
