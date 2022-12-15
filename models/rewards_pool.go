package models

const RewardsPoolTable = "rewards_pool"

type RewardPool struct {
	ID         uint64  `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Available  float64 `gorm:"column:available"      json:"available"`
	Sent       float64 `gorm:"column:sent"           json:"sent"`
	DailyLimit float64 `gorm:"column:daily_limit"    json:"daily_limit"`
}
