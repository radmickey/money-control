package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// AccountsHandler handles accounts-related requests
type AccountsHandler struct {
	proxy *proxy.ServiceProxy
}

// NewAccountsHandler creates a new accounts handler
func NewAccountsHandler(sp *proxy.ServiceProxy) *AccountsHandler {
	return &AccountsHandler{proxy: sp}
}

// CreateAccount creates a new account
func (h *AccountsHandler) CreateAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Accounts.CreateAccount(c.Request.Context(), &accountspb.CreateAccountRequest{
		UserId:      userID,
		Name:        req.Name,
		Type:        stringToAccountType(req.Type),
		Currency:    req.Currency,
		Description: req.Description,
		Icon:        req.Icon,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, resp)
}

// ListAccounts lists accounts
func (h *AccountsHandler) ListAccounts(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountType := c.Query("type")

	resp, err := h.proxy.Accounts.ListAccounts(c.Request.Context(), &accountspb.ListAccountsRequest{
		UserId:   userID,
		Type:     stringToAccountType(accountType),
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"accounts": resp.Accounts,
		"total":    resp.Total,
	})
}

// GetAccount gets an account
func (h *AccountsHandler) GetAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	resp, err := h.proxy.Accounts.GetAccount(c.Request.Context(), &accountspb.GetAccountRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.NotFound(c, "Account not found")
		return
	}

	utils.Success(c, resp)
}

// UpdateAccount updates an account
func (h *AccountsHandler) UpdateAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Currency    string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Accounts.UpdateAccount(c.Request.Context(), &accountspb.UpdateAccountRequest{
		Id:          id,
		UserId:      userID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Currency:    req.Currency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// DeleteAccount deletes an account
func (h *AccountsHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	_, err := h.proxy.Accounts.DeleteAccount(c.Request.Context(), &accountspb.DeleteAccountRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// CreateSubAccount creates a sub-account
func (h *AccountsHandler) CreateSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountID := c.Param("id")

	var req struct {
		Name        string  `json:"name" binding:"required"`
		AssetType   string  `json:"asset_type"`
		Currency    string  `json:"currency"`
		Balance     float64 `json:"balance"`
		Symbol      string  `json:"symbol"`
		Quantity    float64 `json:"quantity"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Default asset type to "cash" if not provided
	assetType := req.AssetType
	if assetType == "" {
		assetType = "cash"
	}

	resp, err := h.proxy.Accounts.CreateSubAccount(c.Request.Context(), &accountspb.CreateSubAccountRequest{
		AccountId:   accountID,
		UserId:      userID,
		Name:        req.Name,
		AssetType:   stringToAssetType(assetType),
		Currency:    req.Currency,
		Balance:     req.Balance,
		Symbol:      req.Symbol,
		Quantity:    req.Quantity,
		Description: req.Description,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, resp)
}

// ListSubAccounts lists sub-accounts
func (h *AccountsHandler) ListSubAccounts(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountID := c.Param("id")

	resp, err := h.proxy.Accounts.ListSubAccounts(c.Request.Context(), &accountspb.ListSubAccountsRequest{
		AccountId: accountID,
		UserId:    userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp.SubAccounts)
}

// GetSubAccount gets a sub-account
func (h *AccountsHandler) GetSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	resp, err := h.proxy.Accounts.GetSubAccount(c.Request.Context(), &accountspb.GetSubAccountRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.NotFound(c, "Sub-account not found")
		return
	}

	utils.Success(c, resp)
}

// UpdateSubAccount updates a sub-account
func (h *AccountsHandler) UpdateSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Balance     float64 `json:"balance"`
		Quantity    float64 `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Accounts.UpdateSubAccount(c.Request.Context(), &accountspb.UpdateSubAccountRequest{
		Id:          id,
		UserId:      userID,
		Name:        req.Name,
		Description: req.Description,
		Balance:     req.Balance,
		Quantity:    req.Quantity,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// DeleteSubAccount deletes a sub-account
func (h *AccountsHandler) DeleteSubAccount(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	_, err := h.proxy.Accounts.DeleteSubAccount(c.Request.Context(), &accountspb.DeleteSubAccountRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// UpdateSubAccountBalance updates sub-account balance
func (h *AccountsHandler) UpdateSubAccountBalance(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req struct {
		Balance  float64 `json:"balance" binding:"required"`
		Quantity float64 `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Accounts.UpdateSubAccountBalance(c.Request.Context(), &accountspb.UpdateSubAccountBalanceRequest{
		Id:       id,
		UserId:   userID,
		Balance:  req.Balance,
		Quantity: req.Quantity,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetNetWorth gets user net worth
func (h *AccountsHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	resp, err := h.proxy.Accounts.GetUserNetWorth(c.Request.Context(), &accountspb.GetUserNetWorthRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Helper functions
func stringToAccountType(s string) accountspb.AccountType {
	switch s {
	case "BANK", "bank":
		return accountspb.AccountType_ACCOUNT_TYPE_BANK
	case "CASH", "cash":
		return accountspb.AccountType_ACCOUNT_TYPE_CASH
	case "INVESTMENT", "investment":
		return accountspb.AccountType_ACCOUNT_TYPE_INVESTMENT
	case "CRYPTO", "crypto":
		return accountspb.AccountType_ACCOUNT_TYPE_CRYPTO
	case "REAL_ESTATE", "real_estate":
		return accountspb.AccountType_ACCOUNT_TYPE_REAL_ESTATE
	case "OTHER", "other":
		return accountspb.AccountType_ACCOUNT_TYPE_OTHER
	default:
		return accountspb.AccountType_ACCOUNT_TYPE_UNSPECIFIED
	}
}

func stringToAssetType(s string) accountspb.AssetType {
	switch s {
	case "BANK", "bank":
		return accountspb.AssetType_ASSET_TYPE_BANK
	case "CASH", "cash":
		return accountspb.AssetType_ASSET_TYPE_CASH
	case "STOCKS", "stocks":
		return accountspb.AssetType_ASSET_TYPE_STOCKS
	case "CRYPTO", "crypto":
		return accountspb.AssetType_ASSET_TYPE_CRYPTO
	case "ETF", "etf":
		return accountspb.AssetType_ASSET_TYPE_ETF
	case "REAL_ESTATE", "real_estate":
		return accountspb.AssetType_ASSET_TYPE_REAL_ESTATE
	case "BONDS", "bonds":
		return accountspb.AssetType_ASSET_TYPE_BONDS
	default:
		return accountspb.AssetType_ASSET_TYPE_OTHER
	}
}
