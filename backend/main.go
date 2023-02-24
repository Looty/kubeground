package main

import (
	"flag"
	"fmt"

	"github.com/Looty/kubeground/cmd/web/level"
	"github.com/Looty/kubeground/database"
	"github.com/Looty/kubeground/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	port            int
	dbFileName      string
	resourcesFolder string
)

func main() {
	flag.IntVar(&port, "port", 4000, "Port of the webserver.")
	flag.StringVar(&dbFileName, "dbfilename", "project.db", "SQLite DB file name.")
	flag.StringVar(&resourcesFolder, "resourcefolder", "./resources", "Resource folder directory.")
	flag.Parse()

	server := server.Config{
		Port:          port,
		DBFile:        dbFileName,
		ResourcesPath: resourcesFolder,
	}

	DB := database.InitDB(server.DBFile)
	database.Setup(DB)
	database.InitLevels(DB, server.ResourcesPath)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}",` +
			`"method":"${method}",` +
			`"uri":"${uri}",` +
			`"status":"${status}",` +
			`"error":"${error}",` +
			`"latency_human":"${latency_human}"}` +
			"\n",
	}))

	e.GET("/", level.ListLevels)
	e.GET("/healthz", level.HealthCheck)

	l := e.Group("/level")

	l.GET("/:id", level.GetLevel)
	l.GET("/validate/:id", level.ValidateLevel)
	l.POST("/start/:id", level.StartLevel)
	l.POST("/stop/:id", level.StopLevel)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", server.Port)))
}
