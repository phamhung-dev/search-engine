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

func Register(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		var data models.UserCreate

		if err := c.Bind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewRegisterBusiness(storage)

		response, err := business.Register(c.Request.Context(), &data)

		if err != nil {
			panic(err)
		}

		c.Output.Status = http.StatusOK
		c.Output.JSON(common.NewSuccessResponse(response), false, false)
	}
}
