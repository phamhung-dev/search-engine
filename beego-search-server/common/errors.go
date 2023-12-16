package common

import "errors"

var (
	ErrDataNotFound          = errors.New("data not found")
	ErrRecordNotFound        = errors.New("record not found")
	ErrIDIsInvalid           = errors.New("id is invalid")
	ErrUserDoesNotHaveAccess = errors.New("user does not have access")
	ErrRefreshTokenIsEmpty   = errors.New("refresh token is empty")
)
