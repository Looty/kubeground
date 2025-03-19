package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func DisplayIndexQuest(c echo.Context) error {
	var cl client.Client
	client.Init(&cl)

	data, err := cl.ClientSet.RESTClient().
		Get().
		AbsPath("/apis/quest.looty.com/v1/quests").
		DoRaw(context.TODO())

	if err != nil {
		log.Printf("Failed to fetch quests in index: %s", data)
		log.Print(err)
		return err
	}

	// Unmarshal the JSON data into the Response struct
	var response internal.QuestList
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("Failed to Unmarshal quest resources")
		log.Print(err)
		return err
	}

	var quests []internal.Quest
	var completedCount int

	// Print the resource name
	for _, item := range response.Quests {
		quests = append(quests, internal.Quest{
			Metadata: internal.QuestMetadata{
				Name:      item.Metadata.Name,
				Namespace: item.Metadata.Namespace,
			},
			Spec: internal.QuestSpec{
				Completed:    item.Spec.Completed,
				Hints:        item.Spec.Hints,
				Instructions: item.Spec.Instructions,
				Level:        item.Spec.Level,
				Manifests:    item.Spec.Manifests,
			},
		})

		if item.Spec.Completed {
			completedCount++
		}
	}

	sort.Slice(quests, func(i, j int) bool {
		return quests[i].Spec.Level < quests[j].Spec.Level
	})

	completion := int(float64(completedCount) / float64(len(quests)) * 100)
	progressBar := template.HTML(fmt.Sprintf(`
	<div class="progress" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="%d">
	  <div class="progress-bar progress-bar-striped bg-success progress-bar-animated" style="width: %d%%;">%d%%</div>
	</div>`, completion, len(quests), completion, completion))

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Quests":             quests,
		"Alerts":             internal.MessageQueue,
		"CompletedQuestsNum": progressBar,
	})
}
