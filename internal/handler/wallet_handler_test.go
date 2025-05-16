package handler

import (
	"bytes"
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	walletID = uuid.New()
)

func TestWalletHandler_OperateDeposit(t *testing.T) {
	router := gin.Default()
	mockService := new(mocks.WalletServiceMock)
	NewWalletHandler(mockService, router)

	reqBody := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "DEPOSIT",
		"amount": 1000
	}`

	mockService.
		On("ProcessOperation", mock.MatchedBy(func(op *model.WalletOperation) bool {
			return op.WalletID == walletID && op.Amount == 1000 && op.OperationType == model.Deposit
		})).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "operation successful")
	mockService.AssertExpectations(t)
}

func TestWalletHandler_OperateWithdraw(t *testing.T) {
	router := gin.Default()
	mockService := new(mocks.WalletServiceMock)
	NewWalletHandler(mockService, router)

	reqBody := `{
		"walletId": "` + walletID.String() + `",
		"operationType": "WITHDRAW",
		"amount": 1000
	}`

	mockService.
		On("ProcessOperation", mock.MatchedBy(func(op *model.WalletOperation) bool {
			return op.WalletID == walletID && op.Amount == 1000 && op.OperationType == model.Withdraw
		})).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "operation successful")
	mockService.AssertExpectations(t)
}

func TestWalletHandler_GetBalance(t *testing.T) {
	router := gin.Default()
	mockService := new(mocks.WalletServiceMock)
	NewWalletHandler(mockService, router)

	mockService.
		On("GetBalance", walletID).
		Return(int64(5000), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "5000")
	mockService.AssertExpectations(t)
}
