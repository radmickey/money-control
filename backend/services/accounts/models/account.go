package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeBank       AccountType = "bank"
	AccountTypeCash       AccountType = "cash"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeCrypto     AccountType = "crypto"
	AccountTypeRealEstate AccountType = "real_estate"
	AccountTypeOther      AccountType = "other"
)

// AssetType represents the type of asset in a sub-account
type AssetType string

const (
	AssetTypeBank       AssetType = "bank"
	AssetTypeCash       AssetType = "cash"
	AssetTypeStocks     AssetType = "stocks"
	AssetTypeCrypto     AssetType = "crypto"
	AssetTypeETF        AssetType = "etf"
	AssetTypeRealEstate AssetType = "real_estate"
	AssetTypeBonds      AssetType = "bonds"
	AssetTypeOther      AssetType = "other"
)

// Account represents a main financial account
type Account struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Type         AccountType    `gorm:"size:50;not null" json:"type"`
	Currency     string         `gorm:"size:3;not null;default:'USD'" json:"currency"`
	TotalBalance float64        `gorm:"type:decimal(20,8);default:0" json:"total_balance"`
	Description  string         `gorm:"size:500" json:"description,omitempty"`
	Icon         string         `gorm:"size:50" json:"icon,omitempty"`
	Color        string         `gorm:"size:7" json:"color,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	SubAccounts []SubAccount `gorm:"foreignKey:AccountID" json:"sub_accounts,omitempty"`
}

// TableName returns the table name for GORM
func (Account) TableName() string {
	return "accounts"
}

// SubAccount represents a sub-account within a main account
type SubAccount struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID   string         `gorm:"type:uuid;not null;index" json:"account_id"`
	UserID      string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	AssetType   AssetType      `gorm:"size:50;not null" json:"asset_type"`
	Currency    string         `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Balance     float64        `gorm:"type:decimal(20,8);default:0" json:"balance"`
	Symbol      string         `gorm:"size:20" json:"symbol,omitempty"`
	Quantity    float64        `gorm:"type:decimal(20,8);default:0" json:"quantity"`
	Description string         `gorm:"size:500" json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

// TableName returns the table name for GORM
func (SubAccount) TableName() string {
	return "sub_accounts"
}

// BalanceHistory tracks historical balance changes
type BalanceHistory struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SubAccountID string    `gorm:"type:uuid;not null;index" json:"sub_account_id"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Balance      float64   `gorm:"type:decimal(20,8);not null" json:"balance"`
	Date         time.Time `gorm:"type:date;not null;index" json:"date"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (BalanceHistory) TableName() string {
	return "balance_history"
}

// RecurringIncome represents a scheduled recurring income
type RecurringIncome struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string         `gorm:"type:uuid;not null;index" json:"user_id"`
	SubAccountID string         `gorm:"type:uuid;not null;index" json:"sub_account_id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Amount       float64        `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency     string         `gorm:"size:3;not null;default:'USD'" json:"currency"`
	DayOfMonth   int            `gorm:"not null" json:"day_of_month"` // 1-31, if day > month days, use last day
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastApplied  *time.Time     `json:"last_applied,omitempty"`
	NextApply    *time.Time     `json:"next_apply,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	SubAccount SubAccount `gorm:"foreignKey:SubAccountID" json:"-"`
}

// TableName returns the table name for GORM
func (RecurringIncome) TableName() string {
	return "recurring_incomes"
}

// CalculateNextApplyDate calculates the next date when this income should be applied
func (r *RecurringIncome) CalculateNextApplyDate() time.Time {
	now := time.Now()
	year, month, _ := now.Date()

	// Get the last day of current month
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, now.Location()).Day()

	// Use the smaller of dayOfMonth or lastDayOfMonth
	day := r.DayOfMonth
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}

	nextDate := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	// If the date has passed this month, move to next month
	if nextDate.Before(now) || nextDate.Equal(now) {
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		lastDayOfNextMonth := time.Date(nextYear, nextMonth+1, 0, 0, 0, 0, 0, now.Location()).Day()
		day = r.DayOfMonth
		if day > lastDayOfNextMonth {
			day = lastDayOfNextMonth
		}
		nextDate = time.Date(nextYear, nextMonth, day, 0, 0, 0, 0, now.Location())
	}

	return nextDate
}

