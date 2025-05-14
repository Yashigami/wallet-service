package handler

import (
	"github.com/Yashigami/wallet-service/internal/model"
	"github.com/Yashigami/wallet-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(s *service.WalletService, serverEngine *gin.Engine) *WalletHandler {
	var h = &WalletHandler{service: s}
	serverEngine.POST("/api/v1/wallet", h.operate)
	serverEngine.GET("/api/v1/wallets/:id", h.getBalance)
	return h
}

type OperationRequest struct {
	WalletID      uuid.UUID           `json:"walletId"`
	OperationType model.OperationType `json:"operationType"`
	Amount        int64               `json:"amount"`
}

func (h *WalletHandler) operate(c *gin.Context) {
	var req OperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.service.ProcessOperation(req.WalletID, req.OperationType, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "err.Error()"})
		return
	}

	c.Status(http.StatusOK)
}
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
