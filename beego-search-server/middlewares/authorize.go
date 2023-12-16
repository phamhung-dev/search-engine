package middlewares

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	"beego-search-server/models"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func Authorize(appctx appcontext.AppContext) web.FilterFunc {
	return func(c *beecontext.Context) {
		currentUser := c.Input.GetData(common.CurrentUser)

		if currentUser == "" || currentUser == nil {
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
	}
}
