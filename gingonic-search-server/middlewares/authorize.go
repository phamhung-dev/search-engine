package middlewares

import (
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/models"

	"github.com/gin-gonic/gin"
)

func Authorize(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := c.Get(common.CurrentUser)

		if !ok {
			panic(common.ErrInternalServer(models.ErrCurrentUserDoesNotExist))
		}

		user := currentUser.(*models.User)
		path := c.Request.URL.Path
		method := c.Request.Method

		authorizeProvider := appctx.GetAuthorizeProvider()

		result, err := authorizeProvider.ValidateRequest(user, path, method)
		if !result || err != nil {
			panic(common.ErrPermissionDenied(common.ErrUserDoesNotHaveAccess))
		}

		c.Next()
	}
}
