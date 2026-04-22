package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260422150028 handles the creation of the webhook_histories table.
type Migration_20260422150028 struct{}

func (m *Migration_20260422150028) Up(schema *database.Schema) error {
	return schema.Create("webhook_histories", func(table *database.Blueprint) {
		table.ID()
		table.String("url").NotNull()
		table.String("event").NotNull()
		table.Text("payload")
		table.Integer("response_code")
		table.Boolean("success").NotNull()
		table.Text("error_message")
		table.Timestamps()
	})
}

func (m *Migration_20260422150028) Down(schema *database.Schema) error {
	return schema.Drop("webhook_histories")
}
