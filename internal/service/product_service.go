package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
)

// ProductService encapsulates business logic for Product CRUD operations.
type ProductService interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, p repository.Pagination) ([]domain.Product, int64, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	products repository.ProductRepository
}

// NewProductService builds a ProductService backed by the given repository.
func NewProductService(products repository.ProductRepository) ProductService {
	return &productService{products: products}
}

func (s *productService) Create(ctx context.Context, req dto.CreateProductRequest) (*domain.Product, error) {
	product := req.ToDomain()
	if err := s.products.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.products.GetByID(ctx, id)
}

func (s *productService) List(ctx context.Context, p repository.Pagination) ([]domain.Product, int64, error) {
	return s.products.List(ctx, p)
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(product)

	if err := s.products.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.products.Delete(ctx, id)
}
