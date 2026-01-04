package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/converters"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	transactionspb "github.com/radmickey/money-control/backend/proto/transactions"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// TransactionsHandler handles transactions-related requests
type TransactionsHandler struct {
	proxy *proxy.ServiceProxy
}

// NewTransactionsHandler creates a new transactions handler
func NewTransactionsHandler(sp *proxy.ServiceProxy) *TransactionsHandler {
	return &TransactionsHandler{proxy: sp}
}

// CreateTransaction creates a new transaction
func (h *TransactionsHandler) CreateTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req struct {
		AccountID      string  `json:"account_id"`
		SubAccountID   string  `json:"sub_account_id"`
		Amount         float64 `json:"amount" binding:"required"`
		Currency       string  `json:"currency"`
		Type           string  `json:"type" binding:"required"`
		Category       string  `json:"category"`
		CustomCategory string  `json:"custom_category"`
		Description    string  `json:"description"`
		Merchant       string  `json:"merchant"`
		Date           string  `json:"date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Transactions.CreateTransaction(c.Request.Context(), &transactionspb.CreateTransactionRequest{
		UserId:         userID,
		SubAccountId:   req.SubAccountID,
		Amount:         req.Amount,
		Currency:       converters.DefaultCurrency(req.Currency),
		Type:           converters.StringToTransactionType(req.Type),
		Category:       converters.StringToTransactionCategory(req.Category),
		CustomCategory: req.CustomCategory,
		Description:    req.Description,
		Merchant:       req.Merchant,
		Date:           converters.ParseDate(req.Date),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, resp)
}

// ListTransactions lists transactions
func (h *TransactionsHandler) ListTransactions(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	subAccountID := c.Query("sub_account_id")

	resp, err := h.proxy.Transactions.ListTransactions(c.Request.Context(), &transactionspb.ListTransactionsRequest{
		UserId:       userID,
		SubAccountId: subAccountID,
		Page:         1,
		PageSize:     50,
		SortBy:       "date",
		SortDesc:     true,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"transactions": resp.Transactions,
		"total":        resp.Total,
	})
}

// GetTransaction gets a transaction
func (h *TransactionsHandler) GetTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	resp, err := h.proxy.Transactions.GetTransaction(c.Request.Context(), &transactionspb.GetTransactionRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.NotFound(c, "Transaction not found")
		return
	}

	utils.Success(c, resp)
}

// UpdateTransaction updates a transaction
func (h *TransactionsHandler) UpdateTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req struct {
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Date        string  `json:"date"`
		Currency    string  `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Transactions.UpdateTransaction(c.Request.Context(), &transactionspb.UpdateTransactionRequest{
		Id:          id,
		UserId:      userID,
		Amount:      req.Amount,
		Type:        converters.StringToTransactionType(req.Type),
		Category:    converters.StringToTransactionCategory(req.Category),
		Description: req.Description,
		Date:        converters.ParseDate(req.Date),
		Currency:    converters.DefaultCurrency(req.Currency),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// DeleteTransaction deletes a transaction
func (h *TransactionsHandler) DeleteTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	_, err := h.proxy.Transactions.DeleteTransaction(c.Request.Context(), &transactionspb.DeleteTransactionRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// GetSummary gets transaction summary
func (h *TransactionsHandler) GetSummary(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	resp, err := h.proxy.Transactions.GetTransactionsSummary(c.Request.Context(), &transactionspb.GetTransactionsSummaryRequest{
		UserId:       userID,
		StartDate:    converters.ParseDate(c.Query("start_date")),
		EndDate:      converters.ParseDate(c.Query("end_date")),
		BaseCurrency: converters.DefaultCurrency(c.Query("currency")),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

