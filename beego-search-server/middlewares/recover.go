package middlewares

import (
	"beego-search-server/common"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func InitRecover() {
	originRecover := web.BConfig.RecoverFunc
	web.BConfig.RecoverFunc = func(c *beecontext.Context, config *web.Config) {
		defer originRecover(c, config)

		if err := recover(); err != nil {
			c.Output.Header("Content-Type", "application/json")

			appErr, ok := err.(*common.ErrorResponse)

			if !ok {
				appErr = common.ErrInternalServer(err.(error))
			}

			c.Output.Status = appErr.StatusCode
			c.Output.JSON(appErr, false, false)

			panic(err)
		}
	}
}
