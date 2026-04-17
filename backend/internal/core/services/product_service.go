package services

import (
	"context"

	"github.com/kodia/framework/backend/internal/core/domain"
	"github.com/kodia/framework/backend/internal/core/ports"
	"github.com/kodia/framework/backend/pkg/pagination"
	"go.uber.org/zap"
)

type productService struct {
	repo ports.ProductRepository
	log  *zap.Logger
}

func NewProductService(repo ports.ProductRepository, log *zap.Logger) ports.ProductService {
	return &productService{
		repo: repo,
		log:  log,
	}
}

func (s *productService) GetAll(ctx context.Context, params pagination.Params) ([]domain.Product, int64, error) {
	return s.repo.FindAll(ctx, params)
}

func (s *productService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
