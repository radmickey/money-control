package models

import (
	"time"

	"gorm.io/gorm"
)

// AssetType represents the type of asset
type AssetType string

const (
	AssetTypeStock      AssetType = "stock"
	AssetTypeCrypto     AssetType = "crypto"
	AssetTypeETF        AssetType = "etf"
	AssetTypeRealEstate AssetType = "real_estate"
	AssetTypeCash       AssetType = "cash"
	AssetTypeBond       AssetType = "bond"
	AssetTypeCommodity  AssetType = "commodity"
	AssetTypeOther      AssetType = "other"
)

// Asset represents a financial asset
type Asset struct {
	ID               string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           string         `gorm:"type:uuid;not null;index" json:"user_id"`
	SubAccountID     *string        `gorm:"type:uuid;index" json:"sub_account_id,omitempty"`
	Symbol           string         `gorm:"size:20;not null;index" json:"symbol"`
	Name             string         `gorm:"size:255" json:"name"`
	Type             AssetType      `gorm:"size:50;not null" json:"type"`
	Quantity         float64        `gorm:"type:decimal(20,8);not null" json:"quantity"`
	PurchasePrice    float64        `gorm:"type:decimal(20,8)" json:"purchase_price"`
	CurrentPrice     float64        `gorm:"type:decimal(20,8)" json:"current_price"`
	Currency         string         `gorm:"size:3;not null;default:'USD'" json:"currency"`
	TotalValue       float64        `gorm:"type:decimal(20,8)" json:"total_value"`
	ProfitLoss       float64        `gorm:"type:decimal(20,8)" json:"profit_loss"`
	ProfitLossPercent float64       `gorm:"type:decimal(10,4)" json:"profit_loss_percent"`
	PurchaseDate     *time.Time     `json:"purchase_date,omitempty"`
	PriceUpdatedAt   *time.Time     `json:"price_updated_at,omitempty"`
	Metadata         *string        `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for GORM
func (Asset) TableName() string {
	return "assets"
}

// CalculateProfitLoss calculates profit/loss for the asset
func (a *Asset) CalculateProfitLoss() {
	if a.PurchasePrice > 0 && a.CurrentPrice > 0 {
		a.TotalValue = a.Quantity * a.CurrentPrice
		invested := a.Quantity * a.PurchasePrice
		a.ProfitLoss = a.TotalValue - invested
		a.ProfitLossPercent = (a.ProfitLoss / invested) * 100
	} else {
		a.TotalValue = a.Quantity * a.CurrentPrice
	}
}

// PriceCache stores cached prices for assets
type PriceCache struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Symbol      string    `gorm:"size:20;uniqueIndex;not null" json:"symbol"`
	AssetType   AssetType `gorm:"size:50;not null" json:"asset_type"`
	Price       float64   `gorm:"type:decimal(20,8);not null" json:"price"`
	Currency    string    `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Change24h   float64   `gorm:"type:decimal(20,8)" json:"change_24h"`
	ChangePercent24h float64 `gorm:"type:decimal(10,4)" json:"change_percent_24h"`
	High24h     float64   `gorm:"type:decimal(20,8)" json:"high_24h"`
	Low24h      float64   `gorm:"type:decimal(20,8)" json:"low_24h"`
	Volume24h   float64   `gorm:"type:decimal(30,8)" json:"volume_24h"`
	MarketCap   float64   `gorm:"type:decimal(30,2)" json:"market_cap"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the table name for GORM
func (PriceCache) TableName() string {
	return "price_cache"
}

// PriceHistory stores historical price data
type PriceHistory struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Symbol    string    `gorm:"size:20;not null;index" json:"symbol"`
	AssetType AssetType `gorm:"size:50;not null" json:"asset_type"`
	Open      float64   `gorm:"type:decimal(20,8)" json:"open"`
	High      float64   `gorm:"type:decimal(20,8)" json:"high"`
	Low       float64   `gorm:"type:decimal(20,8)" json:"low"`
	Close     float64   `gorm:"type:decimal(20,8)" json:"close"`
	Volume    float64   `gorm:"type:decimal(30,8)" json:"volume"`
	Date      time.Time `gorm:"type:date;not null;index" json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (PriceHistory) TableName() string {
	return "price_history"
}

