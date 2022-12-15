package models

const InvitationsTable = "invitations"

type Invitation struct {
	ReferrerID uint64 `gorm:"column:referrer_id" json:"referrer_id"`
	ReferralID uint64 `gorm:"column:referral_id" json:"referral_id"`
}
