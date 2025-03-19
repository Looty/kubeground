package pages

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func DisplayQuest(c echo.Context) error {
	id := c.Param("id")

	questName := ""
	questLevel := ""
	instructions := ""
	hints := ""
	manifests := ""
	completed := false

	var cl client.Client
	client.Init(&cl)

	data, err := cl.ClientSet.RESTClient().
		Get().
		AbsPath("/apis/quest.looty.com/v1/quests").
		DoRaw(context.TODO())

	if err != nil {
		log.Printf("Failed to fetch quest in display: %s", data)
		log.Print(err)
		return err
	}

	// Unmarshal the JSON data into the Response struct
	var response internal.QuestList
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("Failed to Response quest in display")
		log.Print(err)
		return err
	}

	// Print the resource name
	for _, item := range response.Quests {
		if strconv.Itoa(item.Spec.Level) == id {

			questLevel = strconv.Itoa(item.Spec.Level)
			questName = item.Metadata.Name
			completed = item.Spec.Completed
			instructions = item.Spec.Instructions
			hints = item.Spec.Hints
			manifests = item.Spec.Manifests

			log.Println("Resource Name:", questName)
			log.Println("Level:", questLevel)
			log.Println("instructions:", instructions)
			log.Println("Hint:", hints)
			log.Println("Manifests:", manifests)
			log.Println("Completed:", completed)
		}
	}

	return c.Render(http.StatusOK, "quest.html", map[string]interface{}{
		"QuestName":    questName,
		"QuestLevel":   questLevel,
		"Completed":    completed,
		"Instructions": "\n" + instructions,
		"Hints":        "\n" + hints,
		"Manifests":    "\n" + manifests,
	})
}
