package models

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model
	OwnerId   uint   `gorm:"index;not null"`
	OwnerType string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Icon      string
	Color     string
}
