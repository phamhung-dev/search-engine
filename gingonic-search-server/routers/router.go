package routers

import (
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/middlewares"
	searchtransport "gingonic-search-server/modules/search/transport"
	usertransport "gingonic-search-server/modules/user/transport"

	"github.com/gin-gonic/gin"
)

func InitRouter(appctx appcontext.AppContext) *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.Recover(appctx))

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/register", usertransport.Register(appctx))
			v1.POST("/authenticate", usertransport.Authenticate(appctx))
			v1.POST("/refresh-token", usertransport.RefreshToken(appctx))

			me := v1.Group("/me")
			{
				me.Use(middlewares.AccessToken(appctx))
				me.GET("", usertransport.Infor(appctx))
				me.PATCH("/modify", usertransport.Modify(appctx))
				me.PATCH("/change-avatar", usertransport.ChangeAvatar(appctx))
			}

			search := v1.Group("/search")
			{
				search.Use(middlewares.AccessToken(appctx), middlewares.Authorize(appctx))
				search.POST("/files", searchtransport.ListFiles(appctx))
			}
		}
	}

	return router
}
