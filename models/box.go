package models

type Box struct {
	ID        uint64 `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID    uint64 `gorm:"column:user_id"        json:"user_id"`
	Available uint64 `gorm:"column:available"      json:"available"`
	Opened    uint64 `gorm:"column:opened"         json:"opened"`
}
