package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"gorm.io/gorm"
)

type gormProductDoc struct {
	ID              string         `gorm:"column:id;primaryKey"`
	ProductID       string         `gorm:"column:product_id;index;not null"`
	DocType         string         `gorm:"column:doc_type;not null"` // "developer" | "user_guide"
	SectionTitle    string         `gorm:"column:section_title;not null"`
	SectionSlug     string         `gorm:"column:section_slug;not null"`
	Title           string         `gorm:"column:title;not null"`
	Slug            string         `gorm:"column:slug;not null"`
	Content         string         `gorm:"column:content;type:text;not null"`
	SortOrder       int            `gorm:"column:sort_order;default:0"`
	SectionOrder    int            `gorm:"column:section_order;default:0"`
	IsPublished     bool           `gorm:"column:is_published;default:false"`
	PublishedAt     *time.Time     `gorm:"column:published_at"`
	MetaTitle       string         `gorm:"column:meta_title"`
	MetaDescription string         `gorm:"column:meta_description"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormProductDoc) TableName() string { return "product_docs" }

func (g *gormProductDoc) toDomain() *domain.ProductDoc {
	return &domain.ProductDoc{
		ID:              uuid.MustParse(g.ID),
		ProductID:       uuid.MustParse(g.ProductID),
		DocType:         g.DocType,
		SectionTitle:    g.SectionTitle,
		SectionSlug:     g.SectionSlug,
		Title:           g.Title,
		Slug:            g.Slug,
		Content:         g.Content,
		SortOrder:       g.SortOrder,
		SectionOrder:    g.SectionOrder,
		IsPublished:     g.IsPublished,
		PublishedAt:     g.PublishedAt,
		MetaTitle:       g.MetaTitle,
		MetaDescription: g.MetaDescription,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
	}
}

func fromDomainProductDoc(pd *domain.ProductDoc) *gormProductDoc {
	return &gormProductDoc{
		ID:              pd.ID.String(),
		ProductID:       pd.ProductID.String(),
		DocType:         pd.DocType,
		SectionTitle:    pd.SectionTitle,
		SectionSlug:     pd.SectionSlug,
		Title:           pd.Title,
		Slug:            pd.Slug,
		Content:         pd.Content,
		SortOrder:       pd.SortOrder,
		SectionOrder:    pd.SectionOrder,
		IsPublished:     pd.IsPublished,
		PublishedAt:     pd.PublishedAt,
		MetaTitle:       pd.MetaTitle,
		MetaDescription: pd.MetaDescription,
		CreatedAt:       pd.CreatedAt,
		UpdatedAt:       pd.UpdatedAt,
	}
}

type productDocRepository struct {
	db *gorm.DB
}

func NewProductDocRepository(db *gorm.DB) *productDocRepository {
	return &productDocRepository{db: db}
}

func (r *productDocRepository) Create(ctx context.Context, doc *domain.ProductDoc) error {
	gormDoc := fromDomainProductDoc(doc)
	return r.db.WithContext(ctx).Create(gormDoc).Error
}

func (r *productDocRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductDoc, error) {
	var doc gormProductDoc
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return doc.toDomain(), nil
}

func (r *productDocRepository) FindByProduct(ctx context.Context, productID uuid.UUID, docType string) ([]*domain.ProductDoc, error) {
	var docs []gormProductDoc
	query := r.db.WithContext(ctx).Where("product_id = ? AND doc_type = ?", productID.String(), docType)
	if err := query.Order("section_order ASC, sort_order ASC").Find(&docs).Error; err != nil {
		return nil, err
	}
	result := make([]*domain.ProductDoc, len(docs))
	for i, doc := range docs {
		result[i] = doc.toDomain()
	}
	return result, nil
}

func (r *productDocRepository) FindBySlug(ctx context.Context, productID uuid.UUID, docType string, sectionSlug string, docSlug string) (*domain.ProductDoc, error) {
	var doc gormProductDoc
	err := r.db.WithContext(ctx).Where(
		"product_id = ? AND doc_type = ? AND section_slug = ? AND slug = ?",
		productID.String(), docType, sectionSlug, docSlug,
	).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return doc.toDomain(), nil
}

func (r *productDocRepository) Update(ctx context.Context, doc *domain.ProductDoc) error {
	gormDoc := fromDomainProductDoc(doc)
	return r.db.WithContext(ctx).Save(gormDoc).Error
}

func (r *productDocRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&gormProductDoc{}).Error
}
