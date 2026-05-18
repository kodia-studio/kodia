package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gormProduct struct {
	ID              string         `gorm:"column:id;primaryKey"`
	Slug            string         `gorm:"column:slug;uniqueIndex"`
	Name            string         `gorm:"column:name;not null"`
	Tagline         string         `gorm:"column:tagline"`
	Description     string         `gorm:"column:description;type:text"`
	Type            string         `gorm:"column:type;not null"`
	ServiceType     *string        `gorm:"column:service_type"`
	CoverURL        string         `gorm:"column:cover_url"`
	Tags            pq.StringArray `gorm:"column:tags;type:text[]"`
	IsPublished     bool           `gorm:"column:is_published;not null;default:false"`
	SortOrder       int            `gorm:"column:sort_order;default:0"`
	MetaTitle       string         `gorm:"column:meta_title"`
	MetaDescription string         `gorm:"column:meta_description;type:text"`
	OGImageURL      string         `gorm:"column:og_image_url"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormProduct) TableName() string { return "products" }

func (g *gormProduct) toDomain() *domain.Product {
	return &domain.Product{
		ID:              g.ID,
		Slug:            g.Slug,
		Name:            g.Name,
		Tagline:         g.Tagline,
		Description:     g.Description,
		Type:            domain.ProductType(g.Type),
		ServiceType:     (*domain.ServiceType)(g.ServiceType),
		CoverURL:        g.CoverURL,
		Tags:            g.Tags,
		IsPublished:     g.IsPublished,
		SortOrder:       g.SortOrder,
		MetaTitle:       g.MetaTitle,
		MetaDescription: g.MetaDescription,
		OGImageURL:      g.OGImageURL,
		CreatedAt:       g.CreatedAt,
		UpdatedAt:       g.UpdatedAt,
		DeletedAt:       func(t time.Time) *time.Time { if t.IsZero() { return nil }; return &t }(g.DeletedAt.Time),
	}
}

func fromDomainProduct(p *domain.Product) *gormProduct {
	return &gormProduct{
		ID:              p.ID,
		Slug:            p.Slug,
		Name:            p.Name,
		Tagline:         p.Tagline,
		Description:     p.Description,
		Type:            string(p.Type),
		ServiceType:     (*string)(p.ServiceType),
		CoverURL:        p.CoverURL,
		Tags:            p.Tags,
		IsPublished:     p.IsPublished,
		SortOrder:       p.SortOrder,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
		OGImageURL:      p.OGImageURL,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

type gormProductVariant struct {
	ID           string         `gorm:"column:id;primaryKey"`
	ProductID    string         `gorm:"column:product_id;index;not null"`
	Name         string         `gorm:"column:name;not null"`
	Description  string         `gorm:"column:description;type:text"`
	Price        int64          `gorm:"column:price;not null"`
	PromoPrice   *int64         `gorm:"column:promo_price"`
	PromoStartAt *time.Time     `gorm:"column:promo_start_at"`
	PromoEndAt   *time.Time     `gorm:"column:promo_end_at"`
	FileKey      string         `gorm:"column:file_key"`
	IsActive     bool           `gorm:"column:is_active;not null;default:true"`
	SortOrder    int            `gorm:"column:sort_order;default:0"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (gormProductVariant) TableName() string { return "product_variants" }

func (g *gormProductVariant) toDomain() *domain.ProductVariant {
	return &domain.ProductVariant{
		ID:           g.ID,
		ProductID:    g.ProductID,
		Name:         g.Name,
		Description:  g.Description,
		Price:        g.Price,
		PromoPrice:   g.PromoPrice,
		PromoStartAt: g.PromoStartAt,
		PromoEndAt:   g.PromoEndAt,
		FileKey:      g.FileKey,
		IsActive:     g.IsActive,
		SortOrder:    g.SortOrder,
		CreatedAt:    g.CreatedAt,
		UpdatedAt:    g.UpdatedAt,
		DeletedAt:    func(t time.Time) *time.Time { if t.IsZero() { return nil }; return &t }(g.DeletedAt.Time),
	}
}

func fromDomainVariant(v *domain.ProductVariant) *gormProductVariant {
	return &gormProductVariant{
		ID:           v.ID,
		ProductID:    v.ProductID,
		Name:         v.Name,
		Description:  v.Description,
		Price:        v.Price,
		PromoPrice:   v.PromoPrice,
		PromoStartAt: v.PromoStartAt,
		PromoEndAt:   v.PromoEndAt,
		FileKey:      v.FileKey,
		IsActive:     v.IsActive,
		SortOrder:    v.SortOrder,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(fromDomainProduct(product)).Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var g gormProduct
	err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	product := g.toDomain()
	if err := r.db.WithContext(ctx).Find(&[]gormProductVariant{}, "product_id = ?", id).Error; err == nil {
		var variants []gormProductVariant
		r.db.WithContext(ctx).Find(&variants, "product_id = ?", id)
		for i := range variants {
			product.Variants = append(product.Variants, *variants[i].toDomain())
		}
	}
	return product, nil
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var g gormProduct
	err := r.db.WithContext(ctx).First(&g, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	product := g.toDomain()
	var variants []gormProductVariant
	if err := r.db.WithContext(ctx).Find(&variants, "product_id = ?", g.ID).Error; err == nil {
		for i := range variants {
			product.Variants = append(product.Variants, *variants[i].toDomain())
		}
	}
	return product, nil
}

func (r *ProductRepository) List(ctx context.Context, published bool, page, perPage int) ([]*domain.Product, int64, error) {
	var gorms []gormProduct
	var total int64

	query := r.db.WithContext(ctx)
	if published {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Model(&gormProduct{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&gorms).Error; err != nil {
		return nil, 0, err
	}

	products := make([]*domain.Product, len(gorms))
	for i, g := range gorms {
		p := g.toDomain()
		var variants []gormProductVariant
		r.db.WithContext(ctx).Find(&variants, "product_id = ?", g.ID)
		for j := range variants {
			p.Variants = append(p.Variants, *variants[j].toDomain())
		}
		products[i] = p
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Model(&gormProduct{}).
		Where("id = ?", product.ID).
		Updates(fromDomainProduct(product)).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gormProduct{}, "id = ?", id).Error
}

func (r *ProductRepository) CreateVariant(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Create(fromDomainVariant(variant)).Error
}

func (r *ProductRepository) FindVariantByID(ctx context.Context, variantID string) (*domain.ProductVariant, error) {
	var g gormProductVariant
	err := r.db.WithContext(ctx).First(&g, "id = ?", variantID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.toDomain(), nil
}

func (r *ProductRepository) UpdateVariant(ctx context.Context, variant *domain.ProductVariant) error {
	gv := fromDomainVariant(variant)
	return r.db.WithContext(ctx).Model(&gormProductVariant{}).
		Where("id = ?", variant.ID).
		Select("name", "description", "price", "promo_price", "promo_start_at", "promo_end_at", "file_key", "is_active", "sort_order", "updated_at").
		Updates(gv).Error
}

func (r *ProductRepository) DeleteVariant(ctx context.Context, variantID string) error {
	return r.db.WithContext(ctx).Delete(&gormProductVariant{}, "id = ?", variantID).Error
}

func (r *ProductRepository) ListPublished(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error) {
	return r.List(ctx, true, page, perPage)
}

func (r *ProductRepository) ListVariantsByProductID(ctx context.Context, productID string) ([]*domain.ProductVariant, error) {
	var gorms []gormProductVariant
	if err := r.db.WithContext(ctx).Find(&gorms, "product_id = ?", productID).Error; err != nil {
		return nil, err
	}
	variants := make([]*domain.ProductVariant, len(gorms))
	for i := range gorms {
		variants[i] = gorms[i].toDomain()
	}
	return variants, nil
}
