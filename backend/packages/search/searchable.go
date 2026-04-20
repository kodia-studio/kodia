package search

/**
 * Searchable is an interface that models can implement to tell Kodia 
 * how to index them in the search engine.
 */
type Searchable interface {
	// SearchIndex returns the name of the index in Meilisearch/Algolia
	SearchIndex() string

	// SearchID returns the unique identifier for the document
	SearchID() string

	// ToSearchMap returns the map representation of the model for indexing
	ToSearchMap() map[string]interface{}
}
