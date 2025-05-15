package repository_test

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// Подключаемся к тестовой безе данных перед каждым тестом
func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5432 user=postgres password=postgres__ dbname=wallet_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("не удалось подключиться к БД: %v", err)
	}

	// Удаляем старые таблицы
	_ = db.Migrator().DropTable(&model.Wallet{}, &model.WalletOperation{})

	// Создаем заново
	err = db.AutoMigrate(&model.Wallet{}, &model.WalletOperation{})
	if err != nil {
		t.Fatalf("ошибка миграции", err)
	}
	return db
}

// Тест 1: Успешный депозит
func TestCreateOrUpdateWallet_Deposit(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWalletRepository(db)
	walletID := uuid.New()

	op := &model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        1000,
	}

	err := repo.CreateOrUpdateWallet(op)
	assert.NoError(t, err)

	balance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), balance)
}

// Тест 2: Успешное снятие средств
func TestCreateOrUpdateWallet_Withdraw(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWalletRepository(db)
	walletID := uuid.New()

	// Сначала пополняем
	err := repo.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        1000,
	})
	assert.NoError(t, err)

	// Затем снимаем
	err = repo.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        400,
	})
	assert.NoError(t, err)

	balance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(600), balance)
}

// Тест 3: Ошибка при недостаточном балансе
func TestCreateOrUpdateWallet_InsufficientFunds(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWalletRepository(db)
	walletID := uuid.New()

	// Пополняем на 100
	err := repo.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        100,
	})
	assert.NoError(t, err)

	// Пытаемся снять 200
	err = repo.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        200,
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "insufficient funds")
}

func TestGetBalance_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewWalletRepository(db)

	_, err := repo.GetBalance(uuid.New())
	assert.Error(t, err)
	assert.EqualError(t, err, "wallet not found")
}
