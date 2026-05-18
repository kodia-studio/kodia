package ports

import (
	"context"
	"mime/multipart"

	"github.com/kodia-studio/kodia/internal/core/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Product, error)
	List(ctx context.Context, published bool, page, perPage int) ([]*domain.Product, int64, error)
	ListPublished(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error

	CreateVariant(ctx context.Context, variant *domain.ProductVariant) error
	FindVariantByID(ctx context.Context, variantID string) (*domain.ProductVariant, error)
	UpdateVariant(ctx context.Context, variant *domain.ProductVariant) error
	DeleteVariant(ctx context.Context, variantID string) error
	ListVariantsByProductID(ctx context.Context, productID string) ([]*domain.ProductVariant, error)
}

type ProductService interface {
	ListPublic(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Product, error)

	ListAdmin(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, id string, fields map[string]any) (*domain.Product, error)
	Delete(ctx context.Context, id string) error

	AddVariant(ctx context.Context, productID string, variant *domain.ProductVariant) error
	UpdateVariant(ctx context.Context, productID, variantID string, fields map[string]any) (*domain.ProductVariant, error)
	DeleteVariant(ctx context.Context, productID, variantID string) error

	UploadCover(ctx context.Context, productID string, file *multipart.FileHeader) (string, error)
	UploadVariantFile(ctx context.Context, productID, variantID string, file *multipart.FileHeader) (string, error)
}
