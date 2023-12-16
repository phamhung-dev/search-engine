package searchtransport

import (
	"beego-search-server/common"
	"beego-search-server/component/appcontext"
	"beego-search-server/models"
	searchbusiness "beego-search-server/modules/search/business"
	searchstorage "beego-search-server/modules/search/storage"
	"beego-search-server/utils"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func ListFiles(appctx appcontext.AppContext) web.HandleFunc {
	return func(c *beecontext.Context) {
		var (
			filter models.FileFilter
			paging utils.Paging
		)

		if err := c.Bind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := c.Bind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Fullfill()

		var response []models.File

		db := appctx.GetMainDBConnection()

		cacheProvider := appctx.GetCacheProvider()

		storage := searchstorage.NewStorage(db)

		business := searchbusiness.NewListFilesBusiness(storage)

		response, err := business.ListFiles(c.Request.Context(), cacheProvider, &filter, &paging)

		if err != nil {
			panic(err)
		}

		c.Output.Status = http.StatusOK
		c.Output.JSON(common.NewCustomSuccessResponse(response, filter, paging), false, false)
	}
}
