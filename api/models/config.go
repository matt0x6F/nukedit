package models

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Key       string `json:"key" gorm:"primaryKey"`
	Value     string `json:"value"`
	Type      string `json:"type" gorm:"default:'string'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
