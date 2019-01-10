package models

type Group struct {
	Model
	CreatorID   uint   `gorm:"not null" json:"creator_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	Creator User   `json:"creator"`
	Users   []User `gorm:"many2many:group_user" json:"users"`
	Tasks   []Task `json:"tasks"`
}
