package handlers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/converters"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
	insightspb "github.com/radmickey/money-control/backend/proto/insights"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// InsightsHandler handles insights-related requests
type InsightsHandler struct {
	proxy *proxy.ServiceProxy
}

// NewInsightsHandler creates a new insights handler
func NewInsightsHandler(sp *proxy.ServiceProxy) *InsightsHandler {
	return &InsightsHandler{proxy: sp}
}

// GetNetWorth gets current net worth with proper currency conversion
func (h *InsightsHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("baseCurrency")
	if baseCurrency == "" {
		baseCurrency = c.Query("currency")
	}
	baseCurrency = converters.DefaultCurrency(baseCurrency)

	// Get all accounts
	accountsResp, err := h.proxy.Accounts.ListAccounts(c.Request.Context(), &accountspb.ListAccountsRequest{
		UserId:   userID,
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	// Get exchange rates
	rates := h.fetchExchangeRates(c, baseCurrency)

	// Calculate total net worth with conversion
	var totalNetWorth float64
	for _, acc := range accountsResp.Accounts {
		accountCurrency := acc.Currency
		if accountCurrency == "" {
			accountCurrency = baseCurrency
		}

		// Sum all sub-accounts with conversion
		for _, sub := range acc.SubAccounts {
			subCurrency := sub.Currency
			if subCurrency == "" {
				subCurrency = accountCurrency
			}

			// Convert to base currency
			converted := converters.ConvertAmount(sub.Balance, subCurrency, baseCurrency, rates)
			totalNetWorth += converted
		}
	}

	// Get historical data for change calculations (from Insights service)
	resp, err := h.proxy.Insights.GetNetWorth(c.Request.Context(), &insightspb.GetNetWorthRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
	})

	// Return combined result
	result := gin.H{
		"total":             totalNetWorth,
		"currency":          baseCurrency,
		"change24h":         0.0,
		"changePercent24h":  0.0,
		"change7d":          0.0,
		"changePercent7d":   0.0,
		"change30d":         0.0,
		"changePercent30d":  0.0,
		"calculatedAt":      time.Now(),
	}

	// Add change data from Insights if available
	if err == nil && resp != nil {
		result["change24h"] = resp.Change_24H
		result["changePercent24h"] = resp.ChangePercent_24H
		result["change7d"] = resp.Change_7D
		result["changePercent7d"] = resp.ChangePercent_7D
		result["change30d"] = resp.Change_30D
		result["changePercent30d"] = resp.ChangePercent_30D
	}

	utils.Success(c, result)
}

// fetchExchangeRates fetches exchange rates from currency service
func (h *InsightsHandler) fetchExchangeRates(c *gin.Context, baseCurrency string) map[string]float64 {
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

// GetTrends gets net worth trends/history
func (h *InsightsHandler) GetTrends(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := converters.DefaultCurrency(c.Query("currency"))
	period := c.Query("period")
	if period == "" {
		period = "30d"
	}

	resp, err := h.proxy.Insights.GetTrends(c.Request.Context(), &insightspb.GetTrendsRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		Period:       period,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetAllocation gets asset allocation
func (h *InsightsHandler) GetAllocation(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	// Support both "baseCurrency" and "currency" query params
	baseCurrency := c.Query("baseCurrency")
	if baseCurrency == "" {
		baseCurrency = c.Query("currency")
	}
	baseCurrency = converters.DefaultCurrency(baseCurrency)
	groupBy := c.Query("group_by")
	if groupBy == "" {
		groupBy = "asset_type"
	}

	resp, err := h.proxy.Insights.GetAllocation(c.Request.Context(), &insightspb.GetAllocationRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		GroupBy:      groupBy,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetDashboard gets dashboard summary
func (h *InsightsHandler) GetDashboard(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := converters.DefaultCurrency(c.Query("currency"))

	resp, err := h.proxy.Insights.GetDashboardSummary(c.Request.Context(), &insightspb.GetDashboardSummaryRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetCashFlow gets cash flow
func (h *InsightsHandler) GetCashFlow(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := converters.DefaultCurrency(c.Query("currency"))
	period := c.Query("period")
	if period == "" {
		period = "monthly"
	}

	resp, err := h.proxy.Insights.GetCashFlow(c.Request.Context(), &insightspb.GetCashFlowRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		Period:       period,
		StartDate:    converters.ParseDate(c.Query("start_date")),
		EndDate:      converters.ParseDate(c.Query("end_date")),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetNetWorthHistory gets net worth history
func (h *InsightsHandler) GetNetWorthHistory(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := converters.DefaultCurrency(c.Query("currency"))
	period := c.Query("period")
	if period == "" {
		period = "30d"
	}

	resp, err := h.proxy.Insights.GetNetWorthHistory(c.Request.Context(), &insightspb.GetNetWorthHistoryRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		Period:       period,
		StartDate:    converters.ParseDate(c.Query("start_date")),
		EndDate:      converters.ParseDate(c.Query("end_date")),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}
