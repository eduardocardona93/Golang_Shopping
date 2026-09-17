package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/dto"
	"github.com/eduardocardona93/golang_shopping/internal/repository"
	"github.com/eduardocardona93/golang_shopping/internal/service"
)

// UserHandler exposes HTTP endpoints for the User resource.
type UserHandler struct {
	service service.UserService
}

// NewUserHandler builds a UserHandler.
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register wires the User routes onto the given group.
func (h *UserHandler) Register(g *echo.Group) {
	g.POST("/users", h.Create)
	g.GET("/users", h.List)
	g.GET("/users/:id", h.Get)
	g.PUT("/users/:id", h.Update)
	g.DELETE("/users/:id", h.Delete)
}

func (h *UserHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	user, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusCreated, dto.NewUserResponse(user))
}

func (h *UserHandler) Get(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) List(c echo.Context) error {
	page, limit := paginationFromQuery(c)

	users, total, err := h.service.List(c.Request().Context(), repository.Pagination{Page: page, Limit: limit})
	if err != nil {
		return handleError(c, err)
	}

	return respondList(c, dto.NewUserResponseList(users), page, limit, total)
}

func (h *UserHandler) Update(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	var req dto.UpdateUserRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	user, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return handleError(c, err)
	}

	return respondItem(c, http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := pathUUID(c, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
