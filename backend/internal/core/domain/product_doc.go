package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductDoc struct {
	ID              uuid.UUID
	ProductID       uuid.UUID
	DocType         string // "developer" | "user_guide"
	SectionTitle    string
	SectionSlug     string
	Title           string
	Slug            string
	Content         string // raw markdown
	SortOrder       int
	SectionOrder    int
	IsPublished     bool
	PublishedAt     *time.Time
	MetaTitle       string
	MetaDescription string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
