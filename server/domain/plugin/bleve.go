// +build bleve

package plugin

import (
	"fmt"

	"github.com/blevesearch/bleve"
)

// RepoBleve represents an implementation of the Repo interface for Bleve search engine.
type RepoBleve struct {
	index bleve.Index
}

// NewRepository creates an instance of the Bleve repository.
func NewRepository(index bleve.Index) Repository {
	return &RepoBleve{
		index: index,
	}
}

// Search finds plugin manifests in Bleve.
func (r *RepoBleve) Search(query string, pageNumber, resultsPerPage int) ([]*Manifest, error) {
	matchQuery := bleve.NewMatchQuery(query)
	search := bleve.NewSearchRequest(matchQuery)
	search.Size = resultsPerPage
	search.From = pageNumber

	results, err := r.index.Search(search)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Results => %s\n", results)
	return nil, nil
}

// Save indexes plugin metadata in Bleve's index.
func (r *RepoBleve) Save(p *Manifest) error {
	// We use the name of the plugin as the ID to make sure there is always one document
	// in the index for the same plugin.
	return r.index.Index(p.Name, p)
}

// Delete removes plugin from Bleve index.
func (r *RepoBleve) Delete(id string) error {
	return r.index.Delete(id)
}
