package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gertd/go-pluralize"
)

var plrl = pluralize.NewClient()

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

type SuccessResponse struct {
	Data   interface{} `json:"data"`
	Filter interface{} `json:"filter,omitempty"`
	Paging interface{} `json:"paging,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.RootError().Error()
}

func (e *ErrorResponse) RootError() error {
	if err, ok := e.RootErr.(*ErrorResponse); ok {
		return err.RootError()
	}

	return e.RootErr
}

func NewErrorResponse(statusCode int, root error, msg, log, key string) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewCustomErrorResponse(root error, msg, key string) *ErrorResponse {
	if root != nil {
		return NewErrorResponse(http.StatusBadRequest, root, msg, root.Error(), key)
	}

	return NewErrorResponse(http.StatusBadRequest, errors.New(msg), msg, msg, key)
}

func ErrDB(err error) *ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, err, "something went wrong", err.Error(), "ErrDB")
}

func ErrInvalidRequest(err error) *ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, err, "invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternalServer(err error) *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, err, "something went wrong", err.Error(), "ErrInternalServer")
}

func ErrCannotListEntity(entity string, err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("cannot list %s", plrl.Plural(strings.ToLower(entity))),
		err.Error(),
		fmt.Sprintf("ErrCannotList%s", entity),
	)
}

func ErrCannotCreateEntity(entity string, err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("cannot create %s", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("ErrCannotCreate%s", entity),
	)
}

func ErrCannotUpdateEntity(entity string, err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("cannot update %s", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("ErrCannotUpdate%s", entity),
	)
}

func ErrEntityNotFound(entity string, err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusNotFound,
		err,
		fmt.Sprintf("%s not found", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sNotFound", entity),
	)
}

func ErrPermissionDenied(err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusForbidden,
		err,
		"permission denied",
		err.Error(),
		"ErrPermissionDenied",
	)
}

func ErrEntityIsLocked(entity string, err error) *ErrorResponse {
	return NewErrorResponse(
		http.StatusLocked,
		err,
		fmt.Sprintf("%s is locked", strings.ToLower(entity)),
		err.Error(),
		fmt.Sprintf("Err%sIsLocked", entity),
	)
}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{Data: data, Paging: nil, Filter: nil}
}

func NewCustomSuccessResponse(data, filter, paging interface{}) *SuccessResponse {
	return &SuccessResponse{Data: data, Filter: filter, Paging: paging}
}
