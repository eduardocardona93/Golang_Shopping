package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
