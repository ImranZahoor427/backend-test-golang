package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imranzahoor/banking-ledger/internal/service"
	"github.com/imranzahoor/banking-ledger/pkg/errors"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: s}
}

func (h *AccountHandler) RegisterRoutes(rg *gin.RouterGroup) {
	accounts := rg.Group("/accounts")
	accounts.POST("", h.CreateAccount)
	accounts.GET("/:id", h.GetAccount)
}

type createAccountRequest struct {
	OwnerName      string `json:"owner_name" binding:"required"`
	InitialBalance int64  `json:"initial_balance" binding:"gte=0"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), req.OwnerName, req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	account, err := h.accountService.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.ErrAccountNotFound)
		return
	}
	c.JSON(http.StatusOK, account)
}
