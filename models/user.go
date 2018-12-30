package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name          string
	Login         string  `gorm:"unique_index; not null"`
	Password      string  `gorm:"not null"`
	IsAdmin       bool    `gorm:"default:false; not null"`
	IsDisabled    bool    `gorm:"default:false; not null"`
	Groups        []Group `gorm:"many2many:user_group;"`
	CreatorGroups []Group `gorm:"foreignkey:CreatorId"`
	Tasks         []Task
}
