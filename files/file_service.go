package files

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/hooklift/lift-registry/pkg/render"
	identity "github.com/hooklift/uaa/pkg/client"
)

// StorageProvider defines the contract for storage providers.
type StorageProvider interface {
	Upload(ctx context.Context, reader *multipart.Reader) error
	Get(ctx context.Context, filepath string) (io.ReadCloser, error)
}

// Response is the type of the payload sent back as response for uploading files.
type Response struct {
	URLs []string
}

// upload streams up file packages to S3 and returns their URLs once it finishes.
func upload(w http.ResponseWriter, r *http.Request) {
	token, ok := identity.FromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, oka := token.Scopes["admin"]
	_, okw := token.Scopes["write"]

	if !oka || !okw {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	if err := provider.Upload(ctx, reader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w)
}

// getPackage streams the requested file down to the user from the storage provider.
func getPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reader, err := provider.Get(ctx, path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := reader.Close(); err != nil {
			glog.Errorf("failed closing file reader: %+v", err)
		}
	}()

	if _, err := io.Copy(w, reader); err != nil {
		glog.Errorf("error streaming object from storage provider down to the user: %+v", err)
	}
}

var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"POST": upload,
	"GET":  getPackage,
}

var provider StorageProvider

// Handler handles /files requests.
func Handler(h http.Handler) http.Handler {
	registry := map[string]map[string]func(http.ResponseWriter, *http.Request){
		"/files": handlers,
	}

	provider = NewS3()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for p, handlers := range registry {
			if strings.HasPrefix(req.URL.Path, p) {
				if handlerFn, ok := handlers[req.Method]; ok {
					handlerFn(w, req)
					return
				}
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
		}
		h.ServeHTTP(w, req)
	})
}
