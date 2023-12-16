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

func ChangeAvatar(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := c.Get(common.CurrentUser)

		if !ok {
			panic(common.ErrInternalServer(models.ErrCurrentUserDoesNotExist))
		}

		user := currentUser.(*models.User)

		file, header, err := c.Request.FormFile("avatar")
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		data := models.UserAvatar{
			File:   &file,
			Header: header,
		}

		db := appctx.GetMainDBConnection()

		objectStorageProvider := appctx.GetObjectStorageProvider()

		storage := userstorage.NewStorage(db)

		business := userbusiness.NewChangeAvatarBusiness(storage)

		response, err := business.ChangeAvatar(c.Request.Context(), objectStorageProvider, user.ID, &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(response))
	}
}
