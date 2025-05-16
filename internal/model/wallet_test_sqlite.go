//go:build test
// +build test

// Структуры используются при тестировании

package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type WalletOperation struct {
	ID            uuid.UUID     `gorm:"type:char(36);primary_key"`
	WalletID      uuid.UUID     `gorm:"type:char(36);not null;index"`
	OperationType OperationType `gorm:"type:varchar(10);not null"`
	Amount        int64         `gorm:"not null"`
}

func (op *WalletOperation) BeforeCreate(tx *gorm.DB) (err error) {
	if op.ID == uuid.Nil {
		op.ID = uuid.New()
	}
	return
}

type Wallet struct {
	ID      uuid.UUID `gorm:"type:char(36);primary_key"`
	Balance int64     `gorm:"not null"`
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

func (Wallet) TableName() string {
	return "wallets"
}
