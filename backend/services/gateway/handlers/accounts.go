package handlers

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/converters"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/resilience"
	"github.com/radmickey/money-control/backend/pkg/utils"
	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

const (
	serviceAccounts = "accounts-service"
	serviceCurrency = "currency-service"
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

	// Call with circuit breaker and timeout
	resp, err := resilience.Call(c.Request.Context(), resilience.DefaultCallOptions(serviceAccounts),
		func(ctx context.Context) (*accountspb.Account, error) {
			return h.proxy.Accounts.CreateAccount(ctx, &accountspb.CreateAccountRequest{
				UserId:      userID,
				Name:        req.Name,
				Type:        converters.StringToAccountType(req.Type),
				Currency:    converters.DefaultCurrency(req.Currency),
				Description: req.Description,
				Icon:        req.Icon,
			})
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
	// Support both "currency" and "baseCurrency" query params
	baseCurrency := c.Query("baseCurrency")
	if baseCurrency == "" {
		baseCurrency = c.Query("currency")
	}
	baseCurrency = converters.DefaultCurrency(baseCurrency)

	// Get accounts
	resp, err := h.proxy.Accounts.ListAccounts(c.Request.Context(), &accountspb.ListAccountsRequest{
		UserId:   userID,
		Type:     converters.StringToAccountType(accountType),
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	// Get exchange rates
	rates := h.fetchExchangeRates(c, baseCurrency)

	// Convert accounts
	accounts := h.convertAccountsWithRates(resp.Accounts, baseCurrency, rates)

	utils.Success(c, gin.H{
		"accounts":      accounts,
		"total":         resp.Total,
		"base_currency": baseCurrency,
		"rates":         rates,
	})
}

// fetchExchangeRates fetches exchange rates from currency service
func (h *AccountsHandler) fetchExchangeRates(c *gin.Context, baseCurrency string) map[string]float64 {
	ratesResp, err := h.proxy.Currency.GetMultipleExchangeRates(c.Request.Context(), &currencypb.GetMultipleExchangeRatesRequest{
		BaseCurrency: baseCurrency,
	})

	rates := converters.BuildRatesMap(baseCurrency, nil)
	if err != nil {
		log.Printf("Warning: Failed to get exchange rates: %v", err)
		return rates
	}
	if ratesResp != nil {
		return converters.BuildRatesMap(baseCurrency, ratesResp.Rates)
	}
	return rates
}

// convertAccountsWithRates converts account balances using exchange rates
func (h *AccountsHandler) convertAccountsWithRates(pbAccounts []*accountspb.Account, baseCurrency string, rates map[string]float64) []AccountWithConvertedBalance {
	accounts := make([]AccountWithConvertedBalance, 0, len(pbAccounts))

	for _, acc := range pbAccounts {
		accounts = append(accounts, h.convertSingleAccount(acc, baseCurrency, rates))
	}
	return accounts
}

// convertSingleAccount converts a single account with its sub-accounts
func (h *AccountsHandler) convertSingleAccount(acc *accountspb.Account, baseCurrency string, rates map[string]float64) AccountWithConvertedBalance {
	accountCurrency := converters.DefaultCurrency(acc.Currency)

	subAccounts, convertedTotal, balanceInBase, isMixed := h.convertSubAccounts(
		acc.SubAccounts, accountCurrency, baseCurrency, rates,
	)

	return AccountWithConvertedBalance{
		ID:                    acc.Id,
		UserID:                acc.UserId,
		Name:                  acc.Name,
		Type:                  converters.AccountTypeToString(acc.Type),
		Currency:              accountCurrency,
		Description:           acc.Description,
		Icon:                  acc.Icon,
		TotalBalance:          acc.TotalBalance,
		ConvertedTotalBalance: convertedTotal,
		BalanceInBaseCurrency: balanceInBase,
		DisplayCurrency:       accountCurrency,
		IsMixedCurrency:       isMixed,
		SubAccounts:           subAccounts,
		CreatedAt:             converters.FormatTime(acc.CreatedAt),
		UpdatedAt:             converters.FormatTime(acc.UpdatedAt),
	}
}

// convertSubAccounts converts sub-accounts and calculates totals
func (h *AccountsHandler) convertSubAccounts(subs []*accountspb.SubAccount, accountCurrency, baseCurrency string, rates map[string]float64) ([]SubAccountWithConverted, float64, float64, bool) {
	var convertedTotal, balanceInBase float64
	currencies := make(map[string]bool)
	subAccounts := make([]SubAccountWithConverted, 0, len(subs))

	for _, sub := range subs {
		subCurrency := converters.DefaultCurrency(sub.Currency)
		if subCurrency == "" {
			subCurrency = accountCurrency
		}
		currencies[subCurrency] = true

		// Convert balances
		convertedBalance := converters.ConvertAmount(sub.Balance, subCurrency, accountCurrency, rates)
		convertedTotal += convertedBalance
		balanceInBase += converters.ConvertAmount(sub.Balance, subCurrency, baseCurrency, rates)

		subAccounts = append(subAccounts, SubAccountWithConverted{
			ID:               sub.Id,
			AccountID:        sub.AccountId,
			Name:             sub.Name,
			Currency:         subCurrency,
			Balance:          sub.Balance,
			ConvertedBalance: convertedBalance,
			Description:      sub.Description,
			CreatedAt:        converters.FormatTime(sub.CreatedAt),
			UpdatedAt:        converters.FormatTime(sub.UpdatedAt),
		})
	}

	// Check if mixed currencies
	isMixed := len(currencies) > 1 || (len(currencies) == 1 && !currencies[accountCurrency])

	return subAccounts, convertedTotal, balanceInBase, isMixed
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
		AssetType:   converters.StringToAssetType(assetType),
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
