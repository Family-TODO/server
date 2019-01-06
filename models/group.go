package models

type Group struct {
	Model
	CreatorId   uint   `gorm:"not null" json:"creator_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	Users []User `gorm:"many2many:user_group" json:"users"`
	Tasks []Task `json:"tasks"`
	Tag   Tag    `gorm:"polymorphic:Owner" json:"tag"`
}
