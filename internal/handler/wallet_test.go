package handler_test

import (
	"bytes"
	"fmt"
	"github.com/Yashigami/wallet-service/internal/handler"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/repository"
	"github.com/Yashigami/wallet-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeWalletRepo struct {
	storage map[uuid.UUID]int64
}

var _ repository.WalletRepository = (*fakeWalletRepo)(nil)

func newFakeWalletRepo() *fakeWalletRepo {
	return &fakeWalletRepo{storage: make(map[uuid.UUID]int64)}
}

func (f *fakeWalletRepo) CreateOrUpdateWallet(op *model.WalletOperation) error {
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
		return fmt.Errorf("invalid operation type")
	}
	return nil
}

func (f *fakeWalletRepo) GetBalance(walletID uuid.UUID) (int64, error) {
	balance, ok := f.storage[walletID]
	if !ok {
		return 0, fmt.Errorf("wallet not found")
	}
	return balance, nil
}

// Вспомогательная функция для создания тестового роутера
func setupRouter() *gin.Engine {
	r := gin.Default()

	repo := newFakeWalletRepo()
	service := service.NewWalletService(repo)
	handler.NewWalletHandler(service, r)

	return r
}

// Тест 1: Создание кошелька
func TestCreateWallet(t *testing.T) {
	router := setupRouter()

	payload := []byte(`{"username":"user1", "balance":1000}`)

	req, _ := http.NewRequest("POST", "/wallets", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

}

// Тест 2: Получение баланса кошелька
func TestGetWallet(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/wallets/user1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

}

// Тест 3: Пополнение кошелька
func TestDepositWallet(t *testing.T) {
	router := setupRouter()

	payload := []byte(`{"amount":500}`)

	req, _ := http.NewRequest("POST", "/wallets/user1/deposit", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Тест 4: Списание кошелька
func TestWithdrawWallet(t *testing.T) {
	router := setupRouter()

	payload := []byte(`{"amount":300}`)
	req, _ := http.NewRequest("POST", "/wallets/user1/withdraw", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
