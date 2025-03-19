package pages

import (
	"net/http"

	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func ClearAlerts(c echo.Context) error {
	internal.RemoveAllAlerts()

	// Redirect to the main page to refresh the alerts
	return c.Redirect(http.StatusSeeOther, "/")
}
