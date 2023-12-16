package usertransport

import (
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/models"
	userbusiness "gingonic-search-server/modules/user/business"
	userstorage "gingonic-search-server/modules/user/storage"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Authenticate(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.UserLogin

		if err := c.ShouldBind(&data); err != nil {
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
			c.SetCookie("refresh_token", response.RefreshToken, response.ExpiredIn, "/", domain, false, true)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
