package migrations

// Entry represents a registered migration in execution order.
type Entry struct {
	Name      string
	Migration interface{}
}

// All returns all Go migrations in chronological order.
// Each migration must be registered here to be tracked and executed by the unified migrator.
func All() []Entry {
	return []Entry{
		{
			Name:      "20260422150028_create_webhook_histories",
			Migration: &Migration_20260422150028{},
		},
		{
			Name:      "20260422154949_create_security_elite_tables",
			Migration: &Migration_20260422154949{},
		},
		{
			Name:      "20260422164224_create_failed_jobs_tables",
			Migration: &Migration_20260422164224{},
		},
		{
			Name:      "20260508000000_create_product_docs_table",
			Migration: &Migration_20260508000000{},
		},
		{
			Name:      "20260520000000_create_notifications_table",
			Migration: &Migration_20260520000000{},
		},
	}
}
