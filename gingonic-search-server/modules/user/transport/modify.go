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

func Modify(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := c.Get(common.CurrentUser)

		if !ok {
			panic(common.ErrInternalServer(models.ErrCurrentUserDoesNotExist))
		}

		user := currentUser.(*models.User)

		var data models.UserUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appctx.GetMainDBConnection()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewModifyBusiness(storage)

		response, err := business.Modify(c.Request.Context(), user.ID, &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
