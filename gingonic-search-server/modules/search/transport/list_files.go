package searchtransport

import (
	"encoding/json"
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/models"
	searchbusiness "gingonic-search-server/modules/search/business"
	searchstorage "gingonic-search-server/modules/search/storage"
	"gingonic-search-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListFiles(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			filter models.FileFilter
			paging utils.Paging
			body   map[string]interface{}
		)

		if err := c.ShouldBind(&body); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		jsonData, err := json.Marshal(&body)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := json.Unmarshal(jsonData, &filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := json.Unmarshal(jsonData, &paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Fullfill()

		var response []models.File

		db := appctx.GetMainDBConnection()

		cacheProvider := appctx.GetCacheProvider()

		storage := searchstorage.NewStorage(db)

		business := searchbusiness.NewListFilesBusiness(storage)

		response, err = business.ListFiles(c.Request.Context(), cacheProvider, &filter, &paging)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewCustomSuccessResponse(response, filter, paging))
	}
}
