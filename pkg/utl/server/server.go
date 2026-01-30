package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health returns a 200 OK status with a JSON response {"status": "ok"}.
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
