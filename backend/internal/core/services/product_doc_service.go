package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

type ProductDocService struct {
	repo ports.ProductDocRepository
}

func NewProductDocService(repo ports.ProductDocRepository) *ProductDocService {
	return &ProductDocService{repo: repo}
}

func (s *ProductDocService) Create(ctx context.Context, doc *domain.ProductDoc) error {
	return s.repo.Create(ctx, doc)
}

func (s *ProductDocService) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductDoc, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProductDocService) GetByProduct(ctx context.Context, productID uuid.UUID, docType string) ([]*domain.ProductDoc, error) {
	return s.repo.FindByProduct(ctx, productID, docType)
}

// GetPublicDocs returns all published docs for a product, grouped by section.
func (s *ProductDocService) GetPublicDocs(ctx context.Context, productID uuid.UUID, docType string) (map[string][]*domain.ProductDoc, error) {
	docs, err := s.repo.FindByProduct(ctx, productID, docType)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]*domain.ProductDoc)
	for _, doc := range docs {
		if doc.IsPublished {
			grouped[doc.SectionSlug] = append(grouped[doc.SectionSlug], doc)
		}
	}
	return grouped, nil
}

func (s *ProductDocService) GetDoc(ctx context.Context, productID uuid.UUID, docType string, sectionSlug string, docSlug string) (*domain.ProductDoc, error) {
	doc, err := s.repo.FindBySlug(ctx, productID, docType, sectionSlug, docSlug)
	if err != nil {
		return nil, err
	}
	if !doc.IsPublished {
		return nil, domain.ErrNotFound
	}
	return doc, nil
}

func (s *ProductDocService) Update(ctx context.Context, doc *domain.ProductDoc) error {
	return s.repo.Update(ctx, doc)
}

func (s *ProductDocService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
