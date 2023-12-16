package usertransport

import (
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	userbusiness "gingonic-search-server/modules/user/business"
	userstorage "gingonic-search-server/modules/user/storage"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RefreshToken(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")

		if refreshToken == "" || err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		tokenProvider := appctx.GetTokenProvider()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewRefreshTokenBusiness(storage, tokenProvider)

		response, err := business.RefreshToken(c.Request.Context(), refreshToken)

		if err != nil {
			panic(err)
		}

		domain := os.Getenv("DOMAIN")
		if domain != "" {
			c.SetCookie("refresh_token", response.RefreshToken, response.ExpiredIn, "/", domain, false, true)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
