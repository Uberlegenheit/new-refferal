package models

import "time"

type TotalStats struct {
	StakeSum        float64 `gorm:"column:total_paid"       json:"total_paid"`
	RedelegationSum float64 `gorm:"column:box_paid"         json:"box_paid"`
	InvitedSum      float64 `gorm:"column:delegators_count" json:"delegators_count"`
	BoxesGiven      uint64  `gorm:"column:boxes_given"      json:"boxes_given"`
	BoxesOpened     uint64  `gorm:"column:boxes_opened"     json:"boxes_opened"`
	BoxesRewards    float64 `gorm:"column:boxes_rewards"    json:"boxes_rewards"`
	BoxesUnpaid     float64 `gorm:"column:boxes_unpaid"     json:"boxes_unpaid"`
}

type RewardPaymentsStats struct {
	ReferrerID uint64    `gorm:"column:user_id"     json:"referrer_id"`
	ReferralID uint64    `gorm:"column:user_id"     json:"referral_id"`
	WalletName string    `gorm:"column:wallet_name" json:"wallet_name"`
	TypeID     uint64    `gorm:"column:type_id"     json:"type_id"`
	Type       string    `gorm:"column:type"        json:"type"`
	Amount     float64   `gorm:"column:amount"      json:"amount"`
	Boxes      uint64    `gorm:"column:boxes"       json:"boxes"`
	Created    time.Time `gorm:"column:created"     json:"created"`
	TxHash     float64   `gorm:"column:tx_hash"     json:"tx_hash"`
}

type TotalStakeStats struct {
	UserID     uint64    `gorm:"column:user_id"     json:"user_id"`
	WalletName string    `gorm:"column:wallet_name" json:"wallet_name"`
	TypeID     uint64    `gorm:"column:type_id"     json:"type_id"`
	Type       string    `gorm:"column:type"        json:"type"`
	Amount     float64   `gorm:"column:amount"      json:"amount"`
	Boxes      uint64    `gorm:"column:boxes"       json:"boxes"`
	Created    time.Time `gorm:"column:created"     json:"created"`
	TxHash     float64   `gorm:"column:tx_hash"     json:"tx_hash"`
}

type FriendStakeStats struct {
	ReferrerID uint64    `gorm:"column:user_id"     json:"referrer_id"`
	ReferralID uint64    `gorm:"column:user_id"     json:"referral_id"`
	WalletName string    `gorm:"column:wallet_name" json:"wallet_name"`
	TypeID     uint64    `gorm:"column:type_id"     json:"type_id"`
	Type       string    `gorm:"column:type"        json:"type"`
	Amount     float64   `gorm:"column:amount"      json:"amount"`
	Boxes      uint64    `gorm:"column:boxes"       json:"boxes"`
	Created    time.Time `gorm:"column:created"     json:"created"`
	TxHash     float64   `gorm:"column:tx_hash"     json:"tx_hash"`
}
