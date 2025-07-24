package adapters

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
)

type S3ClientAdapter struct {
	client *minio.Client
}

func NewS3ClientAdapter(client *minio.Client) *S3ClientAdapter {
	return &S3ClientAdapter{
		client: client,
	}
}

func (a *S3ClientAdapter) Upload(ctx context.Context, bucket string, key string, body []byte) error {
	_, err := a.client.PutObject(ctx, bucket, key, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
