package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// Health returns a simple status "ok" for health checks
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// NOTE: This file may contain other existing code which is not shown here.
// The Health function is added to the existing content of this file.
