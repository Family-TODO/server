package models

import (
	"../config"

	"github.com/jinzhu/gorm"
)

type Group struct {
	Model
	CreatorID   uint   `gorm:"not null" json:"creator_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	Creator User   `json:"creator"`
	Users   []User `gorm:"many2many:group_user" json:"users"`
	Tasks   []Task `json:"tasks"`
}

func GetAllGroups(groups *[]Group, userId uint, count, offset int) {
	db := config.GetDb()

	db.
		Select("DISTINCT groups.*").
		Preload("Creator").
		Preload("Users").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("DISTINCT tasks.*").
				Group("tasks.group_id").
				Order("tasks.updated_at desc").
				Preload("User")
		}).
		Joins("LEFT JOIN group_user gu ON gu.group_id = groups.id").
		Where("groups.creator_id = ? OR gu.user_id = ?", userId, userId).
		Order("groups.updated_at desc").
		Limit(count).
		Offset(offset).
		Find(&groups)
}
