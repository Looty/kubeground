package level

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetLevel(c echo.Context) error {
	var level Level

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Couldn't parse id")
	}

	db, err := gorm.Open(sqlite.Open("project.db"), &gorm.Config{})
	if err != nil {
		panic("couldn't open project.db")
	}

	result := db.First(&level, id)
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "No level was found")
	}

	return c.JSON(http.StatusOK, level)
}
