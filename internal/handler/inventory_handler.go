package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/internal/service"
)

// InventoryHandler exposes HTTP endpoints for the Inventory resource.
type InventoryHandler struct {
	service service.InventoryService
}

// NewInventoryHandler builds an InventoryHandler.
func NewInventoryHandler(s service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

// Register wires the Inventory routes onto the given group.
func (h *InventoryHandler) Register(g *echo.Group) {
	g.POST("/inventories", h.Create)
	g.GET("/inventories", h.List)
	g.GET("/inventories/:id", h.Get)
	g.PUT("/inventories/:id", h.Update)
	g.DELETE("/inventories/:id", h.Delete)
}

func (h *InventoryHandler) Create(c echo.Context) error {
	var req dto.CreateInventoryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	inventory, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusCreated, dto.NewInventoryResponse(inventory))
}

func (h *InventoryHandler) Get(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	inventory, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewInventoryResponse(inventory))
}

func (h *InventoryHandler) List(c echo.Context) error {
	page, limit := paginationFromQuery(c)

	inventories, total, err := h.service.List(c.Request().Context(), repository.Pagination{Page: page, Limit: limit})
	if err != nil {
		return handleError(c, err)
	}

	return respondList(c, dto.NewInventoryResponseList(inventories), page, limit, total)
}

func (h *InventoryHandler) Update(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateInventoryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	inventory, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewInventoryResponse(inventory))
}

func (h *InventoryHandler) Delete(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
