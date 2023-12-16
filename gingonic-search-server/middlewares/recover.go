package middlewares

import (
	"gingonic-search-server/common"
	"gingonic-search-server/component/appcontext"

	"github.com/gin-gonic/gin"
)

func Recover(appctx appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				appErr, ok := err.(*common.ErrorResponse)

				if !ok {
					appErr = common.ErrInternalServer(err.(error))
				}

				c.AbortWithStatusJSON(appErr.StatusCode, appErr)

				panic(err)
			}
		}()

		c.Next()
	}
}
