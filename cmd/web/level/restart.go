package quest

import (
	"log"
)

func RecreateManifests(id string) string {
	var errorMessage string

	errorMessage = DestroyManifests(id)
	log.Printf("Deleting manifests in questID: %s", id)
	log.Print(errorMessage)

	errorMessage = ApplyManifests(id)
	log.Printf("Applying manifests in questID: %s", id)
	log.Print(errorMessage)

	return errorMessage
}
