package models

import "time"

const PayoutsTable = "payouts"

type Payout struct {
	ID      uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID  uint64    `gorm:"column:user_id"        json:"user_id"`
	Amount  float64   `gorm:"column:amount"         json:"amount"`
	Fee     float64   `gorm:"column:fee"            json:"fee"`
	TxHash  string    `gorm:"column:tx_hash"        json:"tx_hash"`
	Created time.Time `gorm:"column:created"        json:"created"`
}

type PayoutShow struct {
	ID         uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID     uint64    `gorm:"column:user_id"        json:"user_id"`
	WalletName uint64    `gorm:"column:wallet_name"    json:"wallet_name"`
	WalletAddr uint64    `gorm:"column:wallet_address" json:"wallet_address"`
	Amount     float64   `gorm:"column:amount"         json:"amount"`
	Fee        float64   `gorm:"column:fee"            json:"fee"`
	TxHash     string    `gorm:"column:tx_hash"        json:"tx_hash"`
	Created    time.Time `gorm:"column:created"        json:"created"`
}
