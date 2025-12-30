package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/transactions/models"
	"github.com/radmickey/money-control/backend/services/transactions/repository"
	"github.com/radmickey/money-control/backend/services/transactions/service"
)

// RegisterHTTPRoutes registers HTTP routes for transactions
func RegisterHTTPRoutes(r *gin.RouterGroup, txService *service.TransactionService) {
	h := &HTTPHandler{txService: txService}

	transactions := r.Group("/transactions")
	{
		transactions.POST("", h.CreateTransaction)
		transactions.GET("", h.ListTransactions)
		transactions.GET("/summary", h.GetTransactionsSummary)
		transactions.GET("/:id", h.GetTransaction)
		transactions.PUT("/:id", h.UpdateTransaction)
		transactions.DELETE("/:id", h.DeleteTransaction)
		transactions.PATCH("/:id/category", h.CategorizeTransaction)
		transactions.POST("/bulk-categorize", h.BulkCategorize)
	}

	rules := r.Group("/category-rules")
	{
		rules.POST("", h.CreateCategoryRule)
		rules.GET("", h.GetCategoryRules)
		rules.DELETE("/:id", h.DeleteCategoryRule)
	}
}

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	txService *service.TransactionService
}

// CreateTransactionRequest represents create transaction request
type CreateTransactionRequest struct {
	SubAccountID         string  `json:"sub_account_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required"`
	Currency             string  `json:"currency"`
	Type                 string  `json:"type" binding:"required"`
	Category             string  `json:"category"`
	CustomCategory       string  `json:"custom_category"`
	Description          string  `json:"description"`
	Merchant             string  `json:"merchant"`
	Date                 string  `json:"date"`
	TransferToSubAccount string  `json:"transfer_to_sub_account_id"`
}

// CreateTransaction creates a new transaction
func (h *HTTPHandler) CreateTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var date time.Time
	if req.Date != "" {
		var err error
		date, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			date, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				utils.BadRequest(c, "Invalid date format")
				return
			}
		}
	}

	var subAccountID *string
	if req.SubAccountID != "" {
		subAccountID = &req.SubAccountID
	}

	var transferToSubAccount *string
	if req.TransferToSubAccount != "" {
		transferToSubAccount = &req.TransferToSubAccount
	}

	input := service.CreateTransactionInput{
		UserID:               userID,
		SubAccountID:         subAccountID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Type:                 models.TransactionType(req.Type),
		Category:             models.TransactionCategory(req.Category),
		CustomCategory:       req.CustomCategory,
		Description:          req.Description,
		Merchant:             req.Merchant,
		Date:                 date,
		TransferToSubAccount: transferToSubAccount,
	}

	tx, err := h.txService.CreateTransaction(c.Request.Context(), input)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, tx)
}

// ListTransactions lists transactions
func (h *HTTPHandler) ListTransactions(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	subAccountID := c.Query("sub_account_id")
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)
	sortBy := c.Query("sort_by")
	sortDesc := c.Query("sort_desc") == "true"

	// Date range filters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var transactions []models.Transaction
	var total int64
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		transactions, total, err = h.txService.ListTransactionsByDateRange(
			c.Request.Context(), userID, startDate, endDate, subAccountID, page, pageSize,
		)
	} else {
		transactions, total, err = h.txService.ListTransactions(
			c.Request.Context(), userID, subAccountID, page, pageSize, sortBy, sortDesc,
		)
	}

	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMeta(c, transactions, utils.PaginationMeta(page, pageSize, int(total)))
}

// GetTransaction gets a transaction by ID
func (h *HTTPHandler) GetTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	tx, err := h.txService.GetTransaction(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			utils.NotFound(c, "Transaction not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, tx)
}

// UpdateTransactionRequest represents update transaction request
type UpdateTransactionRequest struct {
	Amount         float64 `json:"amount"`
	Category       string  `json:"category"`
	CustomCategory string  `json:"custom_category"`
	Description    string  `json:"description"`
	Merchant       string  `json:"merchant"`
	Date           string  `json:"date"`
}

// UpdateTransaction updates a transaction
func (h *HTTPHandler) UpdateTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var date time.Time
	if req.Date != "" {
		var err error
		date, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			date, _ = time.Parse("2006-01-02", req.Date)
		}
	}

	input := service.UpdateTransactionInput{
		ID:             id,
		UserID:         userID,
		Amount:         req.Amount,
		Category:       models.TransactionCategory(req.Category),
		CustomCategory: req.CustomCategory,
		Description:    req.Description,
		Merchant:       req.Merchant,
		Date:           date,
	}

	tx, err := h.txService.UpdateTransaction(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			utils.NotFound(c, "Transaction not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, tx)
}

// DeleteTransaction deletes a transaction
func (h *HTTPHandler) DeleteTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	err := h.txService.DeleteTransaction(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			utils.NotFound(c, "Transaction not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// GetTransactionsSummary gets transaction summary
func (h *HTTPHandler) GetTransactionsSummary(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	summary, err := h.txService.GetTransactionsSummary(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, summary)
}

// CategorizeRequest represents categorize request
type CategorizeRequest struct {
	Category       string `json:"category" binding:"required"`
	CustomCategory string `json:"custom_category"`
}

// CategorizeTransaction categorizes a transaction
func (h *HTTPHandler) CategorizeTransaction(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req CategorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx, err := h.txService.CategorizeTransaction(
		c.Request.Context(), id, userID,
		models.TransactionCategory(req.Category), req.CustomCategory,
	)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			utils.NotFound(c, "Transaction not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, tx)
}

// BulkCategorizeRequest represents bulk categorize request
type BulkCategorizeRequest struct {
	TransactionIDs []string `json:"transaction_ids" binding:"required"`
	Category       string   `json:"category" binding:"required"`
}

// BulkCategorize categorizes multiple transactions
func (h *HTTPHandler) BulkCategorize(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req BulkCategorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updated, failed, err := h.txService.BulkCategorize(
		c.Request.Context(), userID, req.TransactionIDs,
		models.TransactionCategory(req.Category),
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"updated_count": updated,
		"failed_ids":    failed,
	})
}

// CreateCategoryRuleRequest represents create category rule request
type CreateCategoryRuleRequest struct {
	Pattern  string `json:"pattern" binding:"required"`
	Category string `json:"category" binding:"required"`
	Priority int    `json:"priority"`
}

// CreateCategoryRule creates a category rule
func (h *HTTPHandler) CreateCategoryRule(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req CreateCategoryRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	rule, err := h.txService.CreateCategoryRule(
		c.Request.Context(), userID, req.Pattern,
		models.TransactionCategory(req.Category), req.Priority,
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, rule)
}

// GetCategoryRules gets category rules
func (h *HTTPHandler) GetCategoryRules(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	rules, err := h.txService.GetCategoryRules(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, rules)
}

// DeleteCategoryRule deletes a category rule
func (h *HTTPHandler) DeleteCategoryRule(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	err := h.txService.DeleteCategoryRule(c.Request.Context(), id, userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

func parseIntParam(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	var i int
	_, _ = fmt.Sscanf(val, "%d", &i)
	if i <= 0 {
		return defaultVal
	}
	return i
}

