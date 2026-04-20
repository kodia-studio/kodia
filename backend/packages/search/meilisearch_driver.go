package search

import (
	meilisearch "github.com/meilisearch/meilisearch-go"
)

type MeiliSearchDriver struct {
	client meilisearch.ServiceManager
}

func NewMeiliSearchDriver(host, apiKey string) *MeiliSearchDriver {
	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	return &MeiliSearchDriver{client: client}
}

func (d *MeiliSearchDriver) Search(index string, query string, options map[string]interface{}) (interface{}, error) {
	searchRes, err := d.client.Index(index).Search(query, &meilisearch.SearchRequest{
		// Map options if needed, here just basic search
		Limit: 20,
	})
	if err != nil {
		return nil, err
	}
	return searchRes, nil
}

func (d *MeiliSearchDriver) Index(index string, id string, data interface{}) error {
	// Meilisearch handles slice of documents
	documents := []interface{}{data}
	_, err := d.client.Index(index).AddDocuments(documents, nil)
	return err
}

func (d *MeiliSearchDriver) Delete(index string, id string) error {
	_, err := d.client.Index(index).DeleteDocument(id, nil)
	return err
}

func (d *MeiliSearchDriver) CreateIndex(index string, primaryKey string) error {
	_, err := d.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        index,
		PrimaryKey: primaryKey,
	})
	return err
}
