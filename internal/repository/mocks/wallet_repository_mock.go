package mocks

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// WalletRepositoryMock Заглушки для репозитория
type WalletRepositoryMock struct {
	mock.Mock
}

func (m *WalletRepositoryMock) CreateOrUpdateWallet(op *model.WalletOperation) error {
	args := m.Called(op)
	return args.Error(0)
}

func (m *WalletRepositoryMock) GetBalance(walletID uuid.UUID) (int64, error) {
	args := m.Called(walletID)
	return args.Get(0).(int64), args.Error(1)
}
