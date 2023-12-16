package objectstoragepvd

import (
	"beego-search-server/common"
	"context"
	"errors"
	"mime/multipart"
)

type ObjectStorageProvider interface {
	PutObject(ctx context.Context, bucketName string, file *multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

var (
	ErrProviderIsNotConfigured = common.NewCustomErrorResponse(
		errors.New("object storage provider is not configured"),
		"object storage provider is not configured",
		"ErrProviderIsNotConfigured",
	)
)
