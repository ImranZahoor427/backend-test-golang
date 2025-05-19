package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/service"
	"github.com/imranzahoor/banking-ledger/pkg/constants"
	"github.com/imranzahoor/banking-ledger/pkg/errors"
	"github.com/imranzahoor/banking-ledger/pkg/utils"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: s}
}

func (h *TransactionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	txns := rg.Group("/transactions")
	txns.POST("", h.CreateTransaction)

	txns.GET("/account/:id", h.GetTransactionHistory)

}

type createTransactionRequest struct {
	AccountID string `json:"account_id" binding:"required,uuid4"`
	Type      string `json:"type" binding:"required,oneof=deposit withdrawal"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId, err := utils.ParseUUID(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidAmount)
		return
	}

	if req.Type != string(constants.Deposit) && req.Type != string(constants.Withdrawal) {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidTransactionType)
		return
	}

	if req.Type == string(constants.Withdrawal) {
		hasFunds, err := h.transactionService.HasSufficientFunds(c.Request.Context(), req.AccountID, req.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasFunds {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
			return
		}
	}

	txn := &model.Transaction{
		AccountID: accountId,
		Type:      constants.TransactionType(req.Type),
		Amount:    req.Amount,
	}

	err = h.transactionService.EnqueueTransaction(c.Request.Context(), txn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":        "transaction accepted",
		"transaction_id": txn.ID,
	})
}

func (h *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	accountID := c.Param("id")

	// Parse limit and offset from query params, with defaults
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidLimit)
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, errors.ErrInvalidOffset)
		return
	}

	parsedID, err := utils.ParseUUID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrParsingID)
		return
	}

	txns, err := h.transactionService.GetTransactions(parsedID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrFetchingTransactions)
		return
	}

	c.JSON(http.StatusOK, txns)
}
