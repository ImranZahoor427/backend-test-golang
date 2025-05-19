package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imranzahoor/banking-ledger/internal/service"
)

type Handler struct {
	AccountHandler     *AccountHandler
	TransactionHandler *TransactionHandler
}

func NewHandler(accountSvc *service.AccountService, transactionSvc *service.TransactionService) *Handler {
	return &Handler{
		AccountHandler:     NewAccountHandler(accountSvc),
		TransactionHandler: NewTransactionHandler(transactionSvc),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	h.AccountHandler.RegisterRoutes(api)
	h.TransactionHandler.RegisterRoutes(api)
}
