package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
)

// PurchaseItemRequest is a single line item within a purchase request. The
// unit price and tax rate are always snapshotted from the current Product at
// the time of the purchase, so only the product, quantity and an optional
// flat discount are accepted from the client.
type PurchaseItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Cantidad  int       `json:"cantidad" validate:"required,gt=0"`
	Descuento float64   `json:"descuento" validate:"gte=0"`
}

// CreatePurchaseRequest is the payload to register a new Purchase.
type CreatePurchaseRequest struct {
	UserID         uuid.UUID             `json:"user_id" validate:"required"`
	Items          []PurchaseItemRequest `json:"items" validate:"required,min=1,dive"`
	DescuentoFinal float64               `json:"descuento_final" validate:"gte=0"`
}

// UpdatePurchaseRequest replaces a Purchase's date, final discount and full
// item list. Inventory is reconciled against the previous item list.
type UpdatePurchaseRequest struct {
	Fecha          time.Time             `json:"fecha" validate:"required"`
	Items          []PurchaseItemRequest `json:"items" validate:"required,min=1,dive"`
	DescuentoFinal float64               `json:"descuento_final" validate:"gte=0"`
}

// PurchaseItemResponse is the API representation of a purchase line item.
type PurchaseItemResponse struct {
	ID          uuid.UUID        `json:"id"`
	ProductID   uuid.UUID        `json:"product_id"`
	Product     *ProductResponse `json:"product,omitempty"`
	Cantidad    int              `json:"cantidad"`
	PrecioVenta float64          `json:"precio_venta"`
	Descuento   float64          `json:"descuento"`
	Impuesto    float64          `json:"impuesto"`
	Subtotal    float64          `json:"subtotal"`
}

// PurchaseResponse is the API representation of a Purchase.
type PurchaseResponse struct {
	ID             uuid.UUID              `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	User           *UserResponse          `json:"user,omitempty"`
	Items          []PurchaseItemResponse `json:"items"`
	Fecha          time.Time              `json:"fecha"`
	Impuesto       float64                `json:"impuesto"`
	DescuentoFinal float64                `json:"descuento_final"`
	Total          float64                `json:"total"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func NewPurchaseResponse(p *domain.Purchase) PurchaseResponse {
	resp := PurchaseResponse{
		ID:             p.ID,
		UserID:         p.UserID,
		Items:          make([]PurchaseItemResponse, 0, len(p.Items)),
		Fecha:          p.Fecha,
		Impuesto:       p.Impuesto,
		DescuentoFinal: p.DescuentoFinal,
		Total:          p.Total,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	if p.User.ID != uuid.Nil {
		u := NewUserResponse(&p.User)
		resp.User = &u
	}

	for i := range p.Items {
		item := p.Items[i]
		itemResp := PurchaseItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			Cantidad:    item.Cantidad,
			PrecioVenta: item.PrecioVenta,
			Descuento:   item.Descuento,
			Impuesto:    item.Impuesto,
			Subtotal:    item.Subtotal,
		}
		if item.Product.ID != uuid.Nil {
			p := NewProductResponse(&item.Product)
			itemResp.Product = &p
		}
		resp.Items = append(resp.Items, itemResp)
	}

	return resp
}

func NewPurchaseResponseList(purchases []domain.Purchase) []PurchaseResponse {
	out := make([]PurchaseResponse, 0, len(purchases))
	for i := range purchases {
		out = append(out, NewPurchaseResponse(&purchases[i]))
	}
	return out
}
