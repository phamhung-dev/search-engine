package cachepvd

import (
	"beego-search-server/common"
	"context"
	"errors"
)

type CacheProvider interface {
	GetCacheData(ctx context.Context, key string) (string, error)
	SetCacheData(ctx context.Context, key string, data interface{}) error
}

var (
	ErrProviderIsNotConfigured = common.NewCustomErrorResponse(
		errors.New("cache provider is not configured"),
		"cache provider is not configured",
		"ErrProviderIsNotConfigured",
	)
)
