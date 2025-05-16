package repository

import (
	"fmt"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	CreateOrUpdateWallet(op *model.WalletOperation) error
	GetBalance(walletID uuid.UUID) (int64, error)
}

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

// CreateOrUpdateWallet пополнение, снятие с кошелька
func (r *WalletRepo) CreateOrUpdateWallet(op *model.WalletOperation) error {
	switch op.OperationType {
	case model.Deposit:
		// Пополнение - создать или увеличить баланс
		err := r.db.Transaction(func(tx *gorm.DB) error {
			// Попытка создать кошелек, если его нет
			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"balance": gorm.Expr("wallets.balance + ?", op.Amount),
				}),
			}).Create(&model.Wallet{ID: op.WalletID, Balance: op.Amount}).Error
			if err != nil {
				return err
			}
			return tx.Create(op).Error
		})
		return err

	case model.Withdraw:
		// Снятие - атомарное обновление с проверкой баланса
		err := r.db.Transaction(func(tx *gorm.DB) error {
			res := tx.Exec(
				`UPDATE wallets SET balance = balance - ? WHERE id = ? AND balance >= ?`,
				op.Amount, op.WalletID, op.Amount,
			)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("insufficient funds or wallet not found")
			}
			return tx.Create(op).Error
		})
		return err

	default:
		return fmt.Errorf("invalid operation type")
	}
}

// GetBalance Баланс кошелька
func (r *WalletRepo) GetBalance(walletID uuid.UUID) (int64, error) {
	var wallet model.Wallet
	// Получить кошелек
	if err := r.db.First(&wallet, "id = ?", walletID).Error; err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
