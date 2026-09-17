package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// PurchaseService encapsulates the orchestration of Purchase CRUD
// operations: fetching and locking the rows involved, delegating pricing,
// stock and totals rules to the domain layer, and persisting the result.
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

func (s *purchaseService) Create(ctx context.Context, req dto.CreatePurchaseRequest) (*domain.Purchase, error) {
	if _, err := s.users.GetByID(ctx, req.UserID); err != nil {
		return nil, err
	}

	var purchase *domain.Purchase

	err := s.db.Transaction(func(tx *gorm.DB) error {
		items, err := s.reserveStockAndPriceItems(ctx, tx, req.Items)
		if err != nil {
			return err
		}

		purchase, err = domain.NewPurchase(req.UserID, items, req.DescuentoFinal)
		if err != nil {
			return err
		}

		return s.purchases.WithTx(tx).Create(ctx, purchase)
	})
	if err != nil {
		return nil, err
	}

	return s.purchases.GetByID(ctx, purchase.ID)
}

// reserveStockAndPriceItems locks the Inventory row for every requested
// line item, delegates the stock reservation and pricing to the domain
// layer, and persists the resulting quantities. It must run inside an
// active transaction (tx) so that the row locks and stock decrements are
// atomic with the rest of the purchase write.
func (s *purchaseService) reserveStockAndPriceItems(ctx context.Context, tx *gorm.DB, itemReqs []dto.PurchaseItemRequest) ([]domain.PurchaseItem, error) {
	products := s.products.WithTx(tx)
	inventories := s.inventories.WithTx(tx)

	items := make([]domain.PurchaseItem, 0, len(itemReqs))

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

		if err := inventory.Reserve(itemReq.Cantidad); err != nil {
			return nil, err
		}

		item, err := domain.NewPurchaseItem(product, itemReq.Cantidad, itemReq.Descuento)
		if err != nil {
			return nil, err
		}
		items = append(items, item)

		if err := inventories.UpdateQuantity(ctx, inventory.ID, inventory.Cantidad); err != nil {
			return nil, err
		}
	}

	return items, nil
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

		inventory.Release(item.Cantidad)

		if err := inventories.UpdateQuantity(ctx, inventory.ID, inventory.Cantidad); err != nil {
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

		items, err := s.reserveStockAndPriceItems(ctx, tx, req.Items)
		if err != nil {
			return err
		}

		if err := existing.ApplyItems(items, req.DescuentoFinal); err != nil {
			return err
		}
		existing.Fecha = req.Fecha

		if err := s.purchases.WithTx(tx).ReplaceItems(ctx, id, existing.Items); err != nil {
			return err
		}

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
