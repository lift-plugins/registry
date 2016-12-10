package files

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/hooklift/lift-registry/config"
	"github.com/pkg/errors"
)

// S3 implements the storage driver for AWS S3.
type S3 struct {
	uploader   *s3manager.Uploader
	downloader *s3.S3
}

// NewS3 returns a new instance of an S3 storage provider.
func NewS3() StorageProvider {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(sess)
	downloader := s3.New(sess)

	return &S3{
		uploader:   uploader,
		downloader: downloader,
	}
}

// Upload uploads file parts to S3 as they arrive from the client.
func (s *S3) Upload(ctx context.Context, reader *multipart.Reader) error {
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fileName := part.FileName()
		if fileName == "" {
			// Ignore form fields that are not actual files
			continue
		}

		input := &s3manager.UploadInput{
			Bucket: aws.String(config.S3Bucket),
			Key:    aws.String(fileName),
			Body:   part,
		}

		_, err = s.uploader.Upload(input)
		if err != nil {
			return errors.Wrapf(err, "failed uploading %q to S3", fileName)
		}
	}

	return nil
}

// Get streams down a package file from S3.
// The caller must close the reader once it finishes reading from it.
func (s *S3) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.S3Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed downloading %q from S3", key)
	}

	return result.Body, nil
}
