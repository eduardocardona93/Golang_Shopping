package router

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/eduardocardona93/golang_shopping/internal/handler"
)

// Handlers bundles every resource handler the router needs to wire up.
type Handlers struct {
	User      *handler.UserHandler
	Product   *handler.ProductHandler
	Inventory *handler.InventoryHandler
	Purchase  *handler.PurchaseHandler
}

// New builds the Echo instance with all middleware and routes registered.
func New(h Handlers) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErrorHandler

	registerMiddleware(e)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")
	h.User.Register(api)
	h.Product.Register(api)
	h.Inventory.Register(api)
	h.Purchase.Register(api)

	return e
}
