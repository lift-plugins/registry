// +build bleve

package plugin

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "failed searching %q", query)
	}

	fmt.Printf("Results => %s\n", results)
	return nil, nil
}

// Save indexes plugin metadata in Bleve's index.
func (r *RepoBleve) Save(p *Manifest) error {
	if p == nil {
		return errors.New("manifest is required")
	}

	if p.ID == "" {
		return errors.New("ID is required")
	}

	return r.index.Index(p.ID, p)
}

// Delete removes plugin from Bleve index.
func (r *RepoBleve) Delete(id, accountID string) error {
	if id == "" {
		return errors.New("document ID is required")
	}

	if accountID == "" {
		return errors.New("account ID is required")
	}

	doc, err := r.index.Document(id)
	if err != nil {
		return errors.Wrapf(err, "document ID %q does not exist", id)
	}

	for _, f := range doc.Fields {
		if f.Name() == accountID {
			return r.index.Delete(id)
		}
	}

	return nil
}
