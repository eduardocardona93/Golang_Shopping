package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
)

// CreateProductRequest is the payload to create a new Product.
type CreateProductRequest struct {
	Nombre      string  `json:"nombre" validate:"required,min=1,max=150"`
	Descripcion string  `json:"descripcion" validate:"omitempty"`
	SKU         string  `json:"sku" validate:"required,min=1,max=50"`
	Costo       float64 `json:"costo" validate:"required,gte=0"`
	Precio      float64 `json:"precio" validate:"required,gte=0"`
	Impuesto    float64 `json:"impuesto" validate:"gte=0,lte=100"`
}

// UpdateProductRequest is the payload to update an existing Product.
type UpdateProductRequest struct {
	Nombre      string  `json:"nombre" validate:"required,min=1,max=150"`
	Descripcion string  `json:"descripcion" validate:"omitempty"`
	SKU         string  `json:"sku" validate:"required,min=1,max=50"`
	Costo       float64 `json:"costo" validate:"required,gte=0"`
	Precio      float64 `json:"precio" validate:"required,gte=0"`
	Impuesto    float64 `json:"impuesto" validate:"gte=0,lte=100"`
}

// ProductResponse is the API representation of a Product.
type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	SKU         string    `json:"sku"`
	Costo       float64   `json:"costo"`
	Precio      float64   `json:"precio"`
	Impuesto    float64   `json:"impuesto"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r CreateProductRequest) ToDomain() *domain.Product {
	return &domain.Product{
		Nombre:      r.Nombre,
		Descripcion: r.Descripcion,
		SKU:         r.SKU,
		Costo:       r.Costo,
		Precio:      r.Precio,
		Impuesto:    r.Impuesto,
	}
}

func (r UpdateProductRequest) ApplyTo(p *domain.Product) {
	p.Nombre = r.Nombre
	p.Descripcion = r.Descripcion
	p.SKU = r.SKU
	p.Costo = r.Costo
	p.Precio = r.Precio
	p.Impuesto = r.Impuesto
}

func NewProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		SKU:         p.SKU,
		Costo:       p.Costo,
		Precio:      p.Precio,
		Impuesto:    p.Impuesto,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func NewProductResponseList(products []domain.Product) []ProductResponse {
	out := make([]ProductResponse, 0, len(products))
	for i := range products {
		out = append(out, NewProductResponse(&products[i]))
	}
	return out
}
