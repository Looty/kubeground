package level

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ListLevels(c echo.Context) error {
	var levels []Level

	db, err := gorm.Open(sqlite.Open("project.db"), &gorm.Config{})
	if err != nil {
		panic("couldn't open project.db")
	}

	result := db.Find(&levels)
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "There are no levels available")
	}

	return c.JSON(http.StatusOK, levels)
}
