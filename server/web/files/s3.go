package files

import (
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hooklift/lift-registry/server/config"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// S3 implements the storage driver for AWS S3.
type S3 struct {
	svc *s3.S3
}

// NewS3 returns a new instance of an S3 storage provider.
func NewS3() StorageProvider {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := s3.New(sess)
	return &S3{svc: svc}
}

// Upload uploads file parts to S3 as they arrive from the client.
func (s *S3) Upload(reader *multipart.Reader) error {
	id := uuid.NewV4().String()

	params := &s3.UploadPartInput{
		Bucket:   aws.String(config.S3Bucket),
		UploadId: aws.String(id),
	}

	var partNumber int64
	for {
		partNumber++
		params.PartNumber = aws.Int64(partNumber)

		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			// Ignore form fields that are not actual files
			continue
		}

		params.Key = aws.String(part.FileName())
		// parts do not implement ReadSeekCloser, so by doing the cast we loose request signing and retries
		// from the AWS SDK. Since we are expecting small files (< 100mb) this shouldn't be a concern for now.
		// The main benefit is simpler code and no buffering.
		params.Body = aws.ReadSeekCloser(part)

		_, err = s.svc.UploadPart(params)
		if err != nil {
			return errors.Wrapf(err, "failed uploading part %d to S3", partNumber)
		}
	}

	return nil
}

// Get streams down a package file from S3.
// The caller must close the reader once it finishes reading from it.
func (s *S3) Get(key string) (io.Reader, error) {
	result, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.S3Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed getting file %q from S3", key)
	}

	return result.Body, nil
}
