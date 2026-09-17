package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Purchase represents a sale transaction (Compra) composed of multiple PurchaseItems.
type Purchase struct {
	ID     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User   User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Items  []PurchaseItem `gorm:"foreignKey:PurchaseID;references:ID;constraint:OnDelete:CASCADE" json:"items"`
	Fecha  time.Time      `gorm:"not null" json:"fecha"`
	// Impuesto is the total tax amount accumulated across all items.
	Impuesto float64 `gorm:"type:numeric(12,2);not null;default:0" json:"impuesto"`
	// DescuentoFinal is a flat monetary discount applied to the purchase total.
	DescuentoFinal float64        `gorm:"type:numeric(12,2);not null;default:0" json:"descuento_final"`
	Total          float64        `gorm:"type:numeric(12,2);not null;default:0" json:"total"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Purchase) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (Purchase) TableName() string {
	return "purchases"
}

// PurchaseItem represents a line item within a Purchase (ListaProducto).
type PurchaseItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PurchaseID uuid.UUID `gorm:"type:uuid;not null" json:"purchase_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Cantidad   int       `gorm:"not null" json:"cantidad"`
	// PrecioVenta is the unit sale price captured at purchase time.
	PrecioVenta float64 `gorm:"type:numeric(12,2);not null" json:"precio_venta"`
	// Descuento is a flat monetary discount applied to this line item.
	Descuento float64 `gorm:"type:numeric(12,2);not null;default:0" json:"descuento"`
	// Impuesto is the tax rate percentage applied to this line item, captured at purchase time.
	Impuesto  float64        `gorm:"type:numeric(5,2);not null;default:0" json:"impuesto"`
	Subtotal  float64        `gorm:"type:numeric(12,2);not null;default:0" json:"subtotal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *PurchaseItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (PurchaseItem) TableName() string {
	return "purchase_items"
}
