package service

import (
	"context"
	"sync"
	"time"

	"github.com/radmickey/money-control/backend/pkg/cache"
	"github.com/radmickey/money-control/backend/services/assets/models"
	"github.com/radmickey/money-control/backend/services/assets/providers"
	"github.com/radmickey/money-control/backend/services/assets/repository"
)

const (
	priceCacheTTL = 5 * time.Minute
)

// AssetService handles asset business logic
type AssetService struct {
	assetRepo      *repository.AssetRepository
	priceCacheRepo *repository.PriceCacheRepository
	priceManager   *providers.PriceManager
	redisCache     *cache.Cache
}

// NewAssetService creates a new asset service
func NewAssetService(
	assetRepo *repository.AssetRepository,
	priceCacheRepo *repository.PriceCacheRepository,
	priceManager *providers.PriceManager,
	redisCache *cache.Cache,
) *AssetService {
	return &AssetService{
		assetRepo:      assetRepo,
		priceCacheRepo: priceCacheRepo,
		priceManager:   priceManager,
		redisCache:     redisCache,
	}
}

// CreateAssetInput holds input for creating an asset
type CreateAssetInput struct {
	UserID        string
	SubAccountID  *string
	Symbol        string
	Name          string
	Type          models.AssetType
	Quantity      float64
	PurchasePrice float64
	Currency      string
	PurchaseDate  *time.Time
	Metadata      *string
}

