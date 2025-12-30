package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/assets/models"
	"github.com/radmickey/money-control/backend/services/assets/repository"
	"github.com/radmickey/money-control/backend/services/assets/service"
)

// RegisterHTTPRoutes registers HTTP routes for assets
func RegisterHTTPRoutes(r *gin.RouterGroup, assetService *service.AssetService) {
	h := &HTTPHandler{assetService: assetService}

	assets := r.Group("/assets")
	{
		assets.POST("", h.CreateAsset)
		assets.GET("", h.ListAssets)
		assets.GET("/:id", h.GetAsset)
		assets.PUT("/:id", h.UpdateAsset)
		assets.DELETE("/:id", h.DeleteAsset)
		assets.POST("/refresh-prices", h.RefreshPrices)
	}

	prices := r.Group("/prices")
	{
		prices.GET("/:symbol", h.GetPrice)
		prices.POST("/multiple", h.GetMultiplePrices)
	}

	r.GET("/portfolio", h.GetPortfolioPerformance)
	r.GET("/search", h.SearchAssets)
}

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	assetService *service.AssetService
}

// CreateAssetRequest represents create asset request
type CreateAssetRequest struct {
	SubAccountID  string  `json:"sub_account_id"`
	Symbol        string  `json:"symbol" binding:"required"`
	Name          string  `json:"name"`
	Type          string  `json:"type" binding:"required"`
	Quantity      float64 `json:"quantity" binding:"required"`
	PurchasePrice float64 `json:"purchase_price"`
	Currency      string  `json:"currency"`
	PurchaseDate  string  `json:"purchase_date"`
}

// CreateAsset creates a new asset
func (h *HTTPHandler) CreateAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var purchaseDate *time.Time
	if req.PurchaseDate != "" {
		t, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err == nil {
			purchaseDate = &t
		}
	}

	var subAccountID *string
	if req.SubAccountID != "" {
		subAccountID = &req.SubAccountID
	}

	input := service.CreateAssetInput{
		UserID:        userID,
		SubAccountID:  subAccountID,
		Symbol:        req.Symbol,
		Name:          req.Name,
		Type:          models.AssetType(req.Type),
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		Currency:      req.Currency,
		PurchaseDate:  purchaseDate,
	}

	asset, err := h.assetService.CreateAsset(c.Request.Context(), input)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, asset)
}

// ListAssets lists assets
func (h *HTTPHandler) ListAssets(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	subAccountID := c.Query("sub_account_id")
	assetType := models.AssetType(c.Query("type"))
	page := parseIntParam(c, "page", 1)
	pageSize := parseIntParam(c, "page_size", 20)

	assets, total, totalValue, totalPL, err := h.assetService.ListAssets(
		c.Request.Context(), userID, subAccountID, assetType, page, pageSize,
	)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMeta(c, gin.H{
		"assets":            assets,
		"total_value":       totalValue,
		"total_profit_loss": totalPL,
	}, utils.PaginationMeta(page, pageSize, int(total)))
}

// GetAsset gets an asset by ID
func (h *HTTPHandler) GetAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	asset, err := h.assetService.GetAsset(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrAssetNotFound {
			utils.NotFound(c, "Asset not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, asset)
}

// UpdateAssetRequest represents update asset request
type UpdateAssetRequest struct {
	Quantity      float64 `json:"quantity"`
	PurchasePrice float64 `json:"purchase_price"`
	CurrentPrice  float64 `json:"current_price"`
	Name          string  `json:"name"`
}

// UpdateAsset updates an asset
func (h *HTTPHandler) UpdateAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.UpdateAssetInput{
		ID:            id,
		UserID:        userID,
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		CurrentPrice:  req.CurrentPrice,
		Name:          req.Name,
		Metadata:      nil,
	}

	asset, err := h.assetService.UpdateAsset(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrAssetNotFound {
			utils.NotFound(c, "Asset not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, asset)
}

// DeleteAsset deletes an asset
func (h *HTTPHandler) DeleteAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	err := h.assetService.DeleteAsset(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrAssetNotFound {
			utils.NotFound(c, "Asset not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// GetPrice gets price for a symbol
func (h *HTTPHandler) GetPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	assetType := models.AssetType(c.Query("type"))

	price, err := h.assetService.GetAssetPrice(c.Request.Context(), symbol, assetType)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, price)
}

// GetMultiplePricesRequest represents get multiple prices request
type GetMultiplePricesRequest struct {
	Queries []struct {
		Symbol string `json:"symbol"`
		Type   string `json:"type"`
	} `json:"queries"`
}

// GetMultiplePrices gets prices for multiple symbols
func (h *HTTPHandler) GetMultiplePrices(c *gin.Context) {
	var req GetMultiplePricesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	queries := make(map[string]string)
	for _, q := range req.Queries {
		queries[q.Symbol] = q.Type
	}

	prices, failedSymbols, err := h.assetService.GetMultipleAssetPrices(c.Request.Context(), queries)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"prices":         prices,
		"failed_symbols": failedSymbols,
	})
}

// RefreshPricesRequest represents refresh prices request
type RefreshPricesRequest struct {
	AssetIDs []string `json:"asset_ids"`
}

// RefreshPrices refreshes prices for assets
func (h *HTTPHandler) RefreshPrices(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req RefreshPricesRequest
	_ = c.ShouldBindJSON(&req)

	updated, failed, err := h.assetService.RefreshAssetPrices(c.Request.Context(), userID, req.AssetIDs)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"updated_count": updated,
		"failed_ids":    failed,
	})
}

// GetPortfolioPerformance gets portfolio performance
func (h *HTTPHandler) GetPortfolioPerformance(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	perf, err := h.assetService.GetPortfolioPerformance(c.Request.Context(), userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, perf)
}

// SearchAssets searches for assets
func (h *HTTPHandler) SearchAssets(c *gin.Context) {
	query := c.Query("q")
	assetType := models.AssetType(c.Query("type"))
	limit := parseIntParam(c, "limit", 10)

	if query == "" {
		utils.BadRequest(c, "Query parameter 'q' is required")
		return
	}

	results, err := h.assetService.SearchAssets(c.Request.Context(), query, assetType, limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, results)
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

