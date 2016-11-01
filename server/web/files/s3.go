package files

import (
	"io"
	"mime/multipart"

	"github.com/hooklift/lift-registry/server/config"
	"github.com/pkg/errors"
	"github.com/rlmcpherson/s3gof3r"
)

// S3 implements the storage driver for AWS S3.
type S3 struct {
	bucket *s3gof3r.Bucket
}

// NewS3 returns a new instance of an S3 storage provider.
func NewS3() StorageProvider {
	k, err := s3gof3r.EnvKeys() // get S3 keys from environment
	if err != nil {
		panic(err)
	}

	s3 := s3gof3r.New(s3gof3r.DefaultDomain, k)
	return &S3{bucket: s3.Bucket(config.S3Bucket)}
}

// Upload uploads file parts to S3 as they arrive from the client.
func (s *S3) Upload(reader *multipart.Reader) error {
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

		w, err := s.bucket.PutWriter(fileName, nil, s3gof3r.DefaultConfig)
		if err != nil {
			return errors.Wrapf(err, "failed initializing multipart request for %q", fileName)
		}

		if _, err = io.Copy(w, part); err != nil {
			return errors.Wrapf(err, "failed writing %q to S3", fileName)
		}

		if err = w.Close(); err != nil {
			return errors.Wrapf(err, "failed closing put writer for %q", fileName)
		}
	}

	return nil
}

// Get streams down a package file from S3.
// The caller must close the reader once it finishes reading from it.
func (s *S3) Get(key string) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting reader to download %q from S3", key)
	}

	return r, nil
}
