package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260520000000 handles the creation of the notifications table.
type Migration_20260520000000 struct{}

func (m *Migration_20260520000000) Up(schema *database.Schema) error {
	return schema.Create("notifications", func(table *database.Blueprint) {
		table.ID()
		table.String("user_id").NotNull().Index()
		table.String("type").NotNull()
		table.String("title").NotNull()
		table.Text("message").NotNull()
		table.Text("data")
		table.Boolean("is_read").NotNull()
		table.Timestamps()
	})
}

func (m *Migration_20260520000000) Down(schema *database.Schema) error {
	return schema.Drop("notifications")
}
