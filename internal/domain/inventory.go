package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// Inventory tracks stock levels for a Product (Inventario).
type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Cantidad  int       `gorm:"not null;default:0" json:"cantidad"`
	// MinCantidadAlerta is the minimum quantity threshold used to trigger a low-stock alert.
	MinCantidadAlerta int            `gorm:"not null;default:0" json:"min_cantidad_alerta"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (Inventory) TableName() string {
	return "inventories"
}

// Reserve decrements the available quantity by qty, rejecting the
// reservation with apperrors.ErrInsufficientStock if there isn't enough
// stock. Callers are expected to have loaded the Inventory with a row lock
// (see repository.InventoryRepository.GetByProductIDForUpdate) and to
// persist the updated Cantidad afterward.
func (i *Inventory) Reserve(qty int) error {
	if i.Cantidad < qty {
		return apperrors.NewInsufficientStockError(i.ProductID.String(), qty, i.Cantidad)
	}
	i.Cantidad -= qty
	return nil
}

// Release restores qty units back to the available quantity, undoing a
// previous Reserve (e.g. when a purchase is updated or deleted).
func (i *Inventory) Release(qty int) {
	i.Cantidad += qty
}
