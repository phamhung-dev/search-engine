package appcontext

import (
	"gingonic-search-server/component/authorizepvd"
	"gingonic-search-server/component/cachepvd"
	"gingonic-search-server/component/objectstoragepvd"
	"gingonic-search-server/component/tokenpvd"

	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetTokenProvider() tokenpvd.TokenProvider
	GetAuthorizeProvider() authorizepvd.AuthorizeProvider
	GetCacheProvider() cachepvd.CacheProvider
	GetObjectStorageProvider() objectstoragepvd.ObjectStorageProvider
}

type appCtx struct {
	db                    *gorm.DB
	tokenProvider         tokenpvd.TokenProvider
	authorizeProvider     authorizepvd.AuthorizeProvider
	cacheProvider         cachepvd.CacheProvider
	objectStorageProvider objectstoragepvd.ObjectStorageProvider
}

func NewAppContext(
	db *gorm.DB,
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

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
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
