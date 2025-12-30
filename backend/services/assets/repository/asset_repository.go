package repository

import (
	"context"
	"errors"

	"github.com/radmickey/money-control/backend/pkg/database"
	"github.com/radmickey/money-control/backend/services/assets/models"
	"gorm.io/gorm"
)

var (
	ErrAssetNotFound = errors.New("asset not found")
)

// AssetRepository handles database operations for assets
type AssetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

// Create creates a new asset
func (r *AssetRepository) Create(ctx context.Context, asset *models.Asset) error {
	asset.CalculateProfitLoss()
	return r.db.WithContext(ctx).Create(asset).Error
}

// GetByID finds an asset by ID
func (r *AssetRepository) GetByID(ctx context.Context, id, userID string) (*models.Asset, error) {
	var asset models.Asset
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&asset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}
	return &asset, nil
}

// List lists assets for a user with filters
func (r *AssetRepository) List(ctx context.Context, userID string, subAccountID string, assetType models.AssetType, page, pageSize int) ([]models.Asset, int64, error) {
	var assets []models.Asset
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Asset{}).Where("user_id = ?", userID)

	if subAccountID != "" {
		query = query.Where("sub_account_id = ?", subAccountID)
	}
	if assetType != "" {
		query = query.Where("type = ?", assetType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Order("total_value DESC").Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

// GetBySymbol gets assets by symbol for a user
func (r *AssetRepository) GetBySymbol(ctx context.Context, userID, symbol string) ([]models.Asset, error) {
	var assets []models.Asset
	if err := r.db.WithContext(ctx).Where("user_id = ? AND symbol = ?", userID, symbol).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

// GetAllSymbols gets all unique symbols for a user
func (r *AssetRepository) GetAllSymbols(ctx context.Context, userID string) (map[string]models.AssetType, error) {
	type Result struct {
		Symbol string
		Type   models.AssetType
	}

	var results []Result
	if err := r.db.WithContext(ctx).Model(&models.Asset{}).
		Select("DISTINCT symbol, type").
		Where("user_id = ?", userID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	symbols := make(map[string]models.AssetType)
	for _, r := range results {
		symbols[r.Symbol] = r.Type
	}
	return symbols, nil
}

// Update updates an asset
func (r *AssetRepository) Update(ctx context.Context, asset *models.Asset) error {
	asset.CalculateProfitLoss()
	return r.db.WithContext(ctx).Save(asset).Error
}

// UpdatePrice updates the current price and recalculates values
func (r *AssetRepository) UpdatePrice(ctx context.Context, id string, price float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var asset models.Asset
		if err := tx.First(&asset, "id = ?", id).Error; err != nil {
			return err
		}

		asset.CurrentPrice = price
		asset.CalculateProfitLoss()

		return tx.Save(&asset).Error
	})
}

// BulkUpdatePrices updates prices for multiple assets by symbol
func (r *AssetRepository) BulkUpdatePrices(ctx context.Context, prices map[string]float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for symbol, price := range prices {
			if err := tx.Model(&models.Asset{}).
				Where("symbol = ?", symbol).
				Updates(map[string]interface{}{
					"current_price": price,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete soft-deletes an asset
func (r *AssetRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Asset{})
	if result.RowsAffected == 0 {
		return ErrAssetNotFound
	}
	return result.Error
}

// GetTotalValue gets total portfolio value for a user
func (r *AssetRepository) GetTotalValue(ctx context.Context, userID string) (float64, error) {
	var total float64
	if err := r.db.WithContext(ctx).Model(&models.Asset{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(total_value), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetValueByType gets total value grouped by asset type
func (r *AssetRepository) GetValueByType(ctx context.Context, userID string) (map[models.AssetType]float64, error) {
	type Result struct {
		Type  models.AssetType
		Total float64
	}

	var results []Result
	if err := r.db.WithContext(ctx).Model(&models.Asset{}).
		Select("type, SUM(total_value) as total").
		Where("user_id = ?", userID).
		Group("type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	values := make(map[models.AssetType]float64)
	for _, r := range results {
		values[r.Type] = r.Total
	}
	return values, nil
}

// PriceCacheRepository handles price cache operations
type PriceCacheRepository struct {
	db *gorm.DB
}

// NewPriceCacheRepository creates a new price cache repository
func NewPriceCacheRepository(db *gorm.DB) *PriceCacheRepository {
	return &PriceCacheRepository{db: db}
}

// Upsert inserts or updates a price cache entry
func (r *PriceCacheRepository) Upsert(ctx context.Context, cache *models.PriceCache) error {
	return r.db.WithContext(ctx).
		Where("symbol = ?", cache.Symbol).
		Assign(cache).
		FirstOrCreate(cache).Error
}

// Get gets a cached price
func (r *PriceCacheRepository) Get(ctx context.Context, symbol string) (*models.PriceCache, error) {
	var cache models.PriceCache
	if err := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&cache).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cache, nil
}

// GetMultiple gets multiple cached prices
func (r *PriceCacheRepository) GetMultiple(ctx context.Context, symbols []string) (map[string]*models.PriceCache, error) {
	var caches []models.PriceCache
	if err := r.db.WithContext(ctx).Where("symbol IN ?", symbols).Find(&caches).Error; err != nil {
		return nil, err
	}

	result := make(map[string]*models.PriceCache)
	for i := range caches {
		result[caches[i].Symbol] = &caches[i]
	}
	return result, nil
}

