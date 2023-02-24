package database

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/Looty/kubeground/cmd/web/level"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbFileName string) (db *gorm.DB) {
	dbConnection, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return dbConnection
}

func InitLevels(db *gorm.DB, resourcesPath string) {
	path := resourcesPath

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("No levels were found")
	}

	for _, file := range fileInfos {
		if err != nil {
			log.Fatal("No levels were found #2")
		}

		db.Create(&level.Level{
			Path:        filepath.Join(path, file.Name()),
			CompletedAt: time.Time{},
		})
	}
}
