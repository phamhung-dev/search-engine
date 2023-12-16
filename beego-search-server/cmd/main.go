package main

import (
	"beego-search-server/component/appcontext"
	"beego-search-server/component/authorizepvd"
	"beego-search-server/component/cachepvd"
	"beego-search-server/component/objectstoragepvd"
	"beego-search-server/component/tokenpvd"
	"beego-search-server/middlewares"
	"beego-search-server/postgresql"
	"beego-search-server/routers"

	"github.com/beego/beego/v2/server/web"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	middlewares.InitRecover()

	if web.BConfig.RunMode == web.DEV {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	appctx := appcontext.NewAppContext(
		postgresql.NewPostgreSQLConnection(),
		tokenpvd.NewJWTProvider(),
		authorizepvd.NewCasbinProvider(),
		cachepvd.NewRedisProvider(),
		objectstoragepvd.NewMinioProvider(),
	)

	routers.InitApiV1Router(appctx)

	web.Run()
}
