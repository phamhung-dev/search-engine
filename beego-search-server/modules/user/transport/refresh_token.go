package usertransport

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	userbusiness "beego-search-server/modules/user/business"
	userstorage "beego-search-server/modules/user/storage"
	"net/http"
	"os"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func RefreshToken(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		refreshToken := c.Input.Cookie("refresh_token")

		if refreshToken == "" {
			panic(common.ErrInvalidRequest(common.ErrRefreshTokenIsEmpty))
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
			c.Output.Cookie("refresh_token", response.RefreshToken, response.ExpiredIn, "/", domain, false, true)
		}

		c.Output.Status = http.StatusOK
		c.Output.JSON(common.NewSuccessResponse(response), false, false)
	}
}
