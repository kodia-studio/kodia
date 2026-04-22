package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Schema provides the entry point for database schema manipulation.
type Schema struct {
	db *gorm.DB
}

// NewSchema creates a new Schema builder.
func NewSchema(db *gorm.DB) *Schema {
	return &Schema{db: db}
}

// Blueprint defines the structure of a database table.
type Blueprint struct {
	Name    string
	Columns []Column
}

// Column represents a single database column.
type Column struct {
	Name       string
	Type       string
	Nullable   bool
	IsUnique   bool
	PrimaryKey bool
	Default    string
	References string
}

// Create creates a new table on the schema.
func (s *Schema) Create(name string, callback func(t *Blueprint)) error {
	b := &Blueprint{Name: name}
	callback(b)

	// In a real implementation, this would generate SQL or use GORM's Migrator.
	// For Kodia, we'll use GORM's Migrator but provide a more fluent API surface.
	migrator := s.db.Migrator()
	
	// Check if table exists
	if migrator.HasTable(name) {
		return fmt.Errorf("table %s already exists", name)
	}

	// We use raw SQL for precise control in the Blueprint, then execute.
	sql := s.generateCreateTableSQL(b)
	return s.db.Exec(sql).Error
}

// Drop drops a table from the schema.
func (s *Schema) Drop(name string) error {
	return s.db.Migrator().DropTable(name)
}

// -- Blueprint Helpers --

func (b *Blueprint) ID() {
	b.Columns = append(b.Columns, Column{
		Name:       "id",
		Type:       "UUID",
		PrimaryKey: true,
		Nullable:   false,
	})
}

func (b *Blueprint) String(name string) *Column {
	c := Column{Name: name, Type: "VARCHAR(255)"}
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Text(name string) *Column {
	c := Column{Name: name, Type: "TEXT"}
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Integer(name string) *Column {
	c := Column{Name: name, Type: "INTEGER"}
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Boolean(name string) *Column {
	c := Column{Name: name, Type: "BOOLEAN"}
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Binary(name string) *Column {
	c := Column{Name: name, Type: "BYTEA"} // Postgres BYTEA for binary
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Timestamp(name string) *Column {
	c := Column{Name: name, Type: "TIMESTAMP WITH TIME ZONE"}
	b.Columns = append(b.Columns, c)
	return &b.Columns[len(b.Columns)-1]
}

func (b *Blueprint) Timestamps() {
	b.Columns = append(b.Columns, Column{Name: "created_at", Type: "TIMESTAMP WITH TIME ZONE", Nullable: false})
	b.Columns = append(b.Columns, Column{Name: "updated_at", Type: "TIMESTAMP WITH TIME ZONE", Nullable: false})
}

func (b *Blueprint) SoftDeletes() {
	b.Columns = append(b.Columns, Column{Name: "deleted_at", Type: "TIMESTAMP WITH TIME ZONE", Nullable: true})
}

func (c *Column) NotNull() *Column {
	c.Nullable = false
	return c
}

func (c *Column) Unique() *Column {
	c.IsUnique = true
	return c
}

func (c *Column) Index() *Column {
	// For now, we'll mark it but in a real implementation we'd generate a separate CREATE INDEX SQL.
	return c
}

// -- SQL Generation (Basic Implementation) --

func (s *Schema) generateCreateTableSQL(b *Blueprint) string {
	var cols []string
	for _, c := range b.Columns {
		sql := fmt.Sprintf("%s %s", c.Name, c.Type)
		if c.PrimaryKey {
			sql += " PRIMARY KEY"
		}
		if !c.Nullable {
			sql += " NOT NULL"
		}
		if c.IsUnique {
			sql += " UNIQUE"
		}
		cols = append(cols, sql)
	}
	return fmt.Sprintf("CREATE TABLE %s (%s)", b.Name, strings.Join(cols, ", "))
}
