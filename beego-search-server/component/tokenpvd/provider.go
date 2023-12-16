package tokenpvd

import (
	"beego-search-server/common"
	"errors"
	"time"
)

type TokenProvider interface {
	Generate(payload TokenPayload, expiredIn int) (*Token, error)
	ValidateAccessToken(token string) (*TokenPayload, error)
	ValidateRefreshToken(token string) (*TokenPayload, error)
}

type TokenPayload struct {
	UserId string `json:"user_id"`
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiredIn    int       `json:"expired_in"`
	RefreshToken string    `json:"-"`
}

var (
	ErrProviderIsNotConfigured = common.NewCustomErrorResponse(
		errors.New("token provider is not configured"),
		"token provider is not configured",
		"ErrProviderIsNotConfigured",
	)

	ErrTokenNotFound = common.NewCustomErrorResponse(
		errors.New("token not found"),
		"token not found",
		"ErrTokenNotFound",
	)

	ErrInvalidToken = common.NewCustomErrorResponse(
		errors.New("invalid token"),
		"invalid token",
		"ErrInvalidToken",
	)

	ErrEncodingToken = common.NewCustomErrorResponse(
		errors.New("encoding token error"),
		"encoding token error",
		"ErrEncodingToken",
	)
)
