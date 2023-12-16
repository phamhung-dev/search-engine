package searchbusiness

import (
	"beego-search-server/common"
	"beego-search-server/component/cachepvd"
	"beego-search-server/models"
	"beego-search-server/utils"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type ListFilesStorage interface {
	List(context context.Context, filter *models.FileFilter, paging *utils.Paging, relmd ...string) ([]models.File, error)
}

type listFilesBusiness struct {
	storage ListFilesStorage
}

func NewListFilesBusiness(storage ListFilesStorage) *listFilesBusiness {
	return &listFilesBusiness{storage: storage}
}

func (business *listFilesBusiness) ListFiles(context context.Context, cacheProvider cachepvd.CacheProvider, filter *models.FileFilter, paging *utils.Paging) ([]models.File, error) {
	if err := filter.Validate(); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	files := []models.File{}

	cacheKey := utils.GenerateCacheKey(models.FileTableName, filter, paging)

	cachedFiles, err := cacheProvider.GetCacheData(context, cacheKey)
	if (err == nil || err == redis.Nil) && cachedFiles != "" && json.Unmarshal([]byte(cachedFiles), &files) == nil {
		return files, nil
	}

	files, err = business.storage.List(context, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(models.FileEntityName, err)
	}

	for i := range files {
		files[i].Mask(false)
	}

	cacheProvider.SetCacheData(context, cacheKey, files)

	return files, nil
}