// CreateAsset creates a new asset
func (s *AssetService) CreateAsset(ctx context.Context, input CreateAssetInput) (*models.Asset, error) {
	if input.Currency == "" {
		input.Currency = "USD"
	}

	// Try to get current price
	currentPrice := input.PurchasePrice
	if price, err := s.getPrice(ctx, input.Symbol, string(input.Type)); err == nil {
		currentPrice = price.Price
	}

	asset := &models.Asset{
		UserID:        input.UserID,
		SubAccountID:  input.SubAccountID,
		Symbol:        input.Symbol,
		Name:          input.Name,
		Type:          input.Type,
		Quantity:      input.Quantity,
		PurchasePrice: input.PurchasePrice,
		CurrentPrice:  currentPrice,
		Currency:      input.Currency,
		PurchaseDate:  input.PurchaseDate,
		Metadata:      input.Metadata,
	}

	if err := s.assetRepo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// GetAsset returns an asset by ID
func (s *AssetService) GetAsset(ctx context.Context, id, userID string) (*models.Asset, error) {
	return s.assetRepo.GetByID(ctx, id, userID)
}

// ListAssets lists assets for a user
func (s *AssetService) ListAssets(ctx context.Context, userID, subAccountID string, assetType models.AssetType, page, pageSize int) ([]models.Asset, int64, float64, float64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	assets, total, err := s.assetRepo.List(ctx, userID, subAccountID, assetType, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Calculate totals
	var totalValue, totalProfitLoss float64
	for _, asset := range assets {
		totalValue += asset.TotalValue
		totalProfitLoss += asset.ProfitLoss
	}

	return assets, total, totalValue, totalProfitLoss, nil
}

// UpdateAssetInput holds input for updating an asset
type UpdateAssetInput struct {
	ID            string
	UserID        string
	Quantity      float64
	PurchasePrice float64
	CurrentPrice  float64
	Name          string
	Metadata      *string
}

// UpdateAsset updates an asset
func (s *AssetService) UpdateAsset(ctx context.Context, input UpdateAssetInput) (*models.Asset, error) {
	asset, err := s.assetRepo.GetByID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	if input.Quantity > 0 {
		asset.Quantity = input.Quantity
	}
	if input.PurchasePrice > 0 {
		asset.PurchasePrice = input.PurchasePrice
	}
	if input.CurrentPrice > 0 {
		asset.CurrentPrice = input.CurrentPrice
	}
	if input.Name != "" {
		asset.Name = input.Name
	}
	if input.Metadata != nil && *input.Metadata != "" {
		asset.Metadata = input.Metadata
	}

	if err := s.assetRepo.Update(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// DeleteAsset deletes an asset
func (s *AssetService) DeleteAsset(ctx context.Context, id, userID string) error {
	return s.assetRepo.Delete(ctx, id, userID)
}

// GetAssetPrice gets current price for an asset
func (s *AssetService) GetAssetPrice(ctx context.Context, symbol string, assetType models.AssetType) (*providers.PriceData, error) {
	return s.getPrice(ctx, symbol, string(assetType))
}

// GetMultipleAssetPrices gets prices for multiple assets
func (s *AssetService) GetMultipleAssetPrices(ctx context.Context, queries map[string]string) (map[string]*providers.PriceData, []string, error) {
	results := make(map[string]*providers.PriceData)
	var failedSymbols []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for symbol, assetType := range queries {
		wg.Add(1)
		go func(sym, aType string) {
			defer wg.Done()

			price, err := s.getPrice(ctx, sym, aType)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failedSymbols = append(failedSymbols, sym)
				return
			}
			results[sym] = price
		}(symbol, assetType)
	}

	wg.Wait()
	return results, failedSymbols, nil
}

// RefreshAssetPrices refreshes prices for user's assets
func (s *AssetService) RefreshAssetPrices(ctx context.Context, userID string, assetIDs []string) (int, []string, error) {
	// Get all symbols for user
	symbols, err := s.assetRepo.GetAllSymbols(ctx, userID)
	if err != nil {
		return 0, nil, err
	}

	// Filter by asset IDs if provided
	if len(assetIDs) > 0 {
		filteredSymbols := make(map[string]models.AssetType)
		for _, id := range assetIDs {
			asset, err := s.assetRepo.GetByID(ctx, id, userID)
			if err != nil {
				continue
			}
			filteredSymbols[asset.Symbol] = asset.Type
		}
		symbols = filteredSymbols
	}

	// Fetch prices concurrently
	queries := make(map[string]string)
	for symbol, assetType := range symbols {
		queries[symbol] = string(assetType)
	}

	prices, failedSymbols, err := s.GetMultipleAssetPrices(ctx, queries)
	if err != nil {
		return 0, nil, err
	}

	// Update prices in database
	priceUpdates := make(map[string]float64)
	for symbol, priceData := range prices {
		priceUpdates[symbol] = priceData.Price

		// Update price cache
		_ = s.priceCacheRepo.Upsert(ctx, &models.PriceCache{
			Symbol:           symbol,
			AssetType:        symbols[symbol],
			Price:            priceData.Price,
			Currency:         priceData.Currency,
			Change24h:        priceData.Change24h,
			ChangePercent24h: priceData.ChangePercent24h,
			High24h:          priceData.High24h,
			Low24h:           priceData.Low24h,
			Volume24h:        priceData.Volume24h,
			MarketCap:        priceData.MarketCap,
			UpdatedAt:        time.Now(),
		})
	}

	if err := s.assetRepo.BulkUpdatePrices(ctx, priceUpdates); err != nil {
		return 0, nil, err
	}

	return len(prices), failedSymbols, nil
}

// GetAssetHistory gets historical prices for an asset
func (s *AssetService) GetAssetHistory(ctx context.Context, symbol string, assetType models.AssetType, startDate, endDate time.Time) ([]providers.HistoricalPrice, error) {
	var provider providers.PriceProvider
	switch assetType {
	case models.AssetTypeCrypto:
		provider = s.priceManager.GetCryptoProvider()
	default:
		provider = s.priceManager.GetStockProvider()
	}

	return provider.GetHistoricalPrices(ctx, symbol, startDate, endDate)
}

// GetPortfolioPerformance gets portfolio performance metrics
func (s *AssetService) GetPortfolioPerformance(ctx context.Context, userID string) (*PortfolioPerformance, error) {
	totalValue, err := s.assetRepo.GetTotalValue(ctx, userID)
	if err != nil {
		return nil, err
	}

	valueByType, err := s.assetRepo.GetValueByType(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate allocation percentages
	allocations := make([]AssetAllocation, 0, len(valueByType))
	for assetType, value := range valueByType {
		percentage := 0.0
		if totalValue > 0 {
			percentage = (value / totalValue) * 100
		}
		allocations = append(allocations, AssetAllocation{
			Type:       assetType,
			Value:      value,
			Percentage: percentage,
		})
	}

	// Get all assets to calculate total invested and profit/loss
	assets, _, err := s.assetRepo.List(ctx, userID, "", "", 1, 1000)
	if err != nil {
		return nil, err
	}

	var totalInvested, totalProfitLoss float64
	for _, asset := range assets {
		totalInvested += asset.Quantity * asset.PurchasePrice
		totalProfitLoss += asset.ProfitLoss
	}

	profitLossPercent := 0.0
	if totalInvested > 0 {
		profitLossPercent = (totalProfitLoss / totalInvested) * 100
	}

	return &PortfolioPerformance{
		TotalValue:        totalValue,
		TotalInvested:     totalInvested,
		TotalProfitLoss:   totalProfitLoss,
		ProfitLossPercent: profitLossPercent,
		Allocation:        allocations,
	}, nil
}

// SearchAssets searches for assets
func (s *AssetService) SearchAssets(ctx context.Context, query string, assetType models.AssetType, limit int) ([]providers.SearchResult, error) {
	var provider providers.PriceProvider
	switch assetType {
	case models.AssetTypeCrypto:
		provider = s.priceManager.GetCryptoProvider()
	default:
		provider = s.priceManager.GetStockProvider()
	}

	results, err := provider.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func (s *AssetService) getPrice(ctx context.Context, symbol, assetType string) (*providers.PriceData, error) {
	// Check Redis cache first
	cacheKey := cache.AssetPriceKey(symbol)
	var cachedPrice providers.PriceData
	if err := s.redisCache.Get(ctx, cacheKey, &cachedPrice); err == nil {
		return &cachedPrice, nil
	}

	// Fetch from provider
	price, err := s.priceManager.GetPrice(ctx, symbol, assetType)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.redisCache.Set(ctx, cacheKey, price, priceCacheTTL)

	return price, nil
}

// PortfolioPerformance holds portfolio performance data
type PortfolioPerformance struct {
	TotalValue        float64
	TotalInvested     float64
	TotalProfitLoss   float64
	ProfitLossPercent float64
	Allocation        []AssetAllocation
}

// AssetAllocation holds allocation data for an asset type
type AssetAllocation struct {
	Type       models.AssetType
	Value      float64
	Percentage float64
}

