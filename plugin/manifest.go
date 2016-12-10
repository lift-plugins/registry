package plugin

import (
	"context"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Repo should be initialized by a concrete repository implementation.
var Repo Repository

// Repository is the interface to implement in order to retrieve data from a specific repository.
type Repository interface {
	Search(ctx context.Context, query string, pageNumber, resultsPerPage int) ([]*Manifest, error)
	Save(ctx context.Context, p *Manifest) error
	Delete(ctx context.Context, id, accountID string) error
}

// Arch is the CPU architecture for which a plugin package was compiled.
type Arch string

const (
	x86   Arch = "x86"
	x64   Arch = "x64"
	arm   Arch = "arm"
	arm64 Arch = "arm64"
)

// OS is the operating system supported for which a plugin's package was compiled.
type OS string

const (
	windows OS = "windows"
	macOS   OS = "macOS"
	freebsd OS = "freebsd"
	linux   OS = "linux"
)

// Algorithm is the hashing algorithm used to calculate packages checksums.
type Algorithm string

const (
	sha256 Algorithm = "sha256"
	sha512 Algorithm = "sha512"
)

// Author is the plugin author.
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Package is a plugin tarball prepared for a given CPU architecture and operating system.
type Package struct {
	Name      string    `json:"name"`
	Arch      Arch      `json:"arch"`
	OS        OS        `json:"os"`
	Checksum  string    `json:"checksum"`
	Algorithm Algorithm `json:"algorithm"`
}

// Manifest is the document we use to index and return plugin manifest info.
type Manifest struct {
	// Internal document ID
	ID string `json:"_id"`
	// Account publishing the plugin.
	AccountID string `json:"_account_id"`
	// Name is the name given to the plugin
	Name string `json:"name"`
	// FilesURI is the base URL used to download plugin packages
	FilesURI string `json:"files_uri"`
	// Version is the version number given to the plugin
	Version string `json:"version"`
	// Description is a short text describing what the plugin is for.
	Description string `json:"description"`
	// Author is the company or individual who developed the plugin.
	Author Author `json:"author"`
	// License is the license under which the plugin was published.
	License string `json:"license"`
	// Homepage can be the plugin repo or homepage.
	Homepage string `json:"homepage"`
	// Packages is the list of packages this plugin has.
	Packages []*Package `json:"packages"`
	// PublishedAt is the time when this plugin was published.
	PublishedAt time.Time `json:"published_at"`
}

// Search runs the specified query on the index file and returns a list of plugins.
func Search(ctx context.Context, query string, pageNumber, resultsPerPage int) ([]*Manifest, error) {
	if resultsPerPage == 0 {
		resultsPerPage = 10
	}

	if resultsPerPage > 50 {
		resultsPerPage = 50
	}

	return Repo.Search(ctx, query, pageNumber, resultsPerPage)
}

// Publish adds the plugin document into the index.
func Publish(ctx context.Context, p *Manifest) error {
	if p == nil {
		return errors.New("a valid manifest is required")
	}

	if len(p.Packages) == 0 {
		return errors.New("list of packages missing")
	}

	ver, err := version.NewVersion(p.Version)
	if err != nil {
		return errors.Wrap(err, "invalid plugin version")
	}

	p.Version = ver.String()
	p.ID = uuid.NewV4().String()
	p.PublishedAt = time.Now()

	return Repo.Save(ctx, p)
}

// Unpublish removes a plugin from the index.
func Unpublish(ctx context.Context, id, accountID string) error {
	if id == "" {
		return errors.New("document ID is required")
	}

	if accountID == "" {
		return errors.New("account ID is required")
	}

	return Repo.Delete(ctx, id, accountID)
}
