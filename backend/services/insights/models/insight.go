package models

import (
	"time"
)

// Snapshot represents a point-in-time financial snapshot
type Snapshot struct {
	ID             string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         string    `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalNetWorth  float64   `gorm:"type:decimal(20,8);not null" json:"total_net_worth"`
	Currency       string    `gorm:"size:3;not null;default:'USD'" json:"currency"`
	AllocationData string    `gorm:"type:jsonb" json:"allocation_data"`
	Date           time.Time `gorm:"type:date;not null;index" json:"date"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (Snapshot) TableName() string {
	return "snapshots"
}

// AllocationItem represents an allocation breakdown item
type AllocationItem struct {
	Category   string  `json:"category"`
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color,omitempty"`
}

// NetWorthData represents net worth response data
type NetWorthData struct {
	TotalNetWorth       float64   `json:"total_net_worth"`
	Currency            string    `json:"currency"`
	Change24h           float64   `json:"change_24h"`
	ChangePercent24h    float64   `json:"change_percent_24h"`
	Change7d            float64   `json:"change_7d"`
	ChangePercent7d     float64   `json:"change_percent_7d"`
	Change30d           float64   `json:"change_30d"`
	ChangePercent30d    float64   `json:"change_percent_30d"`
	CalculatedAt        time.Time `json:"calculated_at"`
}

// CashFlowPeriod represents cash flow for a period
type CashFlowPeriod struct {
	Label     string    `json:"label"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Income    float64   `json:"income"`
	Expenses  float64   `json:"expenses"`
	Net       float64   `json:"net"`
}

// TrendPoint represents a data point in a trend
type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
	Label string    `json:"label,omitempty"`
}

// BalanceChange represents a balance change for an account
type BalanceChange struct {
	AccountID       string  `json:"account_id"`
	AccountName     string  `json:"account_name"`
	StartingBalance float64 `json:"starting_balance"`
	EndingBalance   float64 `json:"ending_balance"`
	Change          float64 `json:"change"`
	ChangePercent   float64 `json:"change_percent"`
}

// DashboardSummary represents the dashboard summary
type DashboardSummary struct {
	NetWorth               float64          `json:"net_worth"`
	NetWorthChange24h      float64          `json:"net_worth_change_24h"`
	NetWorthChangePercent  float64          `json:"net_worth_change_percent_24h"`
	TotalAssets            float64          `json:"total_assets"`
	TotalLiabilities       float64          `json:"total_liabilities"`
	MonthlyIncome          float64          `json:"monthly_income"`
	MonthlyExpenses        float64          `json:"monthly_expenses"`
	MonthlySavings         float64          `json:"monthly_savings"`
	SavingsRate            float64          `json:"savings_rate"`
	TopAllocations         []AllocationItem `json:"top_allocations"`
	RecentTransactions     []RecentTx       `json:"recent_transactions"`
	TopPerformers          []TopPerformer   `json:"top_performers"`
	Currency               string           `json:"currency"`
	CalculatedAt           time.Time        `json:"calculated_at"`
}

// RecentTx represents a recent transaction for dashboard
type RecentTx struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// TopPerformer represents a top performing asset
type TopPerformer struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
}

