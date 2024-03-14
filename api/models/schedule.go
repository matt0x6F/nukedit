package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID                    int    `json:"id" gorm:"primaryKey"`
	Username              string `json:"username"`
	CronExpression        string `json:"cronExpression" gorm:"column:cron_expression"`
	Posts                 bool   `json:"posts"`
	Comments              bool   `json:"comments"`
	MaxAge                int    `json:"maxAge"`
	UseMaxAge             bool   `json:"useMaxAge"`
	MinScore              int    `json:"minScore"`
	UseMinScore           bool   `json:"useMinScore"`
	ReplacementTextLength int    `json:"replacementTextLength"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}
