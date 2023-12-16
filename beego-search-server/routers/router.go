// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-search-server/component/appcontext"
	"beego-search-server/middlewares"
	searchtransport "beego-search-server/modules/search/transport"
	usertransport "beego-search-server/modules/user/transport"

	"github.com/beego/beego/v2/server/web"
)

func InitApiV1Router(appctx appcontext.AppContext) {
	// auth
	web.Post("/api/v1/register", usertransport.Register(appctx))
	web.Post("/api/v1/authenticate", usertransport.Authenticate(appctx))
	web.Post("/api/v1/refresh-token", usertransport.RefreshToken(appctx))

	// me
	web.InsertFilter("/api/v1/me/*", web.BeforeRouter, middlewares.AccessToken(appctx))
	web.Get("/api/v1/me", usertransport.Infor(appctx))
	web.Patch("/api/v1/me/modify", usertransport.Modify(appctx))
	web.Patch("/api/v1/me/change-avatar", usertransport.ChangeAvatar(appctx))

	// files
	web.InsertFilter("/api/v1/search/*", web.BeforeRouter, middlewares.AccessToken(appctx))
	web.InsertFilter("/api/v1/search/*", web.BeforeRouter, middlewares.Authorize(appctx))
	web.Post("/api/v1/search/files", searchtransport.ListFiles(appctx))
}
