package usertransport

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	"beego-search-server/models"
	userbusiness "beego-search-server/modules/user/business"
	userstorage "beego-search-server/modules/user/storage"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func Modify(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		currentUser := c.Input.GetData(common.CurrentUser)
		if currentUser == "" || currentUser == nil {
			panic(common.ErrInternalServer(models.ErrCurrentUserDoesNotExist))
		}

		user := currentUser.(*models.User)

		var data models.UserUpdate

		if err := c.Bind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewModifyBusiness(storage)

		response, err := business.Modify(c.Request.Context(), user.ID, &data)

		if err != nil {
			panic(err)
		}

		c.Output.Status = http.StatusOK
		c.Output.JSON(common.NewSuccessResponse(response), false, false)
	}
}
