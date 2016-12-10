// +build bleve

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/golang/glog"
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
func (r *RepoBleve) Search(ctx context.Context, query string, pageNumber, resultsPerPage int) ([]*Manifest, error) {
	matchQuery := bleve.NewQueryStringQuery(query)
	matchQuery.SetBoost(1)
	search := bleve.NewSearchRequest(matchQuery)
	search.Size = resultsPerPage
	search.From = pageNumber
	search.SortBy([]string{"-_score", "_id"})
	search.Fields = []string{"*"}

	if err := search.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid search query")
	}

	results, err := r.index.Search(search)
	if err != nil {
		return nil, errors.Wrapf(err, "failed searching %q", query)
	}

	manifests := make([]*Manifest, 0)
	for _, h := range results.Hits {
		fields := h.Fields

		manifest := &Manifest{
			ID:          fields["_id"].(string),
			AccountID:   fields["_account_id"].(string),
			Name:        fields["name"].(string),
			FilesURI:    fields["files_uri"].(string),
			Version:     fields["version"].(string),
			Description: fields["description"].(string),
			Author: Author{
				Name:  fields["author.name"].(string),
				Email: fields["author.email"].(string),
			},
			License:  fields["license"].(string),
			Homepage: fields["homepage"].(string),
		}

		publishedTime, err := time.Parse(time.RFC3339, fields["published_at"].(string))
		if err != nil {
			glog.Errorf("failed parsing published_at field coming from Bleve: %+v", err)
		} else {
			manifest.PublishedAt = publishedTime
		}

		packages := make([]*Package, 0)
		for i, name := range fields["packages.name"].([]interface{}) {
			p := &Package{Name: name.(string)}

			if v, ok := fields["packages.arch"]; ok {
				p.Arch = Arch(v.([]interface{})[i].(string))
			}

			if v, ok := fields["packages.os"]; ok {
				p.OS = OS(v.([]interface{})[i].(string))
			}

			if v, ok := fields["packages.checksum"]; ok {
				p.Checksum = v.([]interface{})[i].(string)
			}

			if v, ok := fields["packages.algorithm"]; ok {
				p.Algorithm = Algorithm(v.([]interface{})[i].(string))
			}

			packages = append(packages, p)
		}

		manifest.Packages = packages
		manifests = append(manifests, manifest)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return nil, errors.Wrap(err, "failed encoding search result")
	}
	return manifests, nil
}

// Save indexes plugin metadata in Bleve's index.
func (r *RepoBleve) Save(ctx context.Context, p *Manifest) error {
	if p == nil {
		return errors.New("manifest is required")
	}

	if p.ID == "" {
		return errors.New("ID is required")
	}

	return r.index.Index(p.ID, p)
}

// Delete removes plugin from Bleve index.
func (r *RepoBleve) Delete(ctx context.Context, id, accountID string) error {
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

	if doc == nil {
		return fmt.Errorf("plugin %q not found", id)
	}

	for _, f := range doc.Fields {
		if f.Name() == "_account_id" && bytes.Equal(f.Value(), []byte(accountID)) {
			return r.index.Delete(id)
		}
	}

	return nil
}
