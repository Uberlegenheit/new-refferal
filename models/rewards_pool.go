package models

const RewardsPoolTable = "rewards_pool"

type RewardsPool struct {
	ID        uint64  `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Available float64 `gorm:"column:available"      json:"available"`
	Sent      float64 `gorm:"column:sent"           json:"sent"`
}
