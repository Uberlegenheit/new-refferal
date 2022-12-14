package models

import "time"

const StakesTable = "stakes"
const FailedStakesTable = "failed_stakes"

type Stake struct {
	ID         uint64    `gorm:"column:id;PRIMARY_KEY"        json:"id"`
	UserID     uint64    `gorm:"column:user_id"               json:"user_id"`
	Amount     float64   `gorm:"column:amount"                json:"amount"`
	Status     bool      `gorm:"column:status;default:true"   json:"status"`
	TypeID     uint64    `gorm:"column:type_id;default:1"     json:"type_id"`
	BoxesGiven uint64    `gorm:"column:boxes_given;default:0" json:"boxes_given"`
	Hash       string    `gorm:"column:tx_hash"               json:"hash"`
	Created    time.Time `gorm:"column:created;default:now()" json:"created"`
}

type StakeShow struct {
	ID      uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID  uint64    `gorm:"column:user_id"        json:"user_id"`
	Amount  float64   `gorm:"column:amount"         json:"amount"`
	Status  bool      `gorm:"column:status"         json:"status"`
	Hash    string    `gorm:"column:tx_hash"        json:"hash"`
	Boxes   uint64    `gorm:"column:boxes"          json:"boxes"`
	Created time.Time `gorm:"column:created"        json:"created"`
}

type FailedStakeShow struct {
	ID         uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserID     uint64    `gorm:"column:user_id"        json:"user_id"`
	WalletName string    `gorm:"column:wallet_name"    json:"wallet_name"`
	WalletAddr string    `gorm:"column:wallet_address" json:"wallet_address"`
	Amount     float64   `gorm:"column:amount"         json:"amount"`
	TypeID     uint64    `gorm:"column:type_id"        json:"type_id"`
	Type       string    `gorm:"column:type"           json:"type"`
	Hash       string    `gorm:"column:tx_hash"        json:"hash"`
	Boxes      uint64    `gorm:"column:boxes_given"    json:"boxes_given"`
	Created    time.Time `gorm:"column:created"        json:"created"`
}

type StakeAndBoxStat struct {
	UserID     uint64  `gorm:"column:user_id"     json:"user_id"`
	TotalStake float64 `gorm:"column:total_stake" json:"total_stake"`
	TotalBoxes int64   `gorm:"column:total_boxes" json:"total_boxes"`
}

type StakeAndProgress struct {
	UserID     uint64  `gorm:"column:user_id"     json:"user_id"`
	TotalStake float64 `gorm:"column:total_stake" json:"total_stake"`
	Progress   float64 `json:"progress"`
}
