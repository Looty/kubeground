package level

import (
	"net/http"
	"strconv"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StopLevel(c echo.Context) error {
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

	db.Model(&level).Updates(map[string]interface{}{
		"Started": "false",
	})

	db.First(&level, id)

	var cl client.Client
	client.Init(&cl)

	internal.ParseResourceFolder(level.Path, cl.Dynamic, cl.ClientSet, false)

	return c.JSON(http.StatusOK, level)
}
