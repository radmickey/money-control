package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

// TransactionCategory represents the category of transaction
type TransactionCategory string

const (
	CategorySalary           TransactionCategory = "salary"
	CategoryInvestmentIncome TransactionCategory = "investment_income"
	CategoryFood             TransactionCategory = "food"
	CategoryTransport        TransactionCategory = "transport"
	CategoryUtilities        TransactionCategory = "utilities"
	CategoryEntertainment    TransactionCategory = "entertainment"
	CategoryShopping         TransactionCategory = "shopping"
	CategoryHealthcare       TransactionCategory = "healthcare"
	CategoryEducation        TransactionCategory = "education"
	CategoryTravel           TransactionCategory = "travel"
	CategoryHousing          TransactionCategory = "housing"
	CategoryInsurance        TransactionCategory = "insurance"
	CategoryTaxes            TransactionCategory = "taxes"
	CategoryGifts            TransactionCategory = "gifts"
	CategoryTransfer         TransactionCategory = "transfer"
	CategoryOther            TransactionCategory = "other"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID                    string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID                string              `gorm:"type:uuid;not null;index" json:"user_id"`
	SubAccountID          *string             `gorm:"type:uuid;index" json:"sub_account_id,omitempty"`
	Amount                float64             `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency              string              `gorm:"size:3;not null;default:'USD'" json:"currency"`
	Type                  TransactionType     `gorm:"size:20;not null" json:"type"`
	Category              TransactionCategory `gorm:"size:50;not null" json:"category"`
	CustomCategory        string              `gorm:"size:100" json:"custom_category,omitempty"`
	Description           string              `gorm:"size:500" json:"description,omitempty"`
	Merchant              string              `gorm:"size:255" json:"merchant,omitempty"`
	Date                  time.Time           `gorm:"not null;index" json:"date"`
	TransferToSubAccount  *string             `gorm:"type:uuid" json:"transfer_to_sub_account_id,omitempty"`
	Metadata              *string             `gorm:"type:jsonb" json:"metadata,omitempty"`
	IsRecurring           bool                `gorm:"default:false" json:"is_recurring"`
	RecurringFrequency    string              `gorm:"size:20" json:"recurring_frequency,omitempty"`
	IsPending             bool                `gorm:"default:false" json:"is_pending"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
	DeletedAt             gorm.DeletedAt      `gorm:"index" json:"-"`
}

// TableName returns the table name for GORM
func (Transaction) TableName() string {
	return "transactions"
}

// CategoryRule stores rules for auto-categorization
type CategoryRule struct {
	ID        string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string              `gorm:"type:uuid;not null;index" json:"user_id"`
	Pattern   string              `gorm:"size:255;not null" json:"pattern"`
	Category  TransactionCategory `gorm:"size:50;not null" json:"category"`
	IsGlobal  bool                `gorm:"default:false" json:"is_global"`
	Priority  int                 `gorm:"default:0" json:"priority"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// TableName returns the table name for GORM
func (CategoryRule) TableName() string {
	return "category_rules"
}

