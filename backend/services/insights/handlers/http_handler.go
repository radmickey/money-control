package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/insights/service"
)

// RegisterHTTPRoutes registers HTTP routes for insights
func RegisterHTTPRoutes(r *gin.RouterGroup, insightService *service.InsightService) {
	h := &HTTPHandler{insightService: insightService}

	insights := r.Group("/insights")
	{
		insights.GET("/net-worth", h.GetNetWorth)
		insights.GET("/net-worth/history", h.GetNetWorthHistory)
		insights.GET("/allocation", h.GetAllocation)
		insights.GET("/dashboard", h.GetDashboardSummary)
		insights.GET("/cash-flow", h.GetCashFlow)
		insights.POST("/snapshots", h.CreateSnapshot)
		insights.GET("/snapshots", h.GetSnapshots)
	}
}

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	insightService *service.InsightService
}

// GetNetWorth gets current net worth
func (h *HTTPHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	netWorth, err := h.insightService.GetNetWorth(c.Request.Context(), userID, baseCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, netWorth)
}

// GetNetWorthHistory gets net worth history
func (h *HTTPHandler) GetNetWorthHistory(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	period := c.Query("period")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	history, err := h.insightService.GetNetWorthHistory(c.Request.Context(), userID, baseCurrency, period, startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"history":  history,
		"currency": baseCurrency,
	})
}

// GetAllocation gets asset allocation
func (h *HTTPHandler) GetAllocation(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	groupBy := c.Query("group_by")

	allocations, totalValue, err := h.insightService.GetAllocation(c.Request.Context(), userID, baseCurrency, groupBy)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"allocations": allocations,
		"total_value": totalValue,
		"currency":    baseCurrency,
	})
}

// GetDashboardSummary gets dashboard summary
func (h *HTTPHandler) GetDashboardSummary(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")

	summary, err := h.insightService.GetDashboardSummary(c.Request.Context(), userID, baseCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, summary)
}

// GetCashFlow gets cash flow analysis
func (h *HTTPHandler) GetCashFlow(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	period := c.Query("period")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	periods, totalIncome, totalExpenses, netCashFlow, err := h.insightService.GetCashFlow(
		c.Request.Context(), userID, baseCurrency, period, startDate, endDate,
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"periods":        periods,
		"total_income":   totalIncome,
		"total_expenses": totalExpenses,
		"net_cash_flow":  netCashFlow,
		"currency":       baseCurrency,
	})
}

// CreateSnapshot creates a snapshot
func (h *HTTPHandler) CreateSnapshot(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")

	snapshot, err := h.insightService.CreateSnapshot(c.Request.Context(), userID, baseCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, snapshot)
}

// GetSnapshots gets historical snapshots
func (h *HTTPHandler) GetSnapshots(c *gin.Context) {
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

	snapshots, err := h.insightService.GetSnapshots(c.Request.Context(), userID, startDate, endDate, 30)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, snapshots)
}

