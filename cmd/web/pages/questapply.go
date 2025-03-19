package pages

import (
	"fmt"
	"log"
	"net/http"

	quest "github.com/Looty/kubeground/cmd/web/level"
	"github.com/labstack/echo/v4"

	"github.com/Looty/kubeground/internal"
)

func QuestApplyQuest(c echo.Context) error {
	id := c.Param("id")
	var err error
	var errMessage, alert string

	log.Printf("Applying manifests in questID: %s", id)
	errMessage = quest.ApplyManifests(id)

	alert = fmt.Sprintf("Quest %s was applied", id)
	internal.AddMessageToQueue("primary", alert, errMessage, false)

	log.Println(err)
	return c.Redirect(http.StatusFound, "/")
}
