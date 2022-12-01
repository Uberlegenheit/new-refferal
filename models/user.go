package models

import "time"

const UsersTable = "users"

type User struct {
	ID            uint64    `gorm:"column:id;PRIMARY_KEY"        json:"id"`
	WalletName    string    `gorm:"column:wallet_name"           json:"wallet_name"`
	WalletAddress string    `gorm:"column:wallet_address"        json:"wallet_address"`
	Role          string    `gorm:"column:role;default:user"     json:"role"`
	Created       time.Time `gorm:"column:created;default:now()" json:"created"`
}
