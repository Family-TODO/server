package models

import (
	"time"
)

type Task struct {
	Model
	UserId  uint       `gorm:"index; not null" json:"user_id"`
	GroupId uint       `gorm:"index" json:"group_id"`
	Name    string     `gorm:"not null" json:"name"`
	IsDone  bool       `gorm:"default:false; not null" json:"is_done"`
	Time    *time.Time `json:"time"`
}
