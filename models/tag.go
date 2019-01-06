package models

type Tag struct {
	Model
	OwnerId   uint   `gorm:"index; not null" json:"owner_id"`
	OwnerType string `gorm:"not null" json:"owner_type"`
	Name      string `gorm:"not null" json:"name"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
}
