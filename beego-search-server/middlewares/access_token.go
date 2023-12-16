package middlewares

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	userbusiness "beego-search-server/modules/user/business"
	userstorage "beego-search-server/modules/user/storage"
	"errors"
	"strings"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func extractTokenFromHeader(authorization string) (string, error) {
	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader
	}

	return parts[1], nil
}

func AccessToken(appctx appcontext.AppContext) web.FilterFunc {
	return func(c *beecontext.Context) {
		accessToken, err := extractTokenFromHeader(c.Input.Header("Authorization"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		tokenProvider := appctx.GetTokenProvider()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewAccessTokenBusiness(storage, tokenProvider)

		response, err := business.AccessToken(c.Request.Context(), accessToken)

		if err != nil {
			panic(err)
		}

		c.Input.SetData(common.CurrentUser, response)
	}
}

var (
	ErrWrongAuthHeader = common.NewCustomErrorResponse(
		errors.New("wrong authen header"),
		"wrong authen header",
		"ErrWrongAuthHeader",
	)
)
