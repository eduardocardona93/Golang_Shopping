package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// httpErrorHandler is the Echo-level fallback for errors that bubble up
// without already writing a response (e.g. 404 route not found, method not
// allowed, or a raw *echo.HTTPError returned before reaching a handler's
// own error mapping). It keeps the JSON error envelope consistent with
// internal/handler.handleError.
func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "internal server error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprint(he.Message)
	}

	if jsonErr := c.JSON(code, map[string]string{"error": message}); jsonErr != nil {
		c.Logger().Error(jsonErr)
	}
}
