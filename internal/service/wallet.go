package service

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) ProcessOperation(op *model.WalletOperation) error {
	return s.repo.CreateOrUpdateWallet(op)
}

func (s *WalletService) GetBalance(walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(walletID)
}
