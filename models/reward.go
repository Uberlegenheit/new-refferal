package models

import "time"

const RewardsTable = "rewards"
const CreatedRewardStatus = "created"
const SentRewardStatus = "sent"
const APYRewardType = 1
const BoxRewardType = 2

type Reward struct {
	ID      uint64    `gorm:"column:id;PRIMARY_KEY"        json:"id"`
	UserID  uint64    `gorm:"column:user_id"               json:"user_id"`
	Status  string    `gorm:"column:status"                json:"status"`
	TypeID  uint64    `gorm:"column:type_id"               json:"type_id"`
	Amount  float64   `gorm:"column:amount"                json:"amount"`
	Hash    string    `gorm:"column:tx_hash;default:"      json:"hash"`
	Created time.Time `gorm:"column:created;default:now()" json:"created"`
}

type RewardShow struct {
	ID         uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	WalletName string    `gorm:"column:wallet_name"    json:"wallet_name"`
	Status     string    `gorm:"column:status"         json:"status"`
	Type       string    `gorm:"column:type"           json:"type"`
	Amount     float64   `gorm:"column:amount"         json:"amount"`
	Hash       string    `gorm:"column:tx_hash"        json:"hash"`
	Created    time.Time `gorm:"column:created"        json:"created"`
}

type TotalReward struct {
	TotalPaid       float64 `gorm:"column:total_paid"       json:"total_paid"`
	BoxPaid         float64 `gorm:"column:box_paid"         json:"box_paid"`
	DelegatorsCount uint64  `gorm:"column:delegators_count" json:"delegators_count"`
}

type InvitationsStats struct {
	UserID        uint64  `gorm:"column:user_id"        json:"user_id"`
	WalletName    string  `gorm:"column:wallet_name"    json:"wallet_name"`
	WalletAddr    string  `gorm:"column:wallet_address" json:"wallet_address"`
	TotalReward   float64 `gorm:"column:total_reward"   json:"total_reward"`
	FriendInvited uint64  `gorm:"column:friend_invited" json:"friend_invited"`
}
