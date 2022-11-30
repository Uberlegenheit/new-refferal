package models

import "time"

const RewardsTable = "rewards"

type Reward struct {
	ID      uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID  uint64    `gorm:"column:user_id"        json:"user_id"`
	Status  string    `gorm:"column:status"         json:"status"`
	TypeID  uint64    `gorm:"column:type_id"        json:"type_id"`
	Amount  string    `gorm:"column:amount"         json:"amount"`
	Hash    string    `gorm:"column:tx_hash"        json:"hash"`
	Created time.Time `gorm:"column:created"        json:"created"`
}

type RewardShow struct {
	ID         uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	WalletName string    `gorm:"column:wallet_name"    json:"wallet_name"`
	Status     string    `gorm:"column:status"         json:"status"`
	Type       string    `gorm:"column:type"           json:"type"`
	Amount     string    `gorm:"column:amount"         json:"amount"`
	Hash       string    `gorm:"column:tx_hash"        json:"hash"`
	Created    time.Time `gorm:"column:created"        json:"created"`
}
