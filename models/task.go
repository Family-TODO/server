package models

import (
	"time"
)

type Task struct {
	Model
	UserID   uint       `gorm:"index; not null" json:"user_id"`
	GroupID  uint       `gorm:"index; not null" json:"group_id"`
	Name     string     `gorm:"not null" json:"name"`
	IsDone   bool       `gorm:"default:false; not null" json:"is_done"`
	DoneTime *time.Time `json:"done_time"`

	User  User  `json:"user"`
	Group Group `json:"group"`
}
