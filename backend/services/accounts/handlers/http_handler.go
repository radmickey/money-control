package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/accounts/models"
	"github.com/radmickey/money-control/backend/services/accounts/repository"
	"github.com/radmickey/money-control/backend/services/accounts/service"
)

// RegisterHTTPRoutes registers HTTP routes for accounts
func RegisterHTTPRoutes(r *gin.RouterGroup, accountService *service.AccountService) {
	h := &HTTPHandler{accountService: accountService}

	accounts := r.Group("/accounts")
	{
		accounts.POST("", h.CreateAccount)
		accounts.GET("", h.ListAccounts)
		accounts.GET("/:id", h.GetAccount)
		accounts.PUT("/:id", h.UpdateAccount)
		accounts.DELETE("/:id", h.DeleteAccount)

		// Sub-accounts
		accounts.POST("/:id/sub-accounts", h.CreateSubAccount)
		accounts.GET("/:id/sub-accounts", h.ListSubAccounts)
	}

	subAccounts := r.Group("/sub-accounts")
	{
		subAccounts.GET("/:id", h.GetSubAccount)
		subAccounts.PUT("/:id", h.UpdateSubAccount)
		subAccounts.DELETE("/:id", h.DeleteSubAccount)
		subAccounts.PATCH("/:id/balance", h.UpdateSubAccountBalance)
	}

	// Net worth and summary
	r.GET("/net-worth", h.GetNetWorth)
	r.GET("/accounts-summary", h.GetAccountsSummary)
}

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	accountService *service.AccountService
}

// CreateAccountRequest represents create account request
type CreateAccountRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// CreateAccount creates a new account
func (h *HTTPHandler) CreateAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.CreateAccountInput{
		UserID:      userID,
		Name:        req.Name,
		Type:        models.AccountType(req.Type),
		Currency:    req.Currency,
		Description: req.Description,
		Icon:        req.Icon,
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), input)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, account)
}

// ListAccounts lists accounts
func (h *HTTPHandler) ListAccounts(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	accountType := models.AccountType(c.Query("type"))
	page := utils.CoalesceInt(parseIntParam(c, "page"), 1)
	pageSize := utils.CoalesceInt(parseIntParam(c, "page_size"), 20)

	accounts, total, err := h.accountService.ListAccounts(c.Request.Context(), userID, accountType, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMeta(c, accounts, utils.PaginationMeta(page, pageSize, int(total)))
}

// GetAccount gets an account by ID
func (h *HTTPHandler) GetAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	account, err := h.accountService.GetAccount(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			utils.NotFound(c, "Account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, account)
}

// UpdateAccountRequest represents update account request
type UpdateAccountRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Currency    string `json:"currency"`
}

// UpdateAccount updates an account
func (h *HTTPHandler) UpdateAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.UpdateAccountInput{
		ID:          id,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Currency:    req.Currency,
	}

	account, err := h.accountService.UpdateAccount(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			utils.NotFound(c, "Account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, account)
}

// DeleteAccount deletes an account
func (h *HTTPHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	err := h.accountService.DeleteAccount(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			utils.NotFound(c, "Account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// CreateSubAccountRequest represents create sub-account request
type CreateSubAccountRequest struct {
	Name        string  `json:"name" binding:"required"`
	AssetType   string  `json:"asset_type" binding:"required"`
	Currency    string  `json:"currency"`
	Balance     float64 `json:"balance"`
	Symbol      string  `json:"symbol"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
}

// CreateSubAccount creates a new sub-account
func (h *HTTPHandler) CreateSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountID := c.Param("id")

	var req CreateSubAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.CreateSubAccountInput{
		AccountID:   accountID,
		UserID:      userID,
		Name:        req.Name,
		AssetType:   models.AssetType(req.AssetType),
		Currency:    req.Currency,
		Balance:     req.Balance,
		Symbol:      req.Symbol,
		Quantity:    req.Quantity,
		Description: req.Description,
	}

	subAccount, err := h.accountService.CreateSubAccount(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			utils.NotFound(c, "Account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, subAccount)
}

// ListSubAccounts lists sub-accounts for an account
func (h *HTTPHandler) ListSubAccounts(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountID := c.Param("id")
	assetType := models.AssetType(c.Query("asset_type"))

	subAccounts, err := h.accountService.ListSubAccounts(c.Request.Context(), accountID, userID, assetType)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, subAccounts)
}

// GetSubAccount gets a sub-account by ID
func (h *HTTPHandler) GetSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	subAccount, err := h.accountService.GetSubAccount(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrSubAccountNotFound {
			utils.NotFound(c, "Sub-account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, subAccount)
}

// UpdateSubAccountRequest represents update sub-account request
type UpdateSubAccountRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Balance     float64 `json:"balance"`
	Quantity    float64 `json:"quantity"`
}

// UpdateSubAccount updates a sub-account
func (h *HTTPHandler) UpdateSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req UpdateSubAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.UpdateSubAccountInput{
		ID:          id,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Balance:     req.Balance,
		Quantity:    req.Quantity,
	}

	subAccount, err := h.accountService.UpdateSubAccount(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrSubAccountNotFound {
			utils.NotFound(c, "Sub-account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, subAccount)
}

// DeleteSubAccount deletes a sub-account
func (h *HTTPHandler) DeleteSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	err := h.accountService.DeleteSubAccount(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrSubAccountNotFound {
			utils.NotFound(c, "Sub-account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// UpdateBalanceRequest represents update balance request
type UpdateBalanceRequest struct {
	Balance  float64 `json:"balance" binding:"required"`
	Quantity float64 `json:"quantity"`
}

// UpdateSubAccountBalance updates sub-account balance
func (h *HTTPHandler) UpdateSubAccountBalance(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	subAccount, err := h.accountService.UpdateSubAccountBalance(c.Request.Context(), id, userID, req.Balance, req.Quantity)
	if err != nil {
		if err == repository.ErrSubAccountNotFound {
			utils.NotFound(c, "Sub-account not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, subAccount)
}

// GetNetWorth gets user's net worth
func (h *HTTPHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	balances, err := h.accountService.GetUserNetWorth(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	assetBalances, err := h.accountService.GetBalanceByAssetType(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"by_currency":   balances,
		"by_asset_type": assetBalances,
	})
}

// GetAccountsSummary gets accounts summary
func (h *HTTPHandler) GetAccountsSummary(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	summary, err := h.accountService.GetAccountsSummary(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, summary)
}

func parseIntParam(c *gin.Context, key string) int {
	val := c.Query(key)
	if val == "" {
		return 0
	}
	var i int
	_, _ = fmt.Sscanf(val, "%d", &i)
	return i
}

