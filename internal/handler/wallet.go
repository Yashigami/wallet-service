package handler

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type WalletHandler struct {
	service service.WalletService
}

// NewWalletHandler Инициализация слушателей
func NewWalletHandler(s service.WalletService, serverEngine *gin.Engine) *WalletHandler {
	var h = &WalletHandler{service: s}
	serverEngine.POST("/api/v1/wallet", h.operate)
	serverEngine.GET("/api/v1/wallets/:id", h.getBalance)
	return h
}

// OperationRequest Структура запроса для выполнения операций с кошельком
type OperationRequest struct {
	WalletID      uuid.UUID           `json:"walletId"`
	OperationType model.OperationType `json:"operationType"`
	Amount        int64               `json:"amount"`
}

// Обработка запроса для выполнения операций с кошельком
func (h *WalletHandler) operate(c *gin.Context) {
	var req OperationRequest
	// Вызов бизнес-логики для выполнения операции
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// Формирование бизнес-структуры
	op := &model.WalletOperation{
		WalletID:      req.WalletID,
		OperationType: req.OperationType,
		Amount:        req.Amount,
	}

	// Вызов бизнес-логики для выполнения операции
	err := h.service.ProcessOperation(op)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "operation successful"})

}

// Обработка запроса на получение текущего баланса по кошельку
func (h *WalletHandler) getBalance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	balance, err := h.service.GetBalance(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
