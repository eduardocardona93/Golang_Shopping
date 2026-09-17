package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

var validate = validator.New()

// bindAndValidate decodes the request body into req and runs struct
// validation tags against it, returning a rich validation error on failure.
func bindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := validate.Struct(req); err != nil {
		var fields []apperrors.Field
		for _, fe := range err.(validator.ValidationErrors) {
			fields = append(fields, apperrors.Field{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}
		return apperrors.NewValidationError(fields)
	}

	return nil
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "es requerido"
	case "gte":
		return "debe ser mayor o igual a " + fe.Param()
	case "gt":
		return "debe ser mayor a " + fe.Param()
	case "lte":
		return "debe ser menor o igual a " + fe.Param()
	case "max":
		return "debe tener como máximo " + fe.Param() + " caracteres"
	case "min":
		return "debe tener al menos " + fe.Param() + " caracteres"
	default:
		return "no es válido"
	}
}

// pathUUID parses a URL path parameter as a uuid.UUID.
func pathUUID(c echo.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid "+name)
	}
	return id, nil
}

// paginationFromQuery reads page/limit query parameters.
func paginationFromQuery(c echo.Context) (page, limit int) {
	page = queryInt(c, "page", 1)
	limit = queryInt(c, "limit", 20)
	return page, limit
}

func queryInt(c echo.Context, name string, fallback int) int {
	raw := c.QueryParam(name)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}
