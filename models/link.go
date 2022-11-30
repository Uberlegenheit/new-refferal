package models

const LinksTable = "links"

type Link struct {
	ID     uint64 `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID uint64 `gorm:"column:user_id"                json:"user_id"`
	Code   string `gorm:"column:code"                   json:"code"`
}
