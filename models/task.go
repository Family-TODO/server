package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Task struct {
	gorm.Model
	UserId  uint   `gorm:"index;not null"`
	GroupId uint   `gorm:"index"`
	Name    string `gorm:"not null"`
	IsDone  bool   `gorm:"default:false; not null"`
	Time    *time.Time
}
