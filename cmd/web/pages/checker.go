package pages

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"
	"github.com/labstack/echo/v4"
)

func DisplayIndexChecker(c echo.Context) error {
	var cl client.Client
	client.Init(&cl)

	data, err := cl.ClientSet.RESTClient().
		Get().
		AbsPath("/apis/checker.looty.com/v1/checkers").
		DoRaw(context.TODO())

	if err != nil {
		log.Printf("Failed to fetch checker in API: %s", data)
		log.Print(err)
		return err
	}

	// Unmarshal the JSON data into the Response struct
	var response internal.CheckerList
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("Failed to Unmarshal checker")
		log.Print(err)
		return err
	}

	var checkers []internal.Checker

	// Print the resource name
	for _, item := range response.Checkers {
		checkers = append(checkers, internal.Checker{
			Metadata: internal.CheckerMetadata{
				Name:      item.Metadata.Name,
				Namespace: "",
			},
			Spec: internal.CheckerSpec{
				QuestRef:   item.Spec.QuestRef,
				Validation: item.Spec.Validation,
			},
		})
		log.Println("Checker Resource Name:", item.Metadata.Name)
	}

	sort.Slice(checkers, func(i, j int) bool {
		return checkers[i].Spec.QuestRef < checkers[j].Spec.QuestRef
	})

	return c.Render(http.StatusOK, "checkers.html", map[string]interface{}{
		"Checkers": checkers,
	})
}
