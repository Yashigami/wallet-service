package service_test

import (
	"fmt"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Фейковая реализация репозитория
type fakeRepo struct {
	storage map[uuid.UUID]int64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{storage: make(map[uuid.UUID]int64)}
}

// Реализация метода CreateOrUpdateWallet
func (f *fakeRepo) CreateOrUpdateWallet(op *model.WalletOperation) error {
	balance := f.storage[op.WalletID]

	switch op.OperationType {
	case model.Deposit:
		f.storage[op.WalletID] = balance + op.Amount
	case model.Withdraw:
		if balance < op.Amount {
			return fmt.Errorf("insufficient funds")
		}
		f.storage[op.WalletID] = balance - op.Amount
	default:
		return fmt.Errorf("invalid operation")
	}
	return nil
}

// Реализация метода GetBalance
func (f *fakeRepo) GetBalance(walletID uuid.UUID) (int64, error) {
	balance, ok := f.storage[walletID]
	if !ok {
		return 0, fmt.Errorf("wallet not found")
	}
	return balance, nil
}

func TestDepositOperation(t *testing.T) {
	// Создаём фейковый репозиторий и сервис
	repo := newFakeRepo()
	svc := service.NewWalletService(repo)

	walletID := uuid.New()

	// Выполняем депозит
	err := svc.ProcessOperation(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        1000,
	})

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что баланс стал 1000
	balance, _ := repo.GetBalance(walletID)
	assert.Equal(t, int64(1000), balance)
}

func TestWithdrawOperation(t *testing.T) {
	repo := newFakeRepo()
	svc := service.NewWalletService(repo)

	walletID := uuid.New()
	// Сначала закинем денег
	repo.storage[walletID] = 1000

	// Снимаем 400
	err := svc.ProcessOperation(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        400,
	})

	assert.NoError(t, err)
	balance, _ := repo.GetBalance(walletID)
	assert.Equal(t, int64(600), balance)
}

func TestWithdrawInsufficientFunds(t *testing.T) {
	repo := newFakeRepo()
	svc := service.NewWalletService(repo)

	walletID := uuid.New()
	repo.storage[walletID] = 100

	// Пытаемся снять больше, чем есть
	err := svc.ProcessOperation(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Withdraw,
		Amount:        200,
	})

	// Ожидаем ошибку
	assert.Error(t, err)
	assert.EqualError(t, err, "insufficient funds")
}

func TestInvalidOperation(t *testing.T) {
	repo := newFakeRepo()
	svc := service.NewWalletService(repo)

	walletID := uuid.New()

	// Передаём некорректный тип операции
	err := svc.ProcessOperation(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: "INVALID", // неправильный тип
		Amount:        500,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid operation")
}
