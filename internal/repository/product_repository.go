package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// ProductRepository defines persistence operations for domain.Product.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, p Pagination) ([]domain.Product, int64, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	// WithTx returns a copy of the repository bound to the given transaction.
	WithTx(tx *gorm.DB) ProductRepository
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository builds a GORM-backed ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return translateWriteError(r.db.WithContext(ctx).Create(product).Error)
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) List(ctx context.Context, p Pagination) ([]domain.Product, int64, error) {
	offset, limit := p.Normalize()

	var products []domain.Product
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	result := r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]interface{}{
			"nombre":      product.Nombre,
			"descripcion": product.Descripcion,
			"sku":         product.SKU,
			"costo":       product.Costo,
			"precio":      product.Precio,
			"impuesto":    product.Impuesto,
		})
	if err := translateWriteError(result.Error); err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
