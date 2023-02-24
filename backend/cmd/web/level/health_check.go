package level

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck
// Summary Show the status of server.
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
