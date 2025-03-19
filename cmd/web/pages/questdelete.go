package pages

import (
	"fmt"
	"log"
	"net/http"

	quest "github.com/Looty/kubeground/cmd/web/level"
	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func QuestDeleteQuest(c echo.Context) error {
	id := c.Param("id")
	var err error
	var errMessage, alert string

	log.Printf("Deleting manifests in questID: %s", id)
	errMessage = quest.DestroyManifests(id)

	alert = fmt.Sprintf("Quest %s was deleted", id)
	internal.AddMessageToQueue("secondary", alert, errMessage, false)

	log.Println(err)
	return c.Redirect(http.StatusFound, "/")
}
