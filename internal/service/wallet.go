package service

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/google/uuid"
)

type WalletService interface {
	ProcessOperation(op *model.WalletOperation) error
	GetBalance(id uuid.UUID) (int64, error)
}

type WalletServiceImpl struct {
	walletRepository repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &WalletServiceImpl{walletRepository: repo}
}

func (s *WalletServiceImpl) ProcessOperation(op *model.WalletOperation) error {
	return s.walletRepository.CreateOrUpdateWallet(op)
}

func (s *WalletServiceImpl) GetBalance(walletID uuid.UUID) (int64, error) {
	return s.walletRepository.GetBalance(walletID)
}
