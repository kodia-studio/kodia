package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/pagination"
	"gorm.io/gorm"
)

// gormProduct is the GORM model for the products table.
// It mirrors the domain.Product entity but with GORM-specific tags.
// We keep this separate to avoid polluting the domain with framework concerns.
type gormProduct struct {
	ID          string     `gorm:"column:id;primaryKey"`
	Name        string     `gorm:"column:name;not null"`
	Description string     `gorm:"column:description"`
	Price       float64    `gorm:"column:price;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`
}

func (gormProduct) TableName() string { return "products" }

// toDomain converts a gormProduct to a domain.Product entity.
func (g *gormProduct) toDomain() *domain.Product {
	return &domain.Product{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Price:       g.Price,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
		DeletedAt:   g.DeletedAt,
	}
}

// fromDomainProduct converts a domain.Product to a gormProduct.
func fromDomainProduct(p *domain.Product) *gormProduct {
	return &gormProduct{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context, params *pagination.Params) ([]domain.Product, int64, error) {
	var items []gormProduct
	var total int64

	query := r.db.WithContext(ctx).Model(&gormProduct{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.Limit()).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	result := make([]domain.Product, len(items))
	for i, item := range items {
		result[i] = *item.toDomain()
	}

	return result, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	var item gormProduct
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return item.toDomain(), nil
}

func (r *productRepository) Create(ctx context.Context, item *domain.Product) error {
	m := fromDomainProduct(item)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *productRepository) Update(ctx context.Context, item *domain.Product) error {
	m := fromDomainProduct(item)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&gormProduct{}).Error
}
