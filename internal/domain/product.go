package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// Product represents an item that can be sold (Producto).
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nombre      string    `gorm:"type:varchar(150);not null" json:"nombre"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	SKU         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Costo       float64   `gorm:"type:numeric(12,2);not null" json:"costo"`
	Precio      float64   `gorm:"type:numeric(12,2);not null" json:"precio"`
	// Impuesto is the default tax rate percentage applied to this product (e.g. 19 = 19%).
	Impuesto  float64        `gorm:"type:numeric(5,2);not null;default:0" json:"impuesto"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (Product) TableName() string {
	return "products"
}

// PriceLine prices `cantidad` units of this product with a flat `descuento`,
// applying the product's current price and tax rate, and returns the tax
// amount and the final subtotal for the line. It rejects a discount larger
// than the line's gross amount.
func (p *Product) PriceLine(cantidad int, descuento float64) (taxAmount, subtotal float64, err error) {
	grossAmount := p.Precio * float64(cantidad)
	if descuento > grossAmount {
		return 0, 0, apperrors.NewValidationError([]apperrors.Field{
			{Field: "descuento", Message: "descuento no puede ser mayor al subtotal del producto"},
		})
	}

	taxableAmount := grossAmount - descuento
	taxAmount = taxableAmount * (p.Impuesto / 100)
	subtotal = taxableAmount + taxAmount

	return taxAmount, subtotal, nil
}
