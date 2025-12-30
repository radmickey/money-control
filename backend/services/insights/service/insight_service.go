package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/radmickey/money-control/backend/services/insights/models"
	"github.com/radmickey/money-control/backend/services/insights/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	assetspb "github.com/radmickey/money-control/backend/proto/assets"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
	transactionspb "github.com/radmickey/money-control/backend/proto/transactions"
)

// ServiceClients holds gRPC clients for other services
type ServiceClients struct {
	AccountsClient     accountspb.AccountsServiceClient
	AssetsClient       assetspb.AssetsServiceClient
	TransactionsClient transactionspb.TransactionsServiceClient
	CurrencyClient     currencypb.CurrencyServiceClient
}

// InsightService handles insights business logic
type InsightService struct {
	snapshotRepo *repository.SnapshotRepository
	clients      *ServiceClients
}

// NewInsightService creates a new insight service
func NewInsightService(
	snapshotRepo *repository.SnapshotRepository,
	accountsURL, assetsURL, transactionsURL, currencyURL string,
) (*InsightService, error) {
	clients := &ServiceClients{}

	// Connect to Accounts service
	if accountsURL != "" {
		conn, err := grpc.Dial(accountsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			clients.AccountsClient = accountspb.NewAccountsServiceClient(conn)
		}
	}

	// Connect to Assets service
	if assetsURL != "" {
		conn, err := grpc.Dial(assetsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			clients.AssetsClient = assetspb.NewAssetsServiceClient(conn)
		}
	}

	// Connect to Transactions service
	if transactionsURL != "" {
		conn, err := grpc.Dial(transactionsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			clients.TransactionsClient = transactionspb.NewTransactionsServiceClient(conn)
		}
	}

	// Connect to Currency service
	if currencyURL != "" {
		conn, err := grpc.Dial(currencyURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			clients.CurrencyClient = currencypb.NewCurrencyServiceClient(conn)
		}
	}

	return &InsightService{
		snapshotRepo: snapshotRepo,
		clients:      clients,
	}, nil
}

// GetNetWorth calculates current net worth
func (s *InsightService) GetNetWorth(ctx context.Context, userID, baseCurrency string) (*models.NetWorthData, error) {
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	// Get current net worth from accounts service
	var totalNetWorth float64
	if s.clients.AccountsClient != nil {
		resp, err := s.clients.AccountsClient.GetUserNetWorth(ctx, &accountspb.GetUserNetWorthRequest{
			UserId:       userID,
			BaseCurrency: baseCurrency,
		})
		if err == nil {
			totalNetWorth = resp.TotalNetWorth
		}
	}

	// Get historical data for comparison
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	var change24h, change7d, change30d float64

	snapshot1d, _ := s.snapshotRepo.GetByDate(ctx, userID, yesterday)
	if snapshot1d != nil {
		change24h = totalNetWorth - snapshot1d.TotalNetWorth
	}

	snapshot7d, _ := s.snapshotRepo.GetByDate(ctx, userID, weekAgo)
	if snapshot7d != nil {
		change7d = totalNetWorth - snapshot7d.TotalNetWorth
	}

	snapshot30d, _ := s.snapshotRepo.GetByDate(ctx, userID, monthAgo)
	if snapshot30d != nil {
		change30d = totalNetWorth - snapshot30d.TotalNetWorth
	}

	return &models.NetWorthData{
		TotalNetWorth:    totalNetWorth,
		Currency:         baseCurrency,
		Change24h:        change24h,
		ChangePercent24h: calculatePercent(change24h, totalNetWorth-change24h),
		Change7d:         change7d,
		ChangePercent7d:  calculatePercent(change7d, totalNetWorth-change7d),
		Change30d:        change30d,
		ChangePercent30d: calculatePercent(change30d, totalNetWorth-change30d),
		CalculatedAt:     now,
	}, nil
}

// GetNetWorthHistory gets net worth history
func (s *InsightService) GetNetWorthHistory(ctx context.Context, userID, baseCurrency, period string, startDate, endDate time.Time) ([]models.TrendPoint, error) {
	// Calculate date range based on period if not provided
	now := time.Now()
	if startDate.IsZero() {
		switch period {
		case "1w":
			startDate = now.AddDate(0, 0, -7)
		case "1m":
			startDate = now.AddDate(0, -1, 0)
		case "3m":
			startDate = now.AddDate(0, -3, 0)
		case "6m":
			startDate = now.AddDate(0, -6, 0)
		case "1y":
			startDate = now.AddDate(-1, 0, 0)
		default:
			startDate = now.AddDate(0, -1, 0)
		}
	}
	if endDate.IsZero() {
		endDate = now
	}

	return s.snapshotRepo.GetNetWorthHistory(ctx, userID, startDate, endDate)
}

// GetAllocation gets asset allocation breakdown
func (s *InsightService) GetAllocation(ctx context.Context, userID, baseCurrency, groupBy string) ([]models.AllocationItem, float64, error) {
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	var allocations []models.AllocationItem
	var totalValue float64

	// Get allocation from accounts service
	if s.clients.AccountsClient != nil {
		resp, err := s.clients.AccountsClient.GetUserNetWorth(ctx, &accountspb.GetUserNetWorthRequest{
			UserId:       userID,
			BaseCurrency: baseCurrency,
		})
		if err == nil {
			totalValue = resp.TotalNetWorth

			// Group by asset type
			for assetType, value := range resp.ByAssetType {
				percentage := 0.0
				if totalValue > 0 {
					percentage = (value / totalValue) * 100
				}
				allocations = append(allocations, models.AllocationItem{
					Category:   assetType,
					Name:       assetType,
					Value:      value,
					Percentage: percentage,
					Color:      getColorForCategory(assetType),
				})
			}
		}
	}

	return allocations, totalValue, nil
}

