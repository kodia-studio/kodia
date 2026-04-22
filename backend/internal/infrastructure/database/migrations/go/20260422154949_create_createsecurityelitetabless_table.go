package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

type Migration_20260422154949 struct{}

func (m *Migration_20260422154949) Up(schema *database.Schema) error {
	// 1. API Keys table
	err := schema.Create("api_keys", func(table *database.Blueprint) {
		table.ID()
		table.String("user_id").NotNull().Index()
		table.String("name").NotNull()
		table.String("key").Unique().NotNull()
		table.Text("scopes") // JSON string field
		table.Timestamp("last_used_at")
		table.Timestamp("expires_at")
		table.Timestamps()
	})
	if err != nil {
		return err
	}

	// 2. WebAuthn Credentials (Passkeys)
	err = schema.Create("webauthn_credentials", func(table *database.Blueprint) {
		table.ID() // Binary blob ID usually stored as string or bytea
		table.String("user_id").NotNull().Index()
		table.Binary("public_key").NotNull()
		table.String("attestation_type")
		table.Text("transports")
		table.Integer("sign_count")
		table.Timestamps()
	})
	if err != nil {
		return err
	}

	// 3. Sessions
	return schema.Create("sessions", func(table *database.Blueprint) {
		table.ID()
		table.String("user_id").NotNull().Index()
		table.String("user_agent")
		table.String("ip_address")
		table.Timestamp("expires_at").NotNull()
		table.Timestamps()
	})
}

func (m *Migration_20260422154949) Down(schema *database.Schema) error {
	schema.Drop("sessions")
	schema.Drop("webauthn_credentials")
	return schema.Drop("api_keys")
}
