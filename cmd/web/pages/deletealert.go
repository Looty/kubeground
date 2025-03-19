package pages

import (
	"net/http"
	"strconv"

	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func DeleteAlert(c echo.Context) error {
	indexStr := c.FormValue("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid index")
	}

	internal.RemoveAlertByIndex(index)

	// Redirect to the main page to refresh the alerts
	return c.Redirect(http.StatusSeeOther, "/")
}
