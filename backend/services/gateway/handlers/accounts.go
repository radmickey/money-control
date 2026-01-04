package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
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

// AccountWithConvertedBalance represents an account with converted total balance
type AccountWithConvertedBalance struct {
	ID                    string                    `json:"id"`
	UserID                string                    `json:"user_id"`
	Name                  string                    `json:"name"`
	Type                  string                    `json:"type"`
	Currency              string                    `json:"currency"`
	Description           string                    `json:"description,omitempty"`
	Icon                  string                    `json:"icon,omitempty"`
	TotalBalance          float64                   `json:"total_balance"`
	ConvertedTotalBalance float64                   `json:"converted_total_balance"`  // In account's currency
	BalanceInBaseCurrency float64                   `json:"balance_in_base_currency"` // In user's base currency (USD)
	DisplayCurrency       string                    `json:"display_currency"`
	IsMixedCurrency       bool                      `json:"is_mixed_currency"`
	SubAccounts           []SubAccountWithConverted `json:"subAccounts,omitempty"`
	CreatedAt             string                    `json:"created_at"`
	UpdatedAt             string                    `json:"updated_at"`
}

// SubAccountWithConverted represents a sub-account with converted balance
type SubAccountWithConverted struct {
	ID               string  `json:"id"`
	AccountID        string  `json:"account_id"`
	Name             string  `json:"name"`
	Currency         string  `json:"currency"`
	Balance          float64 `json:"balance"`
	ConvertedBalance float64 `json:"converted_balance"`
	Description      string  `json:"description,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ListAccounts lists accounts with converted balances
func (h *AccountsHandler) ListAccounts(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	accountType := c.Query("type")
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	// Get accounts
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

	// Get exchange rates
	ratesResp, err := h.proxy.Currency.GetMultipleExchangeRates(c.Request.Context(), &currencypb.GetMultipleExchangeRatesRequest{
		BaseCurrency: baseCurrency,
	})

	rates := make(map[string]float64)
	rates[baseCurrency] = 1.0
	if err != nil {
		log.Printf("Warning: Failed to get exchange rates: %v", err)
	} else if ratesResp != nil {
		for currency, rate := range ratesResp.Rates {
			rates[currency] = rate
		}
	}

	// Convert accounts
	accounts := make([]AccountWithConvertedBalance, 0, len(resp.Accounts))
	for _, acc := range resp.Accounts {
		accountCurrency := acc.Currency
		if accountCurrency == "" {
			accountCurrency = "USD"
		}

		// Convert sub-accounts and calculate totals
		var convertedTotal float64        // Total in account's currency
		var balanceInBaseCurrency float64 // Total in user's base currency (USD)
		var isMixed bool
		subAccounts := make([]SubAccountWithConverted, 0, len(acc.SubAccounts))

		currencies := make(map[string]bool)
		for _, sub := range acc.SubAccounts {
			subCurrency := sub.Currency
			if subCurrency == "" {
				subCurrency = accountCurrency
			}
			currencies[subCurrency] = true

			// Convert sub-account balance to account currency (for display)
			convertedBalance := convertAmount(sub.Balance, subCurrency, accountCurrency, rates)
			convertedTotal += convertedBalance

			// Convert sub-account balance to base currency (for net worth calculation)
			balanceInBase := convertAmount(sub.Balance, subCurrency, baseCurrency, rates)
			balanceInBaseCurrency += balanceInBase

			subAccounts = append(subAccounts, SubAccountWithConverted{
				ID:               sub.Id,
				AccountID:        sub.AccountId,
				Name:             sub.Name,
				Currency:         subCurrency,
				Balance:          sub.Balance,
				ConvertedBalance: convertedBalance,
				Description:      sub.Description,
				CreatedAt:        sub.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        sub.UpdatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
			})
		}

		// Check if mixed currencies
		if len(currencies) > 1 || (len(currencies) == 1 && !currencies[accountCurrency]) {
			isMixed = true
		}

		// Get account type string
		accType := accountTypeToString(acc.Type)

		accounts = append(accounts, AccountWithConvertedBalance{
			ID:                    acc.Id,
			UserID:                acc.UserId,
			Name:                  acc.Name,
			Type:                  accType,
			Currency:              accountCurrency,
			Description:           acc.Description,
			Icon:                  acc.Icon,
			TotalBalance:          acc.TotalBalance,
			ConvertedTotalBalance: convertedTotal,
			BalanceInBaseCurrency: balanceInBaseCurrency,
			DisplayCurrency:       accountCurrency,
			IsMixedCurrency:       isMixed,
			SubAccounts:           subAccounts,
			CreatedAt:             acc.CreatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:             acc.UpdatedAt.AsTime().Format("2006-01-02T15:04:05Z"),
		})
	}

	utils.Success(c, gin.H{
		"accounts":      accounts,
		"total":         resp.Total,
		"base_currency": baseCurrency,
		"rates":         rates,
	})
}

// convertAmount converts amount from one currency to another
func convertAmount(amount float64, from, to string, rates map[string]float64) float64 {
	if from == to {
		return amount
	}

	fromRate := rates[from]
	toRate := rates[to]

	if fromRate == 0 {
		fromRate = 1
	}
	if toRate == 0 {
		toRate = 1
	}

	// Convert: amount in 'from' -> base currency -> 'to' currency
	// rates are relative to base currency (e.g., USD)
	// If base is USD: EUR rate = 0.85 means 1 USD = 0.85 EUR
	// To convert EUR to USD: amount_EUR / 0.85 = amount_USD
	// To convert USD to EUR: amount_USD * 0.85 = amount_EUR
	inBase := amount / fromRate
	return inBase * toRate
}

// accountTypeToString converts account type enum to string
func accountTypeToString(t accountspb.AccountType) string {
	switch t {
	case accountspb.AccountType_ACCOUNT_TYPE_BANK:
		return "bank"
	case accountspb.AccountType_ACCOUNT_TYPE_CASH:
		return "cash"
	case accountspb.AccountType_ACCOUNT_TYPE_INVESTMENT:
		return "investment"
	case accountspb.AccountType_ACCOUNT_TYPE_CRYPTO:
		return "crypto"
	case accountspb.AccountType_ACCOUNT_TYPE_REAL_ESTATE:
		return "real_estate"
	case accountspb.AccountType_ACCOUNT_TYPE_OTHER:
		return "other"
	default:
		return "other"
	}
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
