package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/currency/service"
)

// RegisterHTTPRoutes registers HTTP routes for currency
func RegisterHTTPRoutes(r *gin.RouterGroup, currencyService *service.CurrencyService, jwtManager *auth.JWTManager) {
	h := &HTTPHandler{currencyService: currencyService}

	// Public routes
	currencies := r.Group("/currencies")
	{
		currencies.GET("", h.ListCurrencies)
		currencies.GET("/rates/:base", h.GetRates)
		currencies.GET("/rate/:from/:to", h.GetRate)
		currencies.POST("/convert", h.Convert)
	}

	// Protected routes
	protected := r.Group("/currencies")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.POST("/convert-multiple", h.ConvertMultiple)
		protected.POST("/refresh", h.RefreshRates)
	}
}

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	currencyService *service.CurrencyService
}

// ListCurrencies lists all supported currencies
func (h *HTTPHandler) ListCurrencies(c *gin.Context) {
	includeCrypto := c.Query("include_crypto") == "true"

	currencies, err := h.currencyService.ListSupportedCurrencies(c.Request.Context(), includeCrypto)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, currencies)
}

// GetRates gets all rates for a base currency
func (h *HTTPHandler) GetRates(c *gin.Context) {
	base := c.Param("base")

	rates, updatedAt, err := h.currencyService.GetMultipleExchangeRates(c.Request.Context(), base, nil)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"base":       base,
		"rates":      rates,
		"updated_at": updatedAt,
	})
}

// GetRate gets exchange rate between two currencies
func (h *HTTPHandler) GetRate(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")

	rate, err := h.currencyService.GetExchangeRate(c.Request.Context(), from, to)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, rate)
}

// ConvertRequest represents convert request
type ConvertRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
}

// Convert converts an amount between currencies
func (h *HTTPHandler) Convert(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	converted, rate, err := h.currencyService.ConvertAmount(c.Request.Context(), req.Amount, req.From, req.To)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"original_amount":  req.Amount,
		"from_currency":    req.From,
		"converted_amount": converted,
		"to_currency":      req.To,
		"rate_used":        rate,
		"timestamp":        time.Now(),
	})
}

// ConvertMultipleRequest represents convert multiple request
type ConvertMultipleRequest struct {
	Amounts []struct {
		ID           string  `json:"id"`
		Amount       float64 `json:"amount"`
		FromCurrency string  `json:"from_currency"`
	} `json:"amounts" binding:"required"`
	ToCurrency string `json:"to_currency" binding:"required"`
}

// ConvertMultiple converts multiple amounts
func (h *HTTPHandler) ConvertMultiple(c *gin.Context) {
	var req ConvertMultipleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	amounts := make([]service.AmountToConvert, len(req.Amounts))
	for i, a := range req.Amounts {
		amounts[i] = service.AmountToConvert{
			ID:           a.ID,
			Amount:       a.Amount,
			FromCurrency: a.FromCurrency,
		}
	}

	results, total, err := h.currencyService.ConvertMultipleAmounts(c.Request.Context(), amounts, req.ToCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"converted":   results,
		"total":       total,
		"to_currency": req.ToCurrency,
	})
}

// RefreshRates refreshes exchange rates
func (h *HTTPHandler) RefreshRates(c *gin.Context) {
	baseCurrency := c.Query("base")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	count, err := h.currencyService.RefreshRates(c.Request.Context(), baseCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"success":       true,
		"rates_updated": count,
		"updated_at":    time.Now(),
	})
}

