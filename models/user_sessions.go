package models

import (
	"time"
)

type UserSession struct {
	UserID    uint   `gorm:"primary_key"`
	Token     string `gorm:"primary_key"`
	Ip        string
	Device    string
	Os        string
	CreatedAt time.Time
	UpdatedAt time.Time
}
