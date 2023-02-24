package level

import (
	"time"

	"gorm.io/gorm"
)

type Level struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Path        string         `json:"path"`
	Started     bool           `json:"started" gorm:"default:false"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CompletedAt time.Time      `json:"completedAt"`
}
