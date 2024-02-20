package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID              string `json:"id" gorm:"primaryKey"`
	ClientID        string `json:"client_id"`
	ClientSecret    string `json:"client_secret"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	RequirePassword bool   `json:"require_password"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
