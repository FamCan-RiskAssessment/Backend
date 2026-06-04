package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	constants *bootstrap.Constants
	storage   *bootstrap.MinIO
	client    *minio.Client
	prefixes  map[enum.BucketType]string
}

func NewMinIOStorage(
	constants *bootstrap.Constants,
	storage *bootstrap.MinIO,
) *MinIOStorage {
	prefixes := make(map[enum.BucketType]string)
	prefixes[enum.BucketTypeMamography] = storage.Prefixes.Mamography
	prefixes[enum.BucketTypeCancer] = storage.Prefixes.Cancer
	prefixes[enum.BucketTypeGeneticTest] = storage.Prefixes.GeneticTest
	prefixes[enum.BucketTypeFatherGeneticTest] = storage.Prefixes.FatherGeneticTest
	prefixes[enum.BucketTypeMotherGeneticTest] = storage.Prefixes.MotherGeneticTest
	return &MinIOStorage{
		constants: constants,
		storage:   storage,
		prefixes:  prefixes,
	}
}

func parseMinIOEndpoint(endpoint string) (host string, secure bool, err error) {
	if endpoint == "" {
		return "", false, fmt.Errorf("minio endpoint is empty")
	}
	if strings.HasPrefix(endpoint, "https://") {
		return strings.TrimPrefix(endpoint, "https://"), true, nil
	}
	if strings.HasPrefix(endpoint, "http://") {
		return strings.TrimPrefix(endpoint, "http://"), false, nil
	}
	return endpoint, false, nil
}

func (m *MinIOStorage) setClient(bucketType enum.BucketType) error {
	bucketTypes := enum.GetAllBucketTypes()
	if !slices.Contains(bucketTypes, bucketType) {
		return fmt.Errorf("bucket not exist")
	}
	if m.client != nil {
		return nil
	}

	endpoint, secure, err := parseMinIOEndpoint(m.storage.Endpoint)
	if err != nil {
		return err
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.storage.AccessKey, m.storage.SecretKey, ""),
		Secure: secure,
		Region: m.storage.Region,
	})
	if err != nil {
		return fmt.Errorf("unable to create MinIO client, %w", err)
	}

	m.client = client
	return nil
}

func (m *MinIOStorage) objectKey(bucketType enum.BucketType, key string) string {
	prefix := strings.Trim(m.prefixes[bucketType], "/")
	key = strings.TrimPrefix(key, "/")
	if prefix == "" {
		return key
	}
	return prefix + "/" + key
}

func (m *MinIOStorage) ensureBucket(ctx context.Context) error {
	bucket := m.storage.Bucket

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("unable to check bucket %q, %w", bucket, err)
	}
	if exists {
		return nil
	}

	if err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: m.storage.Region}); err != nil {
		return fmt.Errorf("unable to create bucket %q, %w", bucket, err)
	}
	return nil
}

func (m *MinIOStorage) UploadObject(bucketType enum.BucketType, key string, file *multipart.FileHeader) error {
	if err := m.setClient(bucketType); err != nil {
		return err
	}

	bucket := m.storage.Bucket
	objectKey := m.objectKey(bucketType, key)
	ctx := context.Background()

	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file %q, %w", file.Filename, err)
	}
	defer fileReader.Close()

	if err := m.ensureBucket(ctx); err != nil {
		return err
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = m.client.PutObject(ctx, bucket, objectKey, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %w", file.Filename, bucket, err)
	}
	return nil
}

func (m *MinIOStorage) DeleteObject(bucketType enum.BucketType, key string) error {
	if err := m.setClient(bucketType); err != nil {
		return err
	}

	bucket := m.storage.Bucket
	objectKey := m.objectKey(bucketType, key)

	err := m.client.RemoveObject(context.Background(), bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("unable to delete %q from %q, %w", objectKey, bucket, err)
	}
	return nil
}

func (m *MinIOStorage) GetPresignedURL(bucketType enum.BucketType, objectKey string, expiration time.Duration) (string, error) {
	if err := m.setClient(bucketType); err != nil {
		return "", err
	}

	bucket := m.storage.Bucket
	key := m.objectKey(bucketType, objectKey)

	reqParams := make(url.Values)
	reqParams.Set("response-content-type", "image/jpeg")
	reqParams.Set("response-content-disposition", "inline")

	presignedURL, err := m.client.PresignedGetObject(context.Background(), bucket, key, expiration, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}
