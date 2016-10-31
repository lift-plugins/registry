package files

import (
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/hooklift/lift-registry/server/pkg/render"
)

// StorageProvider defines the contract for storage providers.
type StorageProvider interface {
	Upload(reader *multipart.Reader) error
	Get(filepath string) (io.Reader, error)
}

// Response is the type of the payload sent back as response for uploading files.
type Response struct {
	URLs []string
}

// upload streams up file packages to S3 and returns their URLs once it finishes.
func upload(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := provider.Upload(reader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w)
}

// getPackage streams the requested file down to the user from the storage provider.
func getPackage(w http.ResponseWriter, r *http.Request) {
	reader, err := provider.Get(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
