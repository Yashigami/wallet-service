package service

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/google/uuid"
)

type WalletService struct {
	repo *repository.WalletRepo
}

func NewWalletService(repo *repository.WalletRepo) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) ProcessOperation(walletID uuid.UUID, opType model.OperationType, amount int64) error {
	return s.repo.CreateOrUpdateWallet(&model.WalletOperation{
		WalletID:      walletID,
		OperationType: opType,
		Amount:        amount,
	})
}
func (s *WalletService) GetBalance(walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(walletID)
}
