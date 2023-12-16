package appcontext

import (
	"beego-search-server/component/authorizepvd"
	"beego-search-server/component/cachepvd"
	"beego-search-server/component/objectstoragepvd"
	"beego-search-server/component/tokenpvd"

	"github.com/beego/beego/v2/client/orm"
)

type AppContext interface {
	GetMainDBConnection() orm.Ormer
	GetTokenProvider() tokenpvd.TokenProvider
	GetAuthorizeProvider() authorizepvd.AuthorizeProvider
	GetCacheProvider() cachepvd.CacheProvider
	GetObjectStorageProvider() objectstoragepvd.ObjectStorageProvider
}

type appCtx struct {
	db                    orm.Ormer
	tokenProvider         tokenpvd.TokenProvider
	authorizeProvider     authorizepvd.AuthorizeProvider
	cacheProvider         cachepvd.CacheProvider
	objectStorageProvider objectstoragepvd.ObjectStorageProvider
}

func NewAppContext(
	db orm.Ormer,
	tokenProvider tokenpvd.TokenProvider,
	authorizeProvider authorizepvd.AuthorizeProvider,
	cacheProvider cachepvd.CacheProvider,
	objectStorageProvider objectstoragepvd.ObjectStorageProvider,
) *appCtx {
	return &appCtx{
		db:                    db,
		tokenProvider:         tokenProvider,
		authorizeProvider:     authorizeProvider,
		cacheProvider:         cacheProvider,
		objectStorageProvider: objectStorageProvider,
	}
}

func (ctx *appCtx) GetMainDBConnection() orm.Ormer {
	return ctx.db
}

func (ctx *appCtx) GetTokenProvider() tokenpvd.TokenProvider {
	return ctx.tokenProvider
}

func (ctx *appCtx) GetAuthorizeProvider() authorizepvd.AuthorizeProvider {
	return ctx.authorizeProvider
}

func (ctx *appCtx) GetCacheProvider() cachepvd.CacheProvider {
	return ctx.cacheProvider
}

func (ctx *appCtx) GetObjectStorageProvider() objectstoragepvd.ObjectStorageProvider {
	return ctx.objectStorageProvider
}
