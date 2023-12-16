package usertransport

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	"beego-search-server/models"
	userbusiness "beego-search-server/modules/user/business"
	userstorage "beego-search-server/modules/user/storage"
	"net/http"
	"os"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func Authenticate(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		var data models.UserLogin

		if err := c.Bind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		tokenProvider := appctx.GetTokenProvider()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewAuthenticateBusiness(storage, tokenProvider)

		response, err := business.Authenticate(c.Request.Context(), &data)

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
