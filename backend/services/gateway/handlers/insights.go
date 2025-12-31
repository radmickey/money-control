package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	insightspb "github.com/radmickey/money-control/backend/proto/insights"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// InsightsHandler handles insights-related requests
type InsightsHandler struct {
	proxy *proxy.ServiceProxy
}

// NewInsightsHandler creates a new insights handler
func NewInsightsHandler(sp *proxy.ServiceProxy) *InsightsHandler {
	return &InsightsHandler{proxy: sp}
}

// GetNetWorth gets current net worth
func (h *InsightsHandler) GetNetWorth(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	resp, err := h.proxy.Insights.GetNetWorth(c.Request.Context(), &insightspb.GetNetWorthRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetTrends gets net worth trends/history
func (h *InsightsHandler) GetTrends(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}
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
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}
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
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

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
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}
	period := c.Query("period")
	if period == "" {
		period = "monthly"
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate *timestamppb.Timestamp
	if startDateStr != "" {
		t, _ := time.Parse("2006-01-02", startDateStr)
		startDate = timestamppb.New(t)
	}
	if endDateStr != "" {
		t, _ := time.Parse("2006-01-02", endDateStr)
		endDate = timestamppb.New(t)
	}

	resp, err := h.proxy.Insights.GetCashFlow(c.Request.Context(), &insightspb.GetCashFlowRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		Period:       period,
		StartDate:    startDate,
		EndDate:      endDate,
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
	baseCurrency := c.Query("currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}
	period := c.Query("period")
	if period == "" {
		period = "30d"
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate *timestamppb.Timestamp
	if startDateStr != "" {
		t, _ := time.Parse("2006-01-02", startDateStr)
		startDate = timestamppb.New(t)
	}
	if endDateStr != "" {
		t, _ := time.Parse("2006-01-02", endDateStr)
		endDate = timestamppb.New(t)
	}

	resp, err := h.proxy.Insights.GetNetWorthHistory(c.Request.Context(), &insightspb.GetNetWorthHistoryRequest{
		UserId:       userID,
		BaseCurrency: baseCurrency,
		Period:       period,
		StartDate:    startDate,
		EndDate:      endDate,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}
