package registry

// Repository is the interface to implement in order to retrieve data from a specific repository.
type Repository interface {
	Search(query string, pageNumber, resultsPerPage int) ([]*PluginManifest, error)
	Publish(p *PluginManifest) error
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

// Package is a plugin tarball prepared for a given CPU architecture and operating system.
type Package struct {
	Arch      Arch
	OS        OS
	URL       string
	Checksum  string
	Algorithm Algorithm
}

// PluginManifest is the document we use to index and return plugin manifest info.
type PluginManifest struct {
	// Name is the name given to the plugin
	Name string
	// Version is the version number given to the plugin
	Version string
	// Description is a short text describing what the plugin is for.
	Description string
	// Author is the company or individual who developed the plugin.
	Author string
	// License is the license under which the plugin was published.
	License string
	// Homepage can be the plugin repo or homepage.
	Homepage string
	// Packages is the list of packages this plugin has.
	Packages []*Package
}

// Search runs the specified query on the index file and returns a list of plugins.
func Search(query string, pageNumber, resultsPerPage int) ([]*PluginManifest, error) {
	return nil, nil
}

// Publish adds the plugin document into the index.
func Publish(p *PluginManifest) error {
	return nil
}
