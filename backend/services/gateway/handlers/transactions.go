package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	transactionspb "github.com/radmickey/money-control/backend/proto/transactions"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	var date *timestamppb.Timestamp
	if req.Date != "" {
		t, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			t, _ = time.Parse("2006-01-02", req.Date)
		}
		date = timestamppb.New(t)
	}

	resp, err := h.proxy.Transactions.CreateTransaction(c.Request.Context(), &transactionspb.CreateTransactionRequest{
		UserId:         userID,
		SubAccountId:   req.SubAccountID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Type:           stringToTransactionType(req.Type),
		Category:       stringToTransactionCategory(req.Category),
		CustomCategory: req.CustomCategory,
		Description:    req.Description,
		Merchant:       req.Merchant,
		Date:           date,
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
		Category    string  `json:"category"`
		Description string  `json:"description"`
		Date        string  `json:"date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var date *timestamppb.Timestamp
	if req.Date != "" {
		t, _ := time.Parse("2006-01-02", req.Date)
		date = timestamppb.New(t)
	}

	resp, err := h.proxy.Transactions.UpdateTransaction(c.Request.Context(), &transactionspb.UpdateTransactionRequest{
		Id:          id,
		UserId:      userID,
		Amount:      req.Amount,
		Category:    stringToTransactionCategory(req.Category),
		Description: req.Description,
		Date:        date,
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
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	baseCurrency := c.Query("currency")

	var startDate, endDate *timestamppb.Timestamp
	if startDateStr != "" {
		t, _ := time.Parse("2006-01-02", startDateStr)
		startDate = timestamppb.New(t)
	}
	if endDateStr != "" {
		t, _ := time.Parse("2006-01-02", endDateStr)
		endDate = timestamppb.New(t)
	}

	resp, err := h.proxy.Transactions.GetTransactionsSummary(c.Request.Context(), &transactionspb.GetTransactionsSummaryRequest{
		UserId:       userID,
		StartDate:    startDate,
		EndDate:      endDate,
		BaseCurrency: baseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Helper functions
func stringToTransactionType(s string) transactionspb.TransactionType {
	switch s {
	case "INCOME", "income":
		return transactionspb.TransactionType_TRANSACTION_TYPE_INCOME
	case "EXPENSE", "expense":
		return transactionspb.TransactionType_TRANSACTION_TYPE_EXPENSE
	case "TRANSFER", "transfer":
		return transactionspb.TransactionType_TRANSACTION_TYPE_TRANSFER
	default:
		return transactionspb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}
}

func stringToTransactionCategory(s string) transactionspb.TransactionCategory {
	switch s {
	case "SALARY", "salary":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_SALARY
	case "FOOD", "food":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_FOOD
	case "TRANSPORT", "transport":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_TRANSPORT
	case "SHOPPING", "shopping":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_SHOPPING
	case "ENTERTAINMENT", "entertainment":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_ENTERTAINMENT
	case "UTILITIES", "utilities":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_UTILITIES
	case "HEALTHCARE", "healthcare":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_HEALTHCARE
	case "TRANSFER", "transfer":
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_TRANSFER
	default:
		return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_OTHER
	}
}
