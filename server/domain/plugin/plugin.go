package plugin

import (
	"errors"
	"time"
)

// Repo should be initialized by a concrete repository implementation.
var Repo Repository

// Repository is the interface to implement in order to retrieve data from a specific repository.
type Repository interface {
	Search(query string, pageNumber, resultsPerPage int) ([]*Manifest, error)
	Save(p *Manifest) error
	Delete(id string) error
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

// Version is the plugin version being published.
type Version struct {
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
	Patch uint32 `json:"patch"`
}

// Package is a plugin tarball prepared for a given CPU architecture and operating system.
type Package struct {
	Arch      Arch      `json:"arch"`
	OS        OS        `json:"os"`
	URL       string    `json:"url"`
	Checksum  string    `json:"checksum"`
	Algorithm Algorithm `json:"algorithm"`
}

// Manifest is the document we use to index and return plugin manifest info.
type Manifest struct {
	// Name is the name given to the plugin
	Name string `json:"name"`
	// Version is the version number given to the plugin
	Version Version `json:"version"`
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
func Search(query string, pageNumber, resultsPerPage int) ([]*Manifest, error) {
	return Repo.Search(query, pageNumber, resultsPerPage)
}

// Publish adds the plugin document into the index.
func Publish(p *Manifest) error {
	if p == nil {
		return errors.New("a valid manifest is required")
	}

	p.PublishedAt = time.Now()

	return Repo.Save(p)
}

// Unpublish removes a plugin from the index.
func Unpublish(id string) error {
	return Repo.Delete(id)
	// TODO(c4milo): Remove packages from files.
}
