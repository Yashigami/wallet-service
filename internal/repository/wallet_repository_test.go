package repository

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// Подключаемся к тестовой безе данных перед каждым тестом
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Wallet{}, &model.WalletOperation{})
	assert.NoError(t, err)

	return db
}

func TestCreateOrUpdateWallet_Deposit(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWalletRepository(db)

	walletID := uuid.New()
	op := &model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        100,
	}

	err := repo.CreateOrUpdateWallet(op)
	assert.NoError(t, err)

	var wallet model.Wallet
	err = db.First(&wallet, "id = ?", walletID).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(100), wallet.Balance)
}

func TestCreateOrUpdateWallet_Withdraw_Success(t *testing.T) {
	db := setupTestDB(t)
	walletRepository := NewWalletRepository(db)

	walletID := uuid.New()

	// Первичный депозит
	err := walletRepository.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        200,
	})
	assert.NoError(t, err)

	// Списание
	err = walletRepository.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        150,
	})
	assert.NoError(t, err)

	balance, err := walletRepository.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), balance)
}

func TestCreateOrUpdateWallet_Withdraw_InsufficientFunds(t *testing.T) {
	db := setupTestDB(t)
	walletRepository := NewWalletRepository(db)

	walletID := uuid.New()

	// Пытаемся списать с нулевым балансом
	err := walletRepository.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        50,
	})
	assert.ErrorContains(t, err, "insufficient funds")
}
