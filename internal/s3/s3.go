package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	conf "github.com/refine-software/afrad-api/config"
)

type S3 interface {
	UploadImage(
		ctx *gin.Context,
		file multipart.File,
		fileHeader *multipart.FileHeader,
	) (string, error)

	DeleteImageByURL(ctx *gin.Context, fileURL string) error
}

type S3Storage struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewS3Storage initializes a new S3 client and wraps it
func NewS3Storage(env *conf.Env) (*S3Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(env.S3Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(env.S3AccessKey, env.S3SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Storage{
		client:     client,
		bucketName: env.S3Bucket,
		region:     env.S3Region,
	}, nil
}

// UploadImage uploads a multipart file to the S3 bucket
func (s *S3Storage) UploadImage(
	ctx *gin.Context,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(file)
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(fileHeader.Filename)
	objectKey := fmt.Sprintf("images/%d%s", time.Now().UnixNano(), fileExt)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, objectKey)
	return url, nil
}

// DeleteImageByURL removes a file from the S3 bucket using the full S3 URL
func (s *S3Storage) DeleteImageByURL(ctx *gin.Context, fileURL string) error {
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucketName, s.region)
	if !strings.HasPrefix(fileURL, prefix) {
		return fmt.Errorf("invalid S3 URL: %s", fileURL)
	}
	objectKey := strings.TrimPrefix(fileURL, prefix)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	waiter := s3.NewObjectNotExistsWaiter(s.client)
	err = waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	}, time.Minute)
	if err != nil {
		return fmt.Errorf("object still exists after deletion: %w", err)
	}

	return nil
}
