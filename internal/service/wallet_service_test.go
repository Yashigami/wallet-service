package service

import (
	"errors"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	walletID = uuid.New()
)

// Тестирование успешной операции пополнения
func TestProcessOperation_DepositSuccess(t *testing.T) {
	walletRepository := new(mocks.WalletRepositoryMock)
	walletService := NewWalletService(walletRepository)

	walletOperation := &model.WalletOperation{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        100,
	}

	walletRepository.On("CreateOrUpdateWallet", walletOperation).Return(nil)

	err := walletService.ProcessOperation(walletOperation)

	assert.NoError(t, err)
	walletRepository.AssertExpectations(t)
}

// Тестирование провальной операции списания
func TestProcessOperation_WithdrawError(t *testing.T) {
	walletRepository := new(mocks.WalletRepositoryMock)
	walletService := NewWalletService(walletRepository)

	op := &model.WalletOperation{
		WalletID:      uuid.New(),
		OperationType: model.Withdraw,
		Amount:        50,
	}

	walletRepository.On("CreateOrUpdateWallet", op).Return(errors.New("db error"))

	err := walletService.ProcessOperation(op)

	assert.EqualError(t, err, "db error")
	walletRepository.AssertExpectations(t)
}

// Тестирование успешной операции получения кошелька
func TestGetBalance_Success(t *testing.T) {
	repo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(repo)

	repo.On("GetBalance", walletID).Return(int64(500), nil)

	balance, err := svc.GetBalance(walletID)

	assert.NoError(t, err)
	assert.Equal(t, int64(500), balance)
	repo.AssertExpectations(t)
}

// Тестирование провальной операции получения кошелька
func TestGetBalance_NotFound(t *testing.T) {
	repo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(repo)

	walletID := uuid.New()

	repo.On("GetBalance", walletID).Return(int64(0), errors.New("wallet not found"))

	balance, err := svc.GetBalance(walletID)

	assert.EqualError(t, err, "wallet not found")
	assert.Equal(t, int64(0), balance)
	repo.AssertExpectations(t)
}
