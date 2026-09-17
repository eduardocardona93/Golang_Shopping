package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// errorResponse is the standard envelope for error responses.
type errorResponse struct {
	Error  string            `json:"error"`
	Fields []apperrors.Field `json:"fields,omitempty"`
	Detail *stockErrorDetail `json:"detail,omitempty"`
}

type stockErrorDetail struct {
	ProductID string `json:"product_id"`
	Requested int    `json:"requested"`
	Available int    `json:"available"`
}

// handleError maps a service/repository error to the appropriate HTTP
// status code and a consistent JSON error body.
func handleError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return c.JSON(httpErr.Code, errorResponse{Error: fmt.Sprint(httpErr.Message)})
	}

	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		return c.JSON(http.StatusUnprocessableEntity, errorResponse{
			Error:  "validation failed",
			Fields: validationErr.Fields,
		})
	}

	var stockErr *apperrors.InsufficientStockError
	if errors.As(err, &stockErr) {
		return c.JSON(http.StatusConflict, errorResponse{
			Error: "insufficient stock",
			Detail: &stockErrorDetail{
				ProductID: stockErr.ProductID,
				Requested: stockErr.Requested,
				Available: stockErr.Available,
			},
		})
	}

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return c.JSON(http.StatusNotFound, errorResponse{Error: "resource not found"})
	case errors.Is(err, apperrors.ErrConflict):
		return c.JSON(http.StatusConflict, errorResponse{Error: "resource already exists"})
	case errors.Is(err, apperrors.ErrValidation):
		return c.JSON(http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
