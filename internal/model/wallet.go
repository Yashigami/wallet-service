//go:build !test
// +build !test

package model

import (
	"github.com/google/uuid"
)

// Структуры используются при запуске

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type WalletOperation struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	WalletID      uuid.UUID     `gorm:"type:uuid;not null;index"`
	OperationType OperationType `gorm:"type:varchar(10); not null"`
	Amount        int64         `gorm:"not null"`
}

type Wallet struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Balance int64     `gorm:"not null"`
}

func (Wallet) TableName() string {
	return "wallets"
}
