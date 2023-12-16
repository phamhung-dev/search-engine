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

func Infor(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := c.Get(common.CurrentUser)

		if !ok {
			panic(common.ErrInternalServer(models.ErrCurrentUserDoesNotExist))
		}

		user := currentUser.(*models.User)

		db := appctx.GetMainDBConnection()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewInforBusiness(storage)

		response, err := business.Infor(c.Request.Context(), user.ID)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
