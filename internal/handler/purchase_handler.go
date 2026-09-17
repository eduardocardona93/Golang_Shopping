package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/internal/service"
)

// PurchaseHandler exposes HTTP endpoints for the Purchase resource.
type PurchaseHandler struct {
	service service.PurchaseService
}

// NewPurchaseHandler builds a PurchaseHandler.
func NewPurchaseHandler(s service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: s}
}

// Register wires the Purchase routes onto the given group.
func (h *PurchaseHandler) Register(g *echo.Group) {
	g.POST("/purchases", h.Create)
	g.GET("/purchases", h.List)
	g.GET("/purchases/:id", h.Get)
	g.PUT("/purchases/:id", h.Update)
	g.DELETE("/purchases/:id", h.Delete)
}

func (h *PurchaseHandler) Create(c echo.Context) error {
	var req dto.CreatePurchaseRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	purchase, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusCreated, dto.NewPurchaseResponse(purchase))
}

func (h *PurchaseHandler) Get(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	purchase, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewPurchaseResponse(purchase))
}

func (h *PurchaseHandler) List(c echo.Context) error {
	page, limit := paginationFromQuery(c)

	purchases, total, err := h.service.List(c.Request().Context(), repository.Pagination{Page: page, Limit: limit})
	if err != nil {
		return handleError(c, err)
	}

	return respondList(c, dto.NewPurchaseResponseList(purchases), page, limit, total)
}

func (h *PurchaseHandler) Update(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdatePurchaseRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	purchase, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewPurchaseResponse(purchase))
}

func (h *PurchaseHandler) Delete(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
