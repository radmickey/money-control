package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	assetspb "github.com/radmickey/money-control/backend/proto/assets"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// AssetsHandler handles assets-related requests
type AssetsHandler struct {
	proxy *proxy.ServiceProxy
}

// NewAssetsHandler creates a new assets handler
func NewAssetsHandler(sp *proxy.ServiceProxy) *AssetsHandler {
	return &AssetsHandler{proxy: sp}
}

// CreateAsset creates a new asset
func (h *AssetsHandler) CreateAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req struct {
		SubAccountID  string  `json:"sub_account_id"`
		Symbol        string  `json:"symbol" binding:"required"`
		Name          string  `json:"name"`
		Type          string  `json:"type" binding:"required"`
		Quantity      float64 `json:"quantity" binding:"required"`
		PurchasePrice float64 `json:"purchase_price"`
		Currency      string  `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Assets.CreateAsset(c.Request.Context(), &assetspb.CreateAssetRequest{
		UserId:        userID,
		SubAccountId:  req.SubAccountID,
		Symbol:        req.Symbol,
		Name:          req.Name,
		Type:          stringToAssetTypeAssets(req.Type),
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		Currency:      req.Currency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, resp)
}

// ListAssets lists assets
func (h *AssetsHandler) ListAssets(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	subAccountID := c.Query("sub_account_id")
	assetType := c.Query("type")

	resp, err := h.proxy.Assets.ListAssets(c.Request.Context(), &assetspb.ListAssetsRequest{
		UserId:       userID,
		SubAccountId: subAccountID,
		Type:         stringToAssetTypeAssets(assetType),
		Page:         1,
		PageSize:     100,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"assets":            resp.Assets,
		"total":             resp.Total,
		"total_value":       resp.TotalValue,
		"total_profit_loss": resp.TotalProfitLoss,
	})
}

// GetAsset gets an asset
func (h *AssetsHandler) GetAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	resp, err := h.proxy.Assets.GetAsset(c.Request.Context(), &assetspb.GetAssetRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.NotFound(c, "Asset not found")
		return
	}

	utils.Success(c, resp)
}

// UpdateAsset updates an asset
func (h *AssetsHandler) UpdateAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	var req struct {
		Quantity      float64 `json:"quantity"`
		PurchasePrice float64 `json:"purchase_price"`
		CurrentPrice  float64 `json:"current_price"`
		Name          string  `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Assets.UpdateAsset(c.Request.Context(), &assetspb.UpdateAssetRequest{
		Id:            id,
		UserId:        userID,
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		CurrentPrice:  req.CurrentPrice,
		Name:          req.Name,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// DeleteAsset deletes an asset
func (h *AssetsHandler) DeleteAsset(c *gin.Context) {
	userID := middleware.MustGetUserID(c)
	id := c.Param("id")

	_, err := h.proxy.Assets.DeleteAsset(c.Request.Context(), &assetspb.DeleteAssetRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.NoContent(c)
}

// GetPrice gets asset price
func (h *AssetsHandler) GetPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	assetType := c.Query("type")

	resp, err := h.proxy.Assets.GetAssetPrice(c.Request.Context(), &assetspb.GetAssetPriceRequest{
		Symbol: symbol,
		Type:   stringToAssetTypeAssets(assetType),
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// SearchAssets searches for assets
func (h *AssetsHandler) SearchAssets(c *gin.Context) {
	query := c.Query("q")
	assetType := c.Query("type")

	resp, err := h.proxy.Assets.SearchAssets(c.Request.Context(), &assetspb.SearchAssetsRequest{
		Query: query,
		Type:  stringToAssetTypeAssets(assetType),
		Limit: 10,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp.Results)
}

// Helper functions
func stringToAssetTypeAssets(s string) assetspb.AssetType {
	switch s {
	case "STOCK", "stock":
		return assetspb.AssetType_ASSET_TYPE_STOCK
	case "CRYPTO", "crypto":
		return assetspb.AssetType_ASSET_TYPE_CRYPTO
	case "ETF", "etf":
		return assetspb.AssetType_ASSET_TYPE_ETF
	case "REAL_ESTATE", "real_estate":
		return assetspb.AssetType_ASSET_TYPE_REAL_ESTATE
	case "CASH", "cash":
		return assetspb.AssetType_ASSET_TYPE_CASH
	case "BOND", "bond":
		return assetspb.AssetType_ASSET_TYPE_BOND
	default:
		return assetspb.AssetType_ASSET_TYPE_OTHER
	}
}
