package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// PurchaseRepository defines persistence operations for domain.Purchase.
type PurchaseRepository interface {
	Create(ctx context.Context, purchase *domain.Purchase) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Purchase, error)
	List(ctx context.Context, p Pagination) ([]domain.Purchase, int64, error)
	// UpdateHeader updates the purchase's own fields (not its items).
	UpdateHeader(ctx context.Context, purchase *domain.Purchase) error
	// ReplaceItems atomically deletes the purchase's existing items and
	// inserts the provided ones.
	ReplaceItems(ctx context.Context, purchaseID uuid.UUID, items []domain.PurchaseItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	// WithTx returns a copy of the repository bound to the given transaction.
	WithTx(tx *gorm.DB) PurchaseRepository
}

type purchaseRepository struct {
	db *gorm.DB
}

// NewPurchaseRepository builds a GORM-backed PurchaseRepository.
func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(ctx context.Context, purchase *domain.Purchase) error {
	return translateWriteError(r.db.WithContext(ctx).Create(purchase).Error)
}

func (r *purchaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Purchase, error) {
	var purchase domain.Purchase
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		First(&purchase, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

func (r *purchaseRepository) List(ctx context.Context, p Pagination) ([]domain.Purchase, int64, error) {
	offset, limit := p.Normalize()

	var purchases []domain.Purchase
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Purchase{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Order("fecha DESC").
		Offset(offset).
		Limit(limit).
		Find(&purchases).Error; err != nil {
		return nil, 0, err
	}

	return purchases, total, nil
}

func (r *purchaseRepository) UpdateHeader(ctx context.Context, purchase *domain.Purchase) error {
	result := r.db.WithContext(ctx).Model(&domain.Purchase{}).
		Where("id = ?", purchase.ID).
		Updates(map[string]interface{}{
			"fecha":           purchase.Fecha,
			"impuesto":        purchase.Impuesto,
			"descuento_final": purchase.DescuentoFinal,
			"total":           purchase.Total,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *purchaseRepository) ReplaceItems(ctx context.Context, purchaseID uuid.UUID, items []domain.PurchaseItem) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("purchase_id = ?", purchaseID).
		Delete(&domain.PurchaseItem{}).Error; err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	for i := range items {
		items[i].PurchaseID = purchaseID
	}

	return translateWriteError(r.db.WithContext(ctx).Create(&items).Error)
}

func (r *purchaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Purchase{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *purchaseRepository) WithTx(tx *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: tx}
}
