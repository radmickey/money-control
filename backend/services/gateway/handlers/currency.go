package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/utils"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// CurrencyHandler handles currency-related requests
type CurrencyHandler struct {
	proxy *proxy.ServiceProxy
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(sp *proxy.ServiceProxy) *CurrencyHandler {
	return &CurrencyHandler{proxy: sp}
}

// ListCurrencies lists supported currencies
func (h *CurrencyHandler) ListCurrencies(c *gin.Context) {
	includeCrypto := c.Query("include_crypto") == "true"

	resp, err := h.proxy.Currency.ListSupportedCurrencies(c.Request.Context(), &currencypb.ListSupportedCurrenciesRequest{
		IncludeCrypto: includeCrypto,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp.Currencies)
}

// GetRates gets exchange rates
func (h *CurrencyHandler) GetRates(c *gin.Context) {
	base := c.Param("base")

	resp, err := h.proxy.Currency.GetMultipleExchangeRates(c.Request.Context(), &currencypb.GetMultipleExchangeRatesRequest{
		BaseCurrency: base,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"base":       resp.BaseCurrency,
		"rates":      resp.Rates,
		"updated_at": resp.UpdatedAt,
	})
}

// GetRate gets a single exchange rate
func (h *CurrencyHandler) GetRate(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")

	resp, err := h.proxy.Currency.GetExchangeRate(c.Request.Context(), &currencypb.GetExchangeRateRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"rate":       resp.Rate.Rate,
		"updated_at": resp.Rate.UpdatedAt,
	})
}

// Convert converts an amount
func (h *CurrencyHandler) Convert(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		From   string  `json:"from" binding:"required"`
		To     string  `json:"to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Currency.ConvertAmount(c.Request.Context(), &currencypb.ConvertAmountRequest{
		Amount:       req.Amount,
		FromCurrency: req.From,
		ToCurrency:   req.To,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"original_amount":  resp.OriginalAmount,
		"from_currency":    resp.FromCurrency,
		"converted_amount": resp.ConvertedAmount,
		"to_currency":      resp.ToCurrency,
		"rate_used":        resp.RateUsed,
	})
}
