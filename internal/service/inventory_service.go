package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
)

// InventoryService encapsulates business logic for Inventory CRUD operations.
type InventoryService interface {
	Create(ctx context.Context, req dto.CreateInventoryRequest) (*domain.Inventory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	List(ctx context.Context, p repository.Pagination) ([]domain.Inventory, int64, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (*domain.Inventory, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type inventoryService struct {
	inventories repository.InventoryRepository
	products    repository.ProductRepository
}

// NewInventoryService builds an InventoryService backed by the given repositories.
func NewInventoryService(inventories repository.InventoryRepository, products repository.ProductRepository) InventoryService {
	return &inventoryService{inventories: inventories, products: products}
}

func (s *inventoryService) Create(ctx context.Context, req dto.CreateInventoryRequest) (*domain.Inventory, error) {
	// Ensure the referenced product exists before creating stock for it.
	if _, err := s.products.GetByID(ctx, req.ProductID); err != nil {
		return nil, err
	}

	inventory := req.ToDomain()
	if err := s.inventories.Create(ctx, inventory); err != nil {
		return nil, err
	}
	return s.inventories.GetByID(ctx, inventory.ID)
}

func (s *inventoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return s.inventories.GetByID(ctx, id)
}

func (s *inventoryService) List(ctx context.Context, p repository.Pagination) ([]domain.Inventory, int64, error) {
	return s.inventories.List(ctx, p)
}

func (s *inventoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (*domain.Inventory, error) {
	inventory, err := s.inventories.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	req.ApplyTo(inventory)

	if err := s.inventories.Update(ctx, inventory); err != nil {
		return nil, err
	}
	return s.inventories.GetByID(ctx, id)
}

func (s *inventoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.inventories.Delete(ctx, id)
}
