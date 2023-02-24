package database

import (
	"github.com/Looty/kubeground/cmd/web/level"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) {
	db.AutoMigrate(&level.Level{})
}
