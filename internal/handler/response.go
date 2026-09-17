package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/dto"
)

// listResponse is the standard envelope for paginated list endpoints.
type listResponse struct {
	Data interface{} `json:"data"`
	Meta dto.Meta    `json:"meta"`
}

// itemResponse is the standard envelope for single-resource responses.
type itemResponse struct {
	Data interface{} `json:"data"`
}

func respondItem(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, itemResponse{Data: data})
}

func respondList(c echo.Context, data interface{}, page, limit int, total int64) error {
	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: dto.NewMeta(page, limit, total),
	})
}
