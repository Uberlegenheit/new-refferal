package models

import "time"

const UsersTable = "users"

type User struct {
	ID            uint64    `gorm:"column:id;PRIMARY_KEY" json:"id"`
	WalletName    string    `gorm:"column:wallet_name"            json:"wallet_name"`
	WalletAddress string    `gorm:"column:wallet_address"         json:"wallet_address"`
	Created       time.Time `gorm:"column:created"                json:"created"`
}
