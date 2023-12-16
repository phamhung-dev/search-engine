package middlewares

import (
	"errors"
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	userbusiness "gingonic-search-server/modules/user/business"
	userstorage "gingonic-search-server/modules/user/storage"
	"strings"

	"github.com/gin-gonic/gin"
)

func extractTokenFromHeader(authorization string) (string, error) {
	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader
	}

	return parts[1], nil
}

func AccessToken(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := extractTokenFromHeader(c.GetHeader("Authorization"))

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

		c.Set(common.CurrentUser, response)

		c.Next()
	}
}

var (
	ErrWrongAuthHeader = common.NewCustomErrorResponse(
		errors.New("wrong authen header"),
		"wrong authen header",
		"ErrWrongAuthHeader",
	)
)
