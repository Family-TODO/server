package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name       string
	Login      string `gorm:"unique_index; not null"`
	Password   string `gorm:"not null"`
	IsAdmin    bool   `gorm:"default:0; not null"`
	IsDisabled bool   `gorm:"default:0; not null"`
}
