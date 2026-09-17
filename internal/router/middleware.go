package router

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// registerMiddleware attaches the standard, industry-common middleware
// stack: request logging, panic recovery, request IDs and permissive CORS
// suitable for an API consumed by separate frontends.
func registerMiddleware(e *echo.Echo) {
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.BodyLimit("2M"))
}
