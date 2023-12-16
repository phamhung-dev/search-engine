package objectstoragepvd

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioObjectStorageProvider struct {
	client *minio.Client
}

func NewMinioProvider() *minioObjectStorageProvider {
	host := os.Getenv("MINIO_HOST")
	port := os.Getenv("MINIO_PORT")
	user := os.Getenv("MINIO_ROOT_USER")
	password := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL, err := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))

	if host == "" || port == "" || user == "" || password == "" || err != nil {
		log.Fatal(ErrProviderIsNotConfigured.Error())
	}

	client, err := minio.New(
		fmt.Sprintf("%s:%s", host, port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(user, password, ""),
			Secure: useSSL,
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &minioObjectStorageProvider{client: client}
}

func (provider *minioObjectStorageProvider) PutObject(ctx context.Context, bucketName string, file *multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = "us-east-1"
	}

	if exists, _ := provider.client.BucketExists(ctx, bucketName); !exists {
		if err := provider.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region}); err != nil {
			return "", err
		}

		policy := `{ "Version": "2012-10-17", "Statement": [ { "Effect": "Allow", "Principal": { "AWS": [ "*" ] }, "Action": [ "s3:GetObject" ], "Resource": [ "arn:aws:s3:::avatars/*" ] } ] }`
		if err := provider.client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
			return "", err
		}
	}

	if _, err := provider.client.PutObject(ctx, bucketName, fileHeader.Filename, *file, fileHeader.Size, minio.PutObjectOptions{ContentType: fileHeader.Header["Content-Type"][0]}); err != nil {
		return "", err
	}

	path := fmt.Sprintf("/%s/%s", bucketName, fileHeader.Filename)
	return path, nil
}
