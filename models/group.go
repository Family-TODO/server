package models

import "github.com/jinzhu/gorm"

type Group struct {
	gorm.Model
	CreatorId   uint   `gorm:"index"`
	Name        string `gorm:"not null"`
	Description string

	Users []User `gorm:"many2many:user_group"`
	Tasks []Task
	Tag   Tag `gorm:"polymorphic:Owner"`
}
