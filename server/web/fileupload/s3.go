package fileupload

import "io"

// S3 implements the storage driver for AWS S3.
type S3 struct{}

// Upload uploads ...
func (s *S3) Upload(reader io.Reader, name string) error {
	return nil
}
