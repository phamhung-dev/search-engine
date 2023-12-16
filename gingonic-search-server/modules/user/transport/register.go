package usertransport

import (
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/models"
	userbusiness "gingonic-search-server/modules/user/business"
	userstorage "gingonic-search-server/modules/user/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewRegisterBusiness(storage)

		response, err := business.Register(c.Request.Context(), &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
