package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	checker "github.com/Looty/kubeground/cmd/web/checker"
	"github.com/Looty/kubeground/cmd/web/config"
	pages "github.com/Looty/kubeground/cmd/web/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Template is a custom renderer for echo that uses html/template
type Template struct {
	templates *template.Template
}

// Render renders the template using html/template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		return err
	}
	return nil
}

func main() {
	// Flags
	cfg := config.ParseFlags()

	// Server
	e := echo.New()
	e.Static("/static", "assets")

	// Logger
	e.Use(middleware.Logger())
	e.Logger.SetOutput(os.Stdout)
	e.Logger.Info("InClusterConnection:", cfg.InClusterConnection)
	e.Logger.Info("Port:", cfg.Port)

	// Read all HTML files in the template directory
	var templateFiles []string
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Only add files with .html extension
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		e.Logger.Fatalf("Error reading template files: %v", err)
	}

	t := &Template{templates: template.Must(template.ParseFiles(templateFiles...))}
	e.Renderer = t

	e.GET("/", pages.DisplayIndexQuest)
	e.GET("/about", pages.DisplayAboutQuest)
	e.GET("/health", pages.HealthCheck)
	e.GET("/checker", pages.DisplayIndexChecker)
	e.GET("/quest/:id", pages.DisplayQuest)
	e.POST("/quest/:id/apply", pages.QuestApplyQuest)
	e.POST("/quest/:id/delete", pages.QuestDeleteQuest)
	e.POST("/quest/:id/restart", pages.QuestRestartQuest)
	e.POST("/quest/:id/validate", checker.ValidateQuest)

	e.GET("/file/:fullpath", pages.GetFilePath)

	e.POST("/remove_alert", pages.DeleteAlert)
	e.GET("/clear_alerts", pages.ClearAlerts)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
