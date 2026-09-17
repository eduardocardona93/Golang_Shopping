package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// PurchaseService encapsulates business logic for Purchase CRUD operations,
// including reserving/releasing stock in Inventory as purchases are
// created, updated or deleted.
type PurchaseService interface {
	Create(ctx context.Context, req dto.CreatePurchaseRequest) (*domain.Purchase, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Purchase, error)
	List(ctx context.Context, p repository.Pagination) ([]domain.Purchase, int64, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdatePurchaseRequest) (*domain.Purchase, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type purchaseService struct {
	db          *gorm.DB
	purchases   repository.PurchaseRepository
	inventories repository.InventoryRepository
	products    repository.ProductRepository
	users       repository.UserRepository
}

// NewPurchaseService builds a PurchaseService. db is used to open the
// transactions that atomically reconcile Purchase records against
// Inventory stock.
func NewPurchaseService(
	db *gorm.DB,
	purchases repository.PurchaseRepository,
	inventories repository.InventoryRepository,
	products repository.ProductRepository,
	users repository.UserRepository,
) PurchaseService {
	return &purchaseService{
		db:          db,
		purchases:   purchases,
		inventories: inventories,
		products:    products,
		users:       users,
	}
}

// lineTotals holds the computed monetary breakdown for a purchase's items.
type lineTotals struct {
	items         []domain.PurchaseItem
	totalImpuesto float64
	totalItems    float64 // sum of item subtotals (already includes tax and item discounts)
}

func (s *purchaseService) Create(ctx context.Context, req dto.CreatePurchaseRequest) (*domain.Purchase, error) {
	if _, err := s.users.GetByID(ctx, req.UserID); err != nil {
		return nil, err
	}

	var purchase domain.Purchase

	err := s.db.Transaction(func(tx *gorm.DB) error {
		totals, err := s.reserveStockAndBuildItems(ctx, tx, req.Items)
		if err != nil {
			return err
		}

		if req.DescuentoFinal > totals.totalItems {
			return apperrors.NewValidationError([]apperrors.Field{
				{Field: "descuento_final", Message: "descuento_final no puede ser mayor al total de la compra"},
			})
		}

		purchase = domain.Purchase{
			UserID:         req.UserID,
			Items:          totals.items,
			Fecha:          time.Now(),
			Impuesto:       totals.totalImpuesto,
			DescuentoFinal: req.DescuentoFinal,
			Total:          totals.totalItems - req.DescuentoFinal,
		}

		return s.purchases.WithTx(tx).Create(ctx, &purchase)
	})
	if err != nil {
		return nil, err
	}

	return s.purchases.GetByID(ctx, purchase.ID)
}

// reserveStockAndBuildItems validates and locks inventory for every
// requested line item, decrements available stock and returns the
// fully-priced PurchaseItem records plus aggregate totals. It must run
// inside an active transaction (tx) so that the row locks and stock
// decrements are atomic with the rest of the purchase write.
func (s *purchaseService) reserveStockAndBuildItems(ctx context.Context, tx *gorm.DB, itemReqs []dto.PurchaseItemRequest) (*lineTotals, error) {
	products := s.products.WithTx(tx)
	inventories := s.inventories.WithTx(tx)

	totals := &lineTotals{items: make([]domain.PurchaseItem, 0, len(itemReqs))}

	for _, itemReq := range itemReqs {
		product, err := products.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}

		// Lock the inventory row for this product so concurrent purchases
		// cannot both read the same available quantity and oversell it.
		inventory, err := inventories.GetByProductIDForUpdate(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}

		if inventory.Cantidad < itemReq.Cantidad {
			return nil, apperrors.NewInsufficientStockError(product.ID.String(), itemReq.Cantidad, inventory.Cantidad)
		}

		grossAmount := product.Precio * float64(itemReq.Cantidad)
		if itemReq.Descuento > grossAmount {
			return nil, apperrors.NewValidationError([]apperrors.Field{
				{Field: "descuento", Message: "descuento no puede ser mayor al subtotal del producto"},
			})
		}

		taxableAmount := grossAmount - itemReq.Descuento
		taxAmount := taxableAmount * (product.Impuesto / 100)
		subtotal := taxableAmount + taxAmount

		totals.items = append(totals.items, domain.PurchaseItem{
			ProductID:   product.ID,
			Cantidad:    itemReq.Cantidad,
			PrecioVenta: product.Precio,
			Descuento:   itemReq.Descuento,
			Impuesto:    product.Impuesto,
			Subtotal:    subtotal,
		})
		totals.totalImpuesto += taxAmount
		totals.totalItems += subtotal

		if err := inventories.UpdateQuantity(ctx, inventory.ID, inventory.Cantidad-itemReq.Cantidad); err != nil {
			return nil, err
		}
	}

	return totals, nil
}

// releaseStock restores previously reserved inventory quantities, used when
// a purchase is updated (its old items are released before the new ones are
// reserved) or deleted.
func (s *purchaseService) releaseStock(ctx context.Context, tx *gorm.DB, items []domain.PurchaseItem) error {
	inventories := s.inventories.WithTx(tx)

	for _, item := range items {
		inventory, err := inventories.GetByProductIDForUpdate(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				// The inventory record for this product no longer exists;
				// nothing to restore.
				continue
			}
			return err
		}

		if err := inventories.UpdateQuantity(ctx, inventory.ID, inventory.Cantidad+item.Cantidad); err != nil {
			return err
		}
	}

	return nil
}

func (s *purchaseService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Purchase, error) {
	return s.purchases.GetByID(ctx, id)
}

func (s *purchaseService) List(ctx context.Context, p repository.Pagination) ([]domain.Purchase, int64, error) {
	return s.purchases.List(ctx, p)
}

func (s *purchaseService) Update(ctx context.Context, id uuid.UUID, req dto.UpdatePurchaseRequest) (*domain.Purchase, error) {
	existing, err := s.purchases.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.releaseStock(ctx, tx, existing.Items); err != nil {
			return err
		}

		totals, err := s.reserveStockAndBuildItems(ctx, tx, req.Items)
		if err != nil {
			return err
		}

		if req.DescuentoFinal > totals.totalItems {
			return apperrors.NewValidationError([]apperrors.Field{
				{Field: "descuento_final", Message: "descuento_final no puede ser mayor al total de la compra"},
			})
		}

		if err := s.purchases.WithTx(tx).ReplaceItems(ctx, id, totals.items); err != nil {
			return err
		}

		existing.Fecha = req.Fecha
		existing.Impuesto = totals.totalImpuesto
		existing.DescuentoFinal = req.DescuentoFinal
		existing.Total = totals.totalItems - req.DescuentoFinal

		return s.purchases.WithTx(tx).UpdateHeader(ctx, existing)
	})
	if err != nil {
		return nil, err
	}

	return s.purchases.GetByID(ctx, id)
}

func (s *purchaseService) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := s.purchases.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.releaseStock(ctx, tx, existing.Items); err != nil {
			return err
		}

		if err := tx.WithContext(ctx).
			Where("purchase_id = ?", id).
			Delete(&domain.PurchaseItem{}).Error; err != nil {
			return err
		}

		return s.purchases.WithTx(tx).Delete(ctx, id)
	})
}
