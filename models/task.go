package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	UserId  uint   `gorm:"index; not null"`
	GroupId uint   `gorm:"index"`
	Name    string `gorm:"not null"`
	IsDone  bool   `gorm:"default:false; not null"`
	Time    *time.Time
}
