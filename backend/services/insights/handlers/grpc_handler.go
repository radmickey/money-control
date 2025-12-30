package handlers

import (
	"context"
	"time"

	pb "github.com/radmickey/money-control/backend/proto/insights"
	"github.com/radmickey/money-control/backend/services/insights/models"
	"github.com/radmickey/money-control/backend/services/insights/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the InsightsServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedInsightsServiceServer
	insightService *service.InsightService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(insightService *service.InsightService) *GRPCHandler {
	return &GRPCHandler{
		insightService: insightService,
	}
}

// GetNetWorth gets the total net worth for a user
func (h *GRPCHandler) GetNetWorth(ctx context.Context, req *pb.GetNetWorthRequest) (*pb.NetWorthResponse, error) {
	netWorth, err := h.insightService.GetNetWorth(ctx, req.UserId, req.BaseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get net worth: %v", err)
	}

	return &pb.NetWorthResponse{
		TotalNetWorth:       netWorth.TotalNetWorth,
		Currency:            netWorth.Currency,
		Change_24H:          netWorth.Change24h,
		ChangePercent_24H:   netWorth.ChangePercent24h,
		Change_7D:           netWorth.Change7d,
		ChangePercent_7D:    netWorth.ChangePercent7d,
		Change_30D:          netWorth.Change30d,
		ChangePercent_30D:   netWorth.ChangePercent30d,
		CalculatedAt:        timestamppb.New(netWorth.CalculatedAt),
	}, nil
}

// GetNetWorthHistory gets net worth history
func (h *GRPCHandler) GetNetWorthHistory(ctx context.Context, req *pb.GetNetWorthHistoryRequest) (*pb.NetWorthHistoryResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	history, err := h.insightService.GetNetWorthHistory(ctx, req.UserId, req.BaseCurrency, req.Period, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get net worth history: %v", err)
	}

	pbHistory := make([]*pb.NetWorthPoint, len(history))
	for i, h := range history {
		pbHistory[i] = &pb.NetWorthPoint{
			Date:  timestamppb.New(h.Date),
			Value: h.Value,
		}
	}

	return &pb.NetWorthHistoryResponse{
		History:  pbHistory,
		Currency: req.BaseCurrency,
	}, nil
}

// GetAllocation gets asset allocation
func (h *GRPCHandler) GetAllocation(ctx context.Context, req *pb.GetAllocationRequest) (*pb.AllocationResponse, error) {
	allocations, totalValue, err := h.insightService.GetAllocation(ctx, req.UserId, req.BaseCurrency, req.GroupBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get allocation: %v", err)
	}

	pbAllocations := make([]*pb.AllocationItem, len(allocations))
	for i, a := range allocations {
		pbAllocations[i] = &pb.AllocationItem{
			Category:   a.Category,
			Name:       a.Name,
			Value:      a.Value,
			Percentage: a.Percentage,
			Color:      a.Color,
		}
	}

	return &pb.AllocationResponse{
		Allocations: pbAllocations,
		TotalValue:  totalValue,
		Currency:    req.BaseCurrency,
	}, nil
}

// GetBalanceChanges gets balance changes
func (h *GRPCHandler) GetBalanceChanges(ctx context.Context, req *pb.GetBalanceChangesRequest) (*pb.BalanceChangesResponse, error) {
	// Simplified - return empty response for now
	return &pb.BalanceChangesResponse{
		Changes:            []*pb.BalanceChange{},
		TotalChange:        0,
		TotalChangePercent: 0,
		Currency:           req.BaseCurrency,
	}, nil
}

// GetTrends gets trends
func (h *GRPCHandler) GetTrends(ctx context.Context, req *pb.GetTrendsRequest) (*pb.TrendsResponse, error) {
	history, err := h.insightService.GetNetWorthHistory(ctx, req.UserId, req.BaseCurrency, req.Period, time.Time{}, time.Time{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get trends: %v", err)
	}

	pbData := make([]*pb.TrendPoint, len(history))
	var sum, min, max float64
	if len(history) > 0 {
		min = history[0].Value
		max = history[0].Value
	}
	for i, h := range history {
		pbData[i] = &pb.TrendPoint{
			Date:  timestamppb.New(h.Date),
			Value: h.Value,
		}
		sum += h.Value
		if h.Value < min {
			min = h.Value
		}
		if h.Value > max {
			max = h.Value
		}
	}

	var avg float64
	if len(history) > 0 {
		avg = sum / float64(len(history))
	}

	trendDirection := "stable"
	if len(history) >= 2 {
		if history[len(history)-1].Value > history[0].Value {
			trendDirection = "up"
		} else if history[len(history)-1].Value < history[0].Value {
			trendDirection = "down"
		}
	}

	return &pb.TrendsResponse{
		Metric:         req.Metric,
		Period:         req.Period,
		Data:           pbData,
		Average:        avg,
		Min:            min,
		Max:            max,
		TrendDirection: trendDirection,
	}, nil
}

// GetDashboardSummary gets dashboard summary
func (h *GRPCHandler) GetDashboardSummary(ctx context.Context, req *pb.GetDashboardSummaryRequest) (*pb.DashboardSummaryResponse, error) {
	summary, err := h.insightService.GetDashboardSummary(ctx, req.UserId, req.BaseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get dashboard summary: %v", err)
	}

	var allocations []*pb.AllocationItem
	for _, a := range summary.TopAllocations {
		allocations = append(allocations, &pb.AllocationItem{
			Category:   a.Category,
			Name:       a.Name,
			Value:      a.Value,
			Percentage: a.Percentage,
			Color:      a.Color,
		})
	}

	var recentTx []*pb.RecentTransaction
	for _, tx := range summary.RecentTransactions {
		recentTx = append(recentTx, &pb.RecentTransaction{
			Id:          tx.ID,
			Description: tx.Description,
			Amount:      tx.Amount,
			Currency:    tx.Currency,
			Category:    tx.Category,
			Date:        timestamppb.New(tx.Date),
		})
	}

	var topPerformers []*pb.TopPerformer
	for _, p := range summary.TopPerformers {
		topPerformers = append(topPerformers, &pb.TopPerformer{
			Symbol:        p.Symbol,
			Name:          p.Name,
			Value:         p.Value,
			Change:        p.Change,
			ChangePercent: p.ChangePercent,
		})
	}

	return &pb.DashboardSummaryResponse{
		NetWorth:                   summary.NetWorth,
		NetWorthChange_24H:         summary.NetWorthChange24h,
		NetWorthChangePercent_24H:  summary.NetWorthChangePercent,
		TotalAssets:                summary.TotalAssets,
		TotalLiabilities:           summary.TotalLiabilities,
		MonthlyIncome:              summary.MonthlyIncome,
		MonthlyExpenses:            summary.MonthlyExpenses,
		MonthlySavings:             summary.MonthlySavings,
		SavingsRate:                summary.SavingsRate,
		TopAllocations:             allocations,
		RecentTransactions:         recentTx,
		TopPerformers:              topPerformers,
		Currency:                   req.BaseCurrency,
		CalculatedAt:               timestamppb.Now(),
	}, nil
}

// GetCashFlow gets cash flow
func (h *GRPCHandler) GetCashFlow(ctx context.Context, req *pb.GetCashFlowRequest) (*pb.CashFlowResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	periods, totalIncome, totalExpenses, netCashFlow, err := h.insightService.GetCashFlow(ctx, req.UserId, req.BaseCurrency, req.Period, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cash flow: %v", err)
	}

	pbPeriods := make([]*pb.CashFlowPeriod, len(periods))
	for i, p := range periods {
		pbPeriods[i] = &pb.CashFlowPeriod{
			Label:     p.Label,
			StartDate: timestamppb.New(p.StartDate),
			EndDate:   timestamppb.New(p.EndDate),
			Income:    p.Income,
			Expenses:  p.Expenses,
			Net:       p.Net,
		}
	}

	return &pb.CashFlowResponse{
		Periods:       pbPeriods,
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		NetCashFlow:   netCashFlow,
		Currency:      req.BaseCurrency,
	}, nil
}

// CreateSnapshot creates a daily snapshot
func (h *GRPCHandler) CreateSnapshot(ctx context.Context, req *pb.CreateSnapshotRequest) (*pb.Snapshot, error) {
	snapshot, err := h.insightService.CreateSnapshot(ctx, req.UserId, req.BaseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create snapshot: %v", err)
	}

	return snapshotToProto(snapshot), nil
}

// GetSnapshots gets snapshots
func (h *GRPCHandler) GetSnapshots(ctx context.Context, req *pb.GetSnapshotsRequest) (*pb.GetSnapshotsResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	snapshots, err := h.insightService.GetSnapshots(ctx, req.UserId, startDate, endDate, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get snapshots: %v", err)
	}

	pbSnapshots := make([]*pb.Snapshot, len(snapshots))
	for i, s := range snapshots {
		pbSnapshots[i] = snapshotToProto(&s)
	}

	return &pb.GetSnapshotsResponse{
		Snapshots: pbSnapshots,
	}, nil
}

// Helper functions
func snapshotToProto(s *models.Snapshot) *pb.Snapshot {
	return &pb.Snapshot{
		Id:            s.ID,
		UserId:        s.UserID,
		TotalNetWorth: s.TotalNetWorth,
		Currency:      s.Currency,
		AllocationData: s.AllocationData,
		Date:          timestamppb.New(s.Date),
		CreatedAt:     timestamppb.New(s.CreatedAt),
	}
}
