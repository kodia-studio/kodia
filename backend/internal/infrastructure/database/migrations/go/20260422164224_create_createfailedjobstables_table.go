package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

type Migration_20260422164224 struct{}

func (m *Migration_20260422164224) Up(schema *database.Schema) error {
	return schema.Create("failed_jobs", func(table *database.Blueprint) {
		table.ID()
		table.String("uuid").Unique().NotNull()
		table.String("connection").NotNull()
		table.String("queue").NotNull()
		table.Text("payload").NotNull()
		table.Text("exception").NotNull()
		table.Timestamp("failed_at").NotNull()
		table.Timestamps()
	})
}

func (m *Migration_20260422164224) Down(schema *database.Schema) error {
	return schema.Drop("failed_jobs")
}
