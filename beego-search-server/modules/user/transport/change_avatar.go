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

func ChangeAvatar(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		currentUser := c.Input.GetData(common.CurrentUser)
		if currentUser == "" || currentUser == nil {
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

		c.Output.Status = http.StatusOK
		c.Output.JSON(common.NewSuccessResponse(response), false, false)
	}
}
