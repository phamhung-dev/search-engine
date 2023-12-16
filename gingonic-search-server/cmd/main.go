package main

import (
	"gingonic-search-server/component/appcontext"
	"gingonic-search-server/component/authorizepvd"
	"gingonic-search-server/component/cachepvd"
	"gingonic-search-server/component/objectstoragepvd"
	"gingonic-search-server/component/tokenpvd"
	"gingonic-search-server/postgresql"
	"gingonic-search-server/routers"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	appctx := appcontext.NewAppContext(
		postgresql.NewPostgreSQLConnection(),
		tokenpvd.NewJWTProvider(),
		authorizepvd.NewCasbinProvider(),
		cachepvd.NewRedisProvider(),
		objectstoragepvd.NewMinioProvider(),
	)

	server := routers.InitRouter(appctx)

	server.Run()
}
