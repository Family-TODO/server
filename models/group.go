package models

import "github.com/jinzhu/gorm"

type Group struct {
	gorm.Model
	CreatorId   uint
	Name        string
	Description string
	Users       []User `gorm:"many2many:user_group;"`
}
