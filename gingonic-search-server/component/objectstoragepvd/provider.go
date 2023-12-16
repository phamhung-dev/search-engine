package objectstoragepvd

import (
	"context"
	"errors"
	"gingonic-search-server/common"
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
