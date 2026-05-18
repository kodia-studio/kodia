package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// StorageProvider is the interface for storage operations
// This same interface is also used by OrderService
type StorageProvider interface {
	UploadStream(ctx context.Context, key string, src io.Reader, size int64) error
	GetURL(ctx context.Context, key string) (string, error)
	GetTemporaryURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}

func makeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, s)
	s = strings.ReplaceAll(s, "--", "-")
	return strings.Trim(s, "-")
}

type ProductService struct {
	repo            ports.ProductRepository
	storageProvider StorageProvider
}

func NewProductService(repo ports.ProductRepository, storageProvider StorageProvider) *ProductService {
	return &ProductService{
		repo:            repo,
		storageProvider: storageProvider,
	}
}

func (s *ProductService) ListPublic(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error) {
	return s.repo.ListPublished(ctx, page, perPage)
}

func (s *ProductService) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *ProductService) ListAdmin(ctx context.Context, page, perPage int) ([]*domain.Product, int64, error) {
	return s.repo.List(ctx, false, page, perPage)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product) error {
	if product.Slug == "" {
		product.Slug = makeSlug(product.Name)
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, id string, fields map[string]any) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for key, value := range fields {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				product.Name = v
				if product.Slug == "" || !strings.Contains(product.Slug, "-") {
					product.Slug = makeSlug(v)
				}
			}
		case "slug":
			if v, ok := value.(string); ok && v != "" {
				product.Slug = v
			}
		case "tagline":
			if v, ok := value.(string); ok {
				product.Tagline = v
			}
		case "description":
			if v, ok := value.(string); ok {
				product.Description = v
			}
		case "type":
			if v, ok := value.(domain.ProductType); ok {
				product.Type = v
			}
		case "service_type":
			if v, ok := value.(domain.ServiceType); ok {
				product.ServiceType = &v
			}
		case "cover_url":
			if v, ok := value.(string); ok {
				product.CoverURL = v
			}
		case "tags":
			if v, ok := value.([]string); ok {
				product.Tags = v
			}
		case "is_published":
			if v, ok := value.(bool); ok {
				product.IsPublished = v
			}
		case "sort_order":
			if v, ok := value.(int); ok {
				product.SortOrder = v
			}
		case "meta_title":
			if v, ok := value.(string); ok {
				product.MetaTitle = v
			}
		case "meta_description":
			if v, ok := value.(string); ok {
				product.MetaDescription = v
			}
		case "og_image_url":
			if v, ok := value.(string); ok {
				product.OGImageURL = v
			}
		}
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) AddVariant(ctx context.Context, productID string, variant *domain.ProductVariant) error {
	variant.ProductID = productID
	return s.repo.CreateVariant(ctx, variant)
}

func (s *ProductService) UpdateVariant(ctx context.Context, productID, variantID string, fields map[string]any) (*domain.ProductVariant, error) {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var variant *domain.ProductVariant
	for i := range product.Variants {
		if product.Variants[i].ID == variantID {
			variant = &product.Variants[i]
			break
		}
	}

	if variant == nil {
		return nil, fmt.Errorf("variant not found")
	}

	for key, value := range fields {
		switch key {
		case "name":
			if v, ok := value.(string); ok {
				variant.Name = v
			}
		case "description":
			if v, ok := value.(string); ok {
				variant.Description = v
			}
		case "price":
			if v, ok := value.(int64); ok {
				variant.Price = v
			}
		case "promo_price":
			if value == nil {
				variant.PromoPrice = nil
			} else if v, ok := value.(int64); ok {
				variant.PromoPrice = &v
			}
		case "promo_start_at":
			if value == nil {
				variant.PromoStartAt = nil
			} else if _, ok := value.(interface{}); ok {
				// Handle time parsing from request
				variant.PromoStartAt = nil
			}
		case "promo_end_at":
			if value == nil {
				variant.PromoEndAt = nil
			} else if _, ok := value.(interface{}); ok {
				variant.PromoEndAt = nil
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				variant.IsActive = v
			}
		case "sort_order":
			if v, ok := value.(int); ok {
				variant.SortOrder = v
			}
		}
	}

	if err := s.repo.UpdateVariant(ctx, variant); err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *ProductService) DeleteVariant(ctx context.Context, productID, variantID string) error {
	return s.repo.DeleteVariant(ctx, variantID)
}

func (s *ProductService) UploadVariantFile(ctx context.Context, productID, variantID string, file *multipart.FileHeader) (string, error) {
	// Get variant to ensure it exists
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return "", err
	}

	var variant *domain.ProductVariant
	for i := range product.Variants {
		if product.Variants[i].ID == variantID {
			variant = &product.Variants[i]
			break
		}
	}

	if variant == nil {
		return "", fmt.Errorf("variant not found")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to R2 private bucket
	objectKey := fmt.Sprintf("products/%s/variants/%s/%s", productID, variantID, file.Filename)
	if err := s.storageProvider.UploadStream(ctx, objectKey, src, file.Size); err != nil {
		return "", err
	}

	// Update variant with file key
	variant.FileKey = objectKey
	if err := s.repo.UpdateVariant(ctx, variant); err != nil {
		return "", err
	}

	return objectKey, nil
}

func (s *ProductService) UploadCover(ctx context.Context, productID string, file *multipart.FileHeader) (string, error) {
	// Get product to ensure it exists
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return "", err
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to R2 public bucket
	objectKey := fmt.Sprintf("products/%s/cover/%s", productID, file.Filename)
	if err := s.storageProvider.UploadStream(ctx, objectKey, src, file.Size); err != nil {
		return "", err
	}

	// Get public URL
	publicURL, err := s.storageProvider.GetURL(ctx, objectKey)
	if err != nil {
		return "", err
	}

	// Update product with cover URL
	product.CoverURL = publicURL
	if err := s.repo.Update(ctx, product); err != nil {
		return "", err
	}

	return publicURL, nil
}
