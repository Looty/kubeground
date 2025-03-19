package pages

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func GetFilePath(c echo.Context) error {
	fullpath := c.Param("fullpath")

	questFilePath, err := os.ReadFile(fullpath)
	if err != nil {
		log.Println("Failed to GetFilePath")
		log.Print(err)
		return err
	}

	fileContent := string(questFilePath)

	c.Response().Header().Set("Content-Type", "application/x-yaml")
	return c.Render(200, fileContent, "")

	// c.Response().Header().Set("Content-Type", "application/x-yaml")
	// return c.String(http.StatusOK, "<html>asd</html>")
}
