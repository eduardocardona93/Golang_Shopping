package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// InventoryRepository defines persistence operations for domain.Inventory.
type InventoryRepository interface {
	Create(ctx context.Context, inventory *domain.Inventory) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error)
	// GetByProductIDForUpdate locks the inventory row for the duration of the
	// enclosing transaction, preventing concurrent purchases from
	// overselling the same product.
	GetByProductIDForUpdate(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error)
	List(ctx context.Context, p Pagination) ([]domain.Inventory, int64, error)
	Update(ctx context.Context, inventory *domain.Inventory) error
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	Delete(ctx context.Context, id uuid.UUID) error
	// WithTx returns a copy of the repository bound to the given transaction.
	WithTx(tx *gorm.DB) InventoryRepository
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository builds a GORM-backed InventoryRepository.
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	return translateWriteError(r.db.WithContext(ctx).Create(inventory).Error)
}

func (r *inventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.WithContext(ctx).Preload("Product").First(&inventory, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.WithContext(ctx).Preload("Product").First(&inventory, "product_id = ?", productID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetByProductIDForUpdate(ctx context.Context, productID uuid.UUID) (*domain.Inventory, error) {
	var inventory domain.Inventory
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&inventory, "product_id = ?", productID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) List(ctx context.Context, p Pagination) ([]domain.Inventory, int64, error) {
	offset, limit := p.Normalize()

	var inventories []domain.Inventory
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Inventory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

func (r *inventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	result := r.db.WithContext(ctx).Model(&domain.Inventory{}).
		Where("id = ?", inventory.ID).
		Updates(map[string]interface{}{
			"cantidad":            inventory.Cantidad,
			"min_cantidad_alerta": inventory.MinCantidadAlerta,
		})
	if err := translateWriteError(result.Error); err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *inventoryRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	result := r.db.WithContext(ctx).Model(&domain.Inventory{}).
		Where("id = ?", id).
		Update("cantidad", quantity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *inventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Inventory{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *inventoryRepository) WithTx(tx *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: tx}
}
