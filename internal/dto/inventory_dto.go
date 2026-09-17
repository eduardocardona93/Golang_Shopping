package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
)

// CreateInventoryRequest is the payload to create a new Inventory record for a product.
type CreateInventoryRequest struct {
	ProductID         uuid.UUID `json:"product_id" validate:"required"`
	Cantidad          int       `json:"cantidad" validate:"gte=0"`
	MinCantidadAlerta int       `json:"min_cantidad_alerta" validate:"gte=0"`
}

// UpdateInventoryRequest is the payload to update an existing Inventory record.
type UpdateInventoryRequest struct {
	Cantidad          int `json:"cantidad" validate:"gte=0"`
	MinCantidadAlerta int `json:"min_cantidad_alerta" validate:"gte=0"`
}

// InventoryResponse is the API representation of an Inventory record.
type InventoryResponse struct {
	ID                uuid.UUID        `json:"id"`
	ProductID         uuid.UUID        `json:"product_id"`
	Product           *ProductResponse `json:"product,omitempty"`
	Cantidad          int              `json:"cantidad"`
	MinCantidadAlerta int              `json:"min_cantidad_alerta"`
	LowStock          bool             `json:"low_stock"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

func (r CreateInventoryRequest) ToDomain() *domain.Inventory {
	return &domain.Inventory{
		ProductID:         r.ProductID,
		Cantidad:          r.Cantidad,
		MinCantidadAlerta: r.MinCantidadAlerta,
	}
}

func (r UpdateInventoryRequest) ApplyTo(i *domain.Inventory) {
	i.Cantidad = r.Cantidad
	i.MinCantidadAlerta = r.MinCantidadAlerta
}

func NewInventoryResponse(i *domain.Inventory) InventoryResponse {
	resp := InventoryResponse{
		ID:                i.ID,
		ProductID:         i.ProductID,
		Cantidad:          i.Cantidad,
		MinCantidadAlerta: i.MinCantidadAlerta,
		LowStock:          i.Cantidad <= i.MinCantidadAlerta,
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}
	if i.Product.ID != uuid.Nil {
		p := NewProductResponse(&i.Product)
		resp.Product = &p
	}
	return resp
}

func NewInventoryResponseList(inventories []domain.Inventory) []InventoryResponse {
	out := make([]InventoryResponse, 0, len(inventories))
	for i := range inventories {
		out = append(out, NewInventoryResponse(&inventories[i]))
	}
	return out
}
