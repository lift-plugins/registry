package domain

import (
	"github.com/blevesearch/bleve"
	"github.com/hooklift/lift-registry/server/domain/plugin"
)

// Init initializes all the domain modules.
func Init(index bleve.Index) {
	// The repository layer compiled is determined by build flags
	plugin.Repo = plugin.NewRepository(index)
}
