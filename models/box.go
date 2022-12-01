package models

type Box struct {
	ID        uint64 `gorm:"column:id;PRIMARY_KEY"      json:"id"`
	UserID    uint64 `gorm:"column:user_id"             json:"user_id"`
	Available uint64 `gorm:"column:available;default:0" json:"available"`
	Opened    uint64 `gorm:"column:opened;default:0"    json:"opened"`
}
