package mocks

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletServiceMock struct {
	mock.Mock
}

func (m *WalletServiceMock) ProcessOperation(op *model.WalletOperation) error {
	args := m.Called(op)
	return args.Error(0)
}

func (m *WalletServiceMock) GetBalance(id uuid.UUID) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}