// GetDashboardSummary gets comprehensive dashboard summary
func (s *InsightService) GetDashboardSummary(ctx context.Context, userID, baseCurrency string) (*models.DashboardSummary, error) {
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	summary := &models.DashboardSummary{
		Currency:     baseCurrency,
		CalculatedAt: time.Now(),
	}

	// Get net worth
	netWorth, err := s.GetNetWorth(ctx, userID, baseCurrency)
	if err == nil {
		summary.NetWorth = netWorth.TotalNetWorth
		summary.NetWorthChange24h = netWorth.Change24h
		summary.NetWorthChangePercent = netWorth.ChangePercent24h
	}

	// Get allocations
	allocations, _, _ := s.GetAllocation(ctx, userID, baseCurrency, "asset_type")
	if len(allocations) > 5 {
		allocations = allocations[:5]
	}
	summary.TopAllocations = allocations

	// Get monthly transactions summary
	if s.clients.TransactionsClient != nil {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		resp, err := s.clients.TransactionsClient.GetTransactionsSummary(ctx, &transactionspb.GetTransactionsSummaryRequest{
			UserId:       userID,
			StartDate:    timestampProto(startOfMonth),
			EndDate:      timestampProto(now),
			BaseCurrency: baseCurrency,
		})
		if err == nil {
			summary.MonthlyIncome = resp.TotalIncome
			summary.MonthlyExpenses = resp.TotalExpenses
			summary.MonthlySavings = resp.TotalIncome - resp.TotalExpenses
			if resp.TotalIncome > 0 {
				summary.SavingsRate = (summary.MonthlySavings / resp.TotalIncome) * 100
			}
		}
	}

	return summary, nil
}

// GetCashFlow gets cash flow analysis
func (s *InsightService) GetCashFlow(ctx context.Context, userID, baseCurrency, period string, startDate, endDate time.Time) ([]models.CashFlowPeriod, float64, float64, float64, error) {
	// Implementation depends on transaction data
	// This is a simplified version
	now := time.Now()
	if startDate.IsZero() {
		startDate = now.AddDate(0, -6, 0) // Last 6 months
	}
	if endDate.IsZero() {
		endDate = now
	}

	var periods []models.CashFlowPeriod
	var totalIncome, totalExpenses, netCashFlow float64

	// Generate monthly periods
	current := startDate
	for current.Before(endDate) {
		monthEnd := current.AddDate(0, 1, 0)
		if monthEnd.After(endDate) {
			monthEnd = endDate
		}

		period := models.CashFlowPeriod{
			Label:     current.Format("Jan 2006"),
			StartDate: current,
			EndDate:   monthEnd,
		}

		// Get transactions for this period
		if s.clients.TransactionsClient != nil {
			resp, err := s.clients.TransactionsClient.GetTransactionsSummary(ctx, &transactionspb.GetTransactionsSummaryRequest{
				UserId:       userID,
				StartDate:    timestampProto(current),
				EndDate:      timestampProto(monthEnd),
				BaseCurrency: baseCurrency,
			})
			if err == nil {
				period.Income = resp.TotalIncome
				period.Expenses = resp.TotalExpenses
				period.Net = resp.TotalIncome - resp.TotalExpenses
				totalIncome += resp.TotalIncome
				totalExpenses += resp.TotalExpenses
			}
		}

		periods = append(periods, period)
		current = monthEnd
	}

	netCashFlow = totalIncome - totalExpenses
	return periods, totalIncome, totalExpenses, netCashFlow, nil
}

// CreateSnapshot creates a daily snapshot
func (s *InsightService) CreateSnapshot(ctx context.Context, userID, baseCurrency string) (*models.Snapshot, error) {
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	// Get current net worth
	netWorth, err := s.GetNetWorth(ctx, userID, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Get allocation
	allocations, _, _ := s.GetAllocation(ctx, userID, baseCurrency, "asset_type")
	allocationJSON, _ := json.Marshal(allocations)

	snapshot := &models.Snapshot{
		UserID:         userID,
		TotalNetWorth:  netWorth.TotalNetWorth,
		Currency:       baseCurrency,
		AllocationData: string(allocationJSON),
		Date:           time.Now().Truncate(24 * time.Hour),
	}

	if err := s.snapshotRepo.Upsert(ctx, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// GetSnapshots gets historical snapshots
func (s *InsightService) GetSnapshots(ctx context.Context, userID string, startDate, endDate time.Time, limit int) ([]models.Snapshot, error) {
	return s.snapshotRepo.GetByDateRange(ctx, userID, startDate, endDate, limit)
}

// Helper functions
func calculatePercent(change, base float64) float64 {
	if base == 0 {
		return 0
	}
	return (change / base) * 100
}

func getColorForCategory(category string) string {
	colors := map[string]string{
		"stocks":      "#4CAF50",
		"crypto":      "#FF9800",
		"bank":        "#2196F3",
		"cash":        "#9C27B0",
		"etf":         "#00BCD4",
		"real_estate": "#795548",
		"bonds":       "#607D8B",
		"other":       "#9E9E9E",
	}
	if color, ok := colors[category]; ok {
		return color
	}
	return "#9E9E9E"
}

func timestampProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
