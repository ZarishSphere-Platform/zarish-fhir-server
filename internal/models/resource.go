package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FHIRResource struct {
	ID           string         `gorm:"primaryKey"`
	ResourceType string         `gorm:"index"`
	Content      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
