package pages

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func DisplayAboutQuest(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", "")
}
