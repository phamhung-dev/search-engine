package authorizepvd

import (
	"beego-search-server/common"
	"beego-search-server/models"
	"errors"
)

type AuthorizeProvider interface {
	ValidateRequest(user *models.User, path string, method string) (bool, error)
}

var (
	ErrProviderIsNotConfigured = common.NewCustomErrorResponse(
		errors.New("authorize provider is not configured"),
		"authorize provider is not configured",
		"ErrProviderIsNotConfigured",
	)
)
