package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/internal/service"
)

// ProductHandler exposes HTTP endpoints for the Product resource.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler builds a ProductHandler.
func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// Register wires the Product routes onto the given group.
func (h *ProductHandler) Register(g *echo.Group) {
	g.POST("/products", h.Create)
	g.GET("/products", h.List)
	g.GET("/products/:id", h.Get)
	g.PUT("/products/:id", h.Update)
	g.DELETE("/products/:id", h.Delete)
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	product, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusCreated, dto.NewProductResponse(product))
}

func (h *ProductHandler) Get(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	product, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewProductResponse(product))
}

func (h *ProductHandler) List(c echo.Context) error {
	page, limit := paginationFromQuery(c)

	products, total, err := h.service.List(c.Request().Context(), repository.Pagination{Page: page, Limit: limit})
	if err != nil {
		return handleError(c, err)
	}

	return respondList(c, dto.NewProductResponseList(products), page, limit, total)
}

func (h *ProductHandler) Update(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateProductRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	product, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewProductResponse(product))
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
