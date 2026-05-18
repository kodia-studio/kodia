package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260508000000 handles the creation of the product_docs table.
type Migration_20260508000000 struct{}

func (m *Migration_20260508000000) Up(schema *database.Schema) error {
	return schema.Create("product_docs", func(table *database.Blueprint) {
		table.ID()
		table.String("product_id").NotNull()
		table.String("doc_type").NotNull()
		table.String("section_title").NotNull()
		table.String("section_slug").NotNull()
		table.String("title").NotNull()
		table.String("slug").NotNull()
		table.Text("content").NotNull()
		table.Integer("sort_order")
		table.Integer("section_order")
		table.Boolean("is_published")
		table.Timestamp("published_at")
		table.String("meta_title")
		table.String("meta_description")
		table.Timestamps()
	})
}

func (m *Migration_20260508000000) Down(schema *database.Schema) error {
	return schema.Drop("product_docs")
}
