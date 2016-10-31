package fileupload

import (
	"io"
	"net/http"

	"github.com/hooklift/uaa/openidc/pkg/render"
)

// StorageProvider defines the contract for storage providers.
type StorageProvider interface {
	Upload(reader io.Reader, name string) error
}

// Handler handles plugin packages uploads.
func Handler(h http.Handler, provider StorageProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/upload" {
			h.ServeHTTP(w, r)
			return
		}

		// Reads multipart body up until 32MB in memory, the rest gets stored temporarly in filesystem.
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		m := r.MultipartForm
		pkgs := m.File["package"]
		for i, header := range pkgs {
			file, err := pkgs[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			if err := provider.Upload(file, header.Filename); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		render.JSON(w, render.Options{
			Status: http.StatusOK,
		})
	})
}
