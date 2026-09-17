package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
