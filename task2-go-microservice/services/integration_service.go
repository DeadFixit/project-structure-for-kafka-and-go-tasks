package services

import (
	bytes "bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IntegrationService struct {
	client     *minio.Client
	bucketName string
}

func NewIntegrationService(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*IntegrationService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &IntegrationService{client: client, bucketName: bucket}, nil
}

func (s *IntegrationService) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
}

func (s *IntegrationService) UploadSampleObject(ctx context.Context, objectName string, content []byte) (minio.UploadInfo, error) {
	if err := s.EnsureBucket(ctx); err != nil {
		return minio.UploadInfo{}, err
	}
	reader := bytes.NewReader(content)
	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{ContentType: "text/plain"})
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}

func (s *IntegrationService) SampleContent(prefix string) []byte {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	return []byte(fmt.Sprintf("integration sample at %s (%s)", timestamp, prefix))
}
