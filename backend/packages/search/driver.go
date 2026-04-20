package search

/**
 * SearchDriver is the interface that all search drivers must implement.
 * This allows Kodia to support multiple search engines (Meilisearch, Algolia, etc.)
 */
type SearchDriver interface {
	// Search performs a search query on the given index
	Search(index string, query string, options map[string]interface{}) (interface{}, error)

	// Index adds or updates a document in the given index
	Index(index string, id string, data interface{}) error

	// Delete removes a document from the given index
	Delete(index string, id string) error

	// CreateIndex initializes a new index if it doesn't exist
	CreateIndex(index string, primaryKey string) error
}
