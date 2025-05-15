package repository

import (
	"fmt"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type WalletRepository interface {
	CreateOrUpdateWallet(op *model.WalletOperation) error
	GetBalance(walletID uuid.UUID) (int64, error)
}

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) CreateOrUpdateWallet(op *model.WalletOperation) error {
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		err := r.db.Transaction(func(tx *gorm.DB) error {
			var wallet model.Wallet
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).FirstOrCreate(&wallet, model.Wallet{ID: op.WalletID}).Error; err != nil {
				return err
			}
			if op.OperationType == model.Withdraw && wallet.Balance < op.Amount {
				return fmt.Errorf("insufficient funds")
			}
			if op.OperationType == model.Deposit {
				wallet.Balance += op.Amount
			} else if op.OperationType == model.Withdraw {
				wallet.Balance -= op.Amount
			} else {
				return fmt.Errorf("invalid operation type")
			}
			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}

			return tx.Create(op).Error
		})
		if err == nil {
			return nil
		}
		if utils.IsRetraybaleError(err) {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		return err
	}
	return fmt.Errorf("max retry attempts exceeded")
}

func (r *WalletRepo) GetBalance(walletID uuid.UUID) (int64, error) {
	var wallet model.Wallet
	if err := r.db.First(&wallet, "id = ?", walletID).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
