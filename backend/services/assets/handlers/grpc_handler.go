package handlers

import (
	"context"
	"time"

	pb "github.com/radmickey/money-control/backend/proto/assets"
	"github.com/radmickey/money-control/backend/services/assets/models"
	"github.com/radmickey/money-control/backend/services/assets/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the AssetsServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedAssetsServiceServer
	assetService *service.AssetService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(assetService *service.AssetService) *GRPCHandler {
	return &GRPCHandler{
		assetService: assetService,
	}
}

// CreateAsset creates a new asset
func (h *GRPCHandler) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.Asset, error) {
	var subAccountID *string
	if req.SubAccountId != "" {
		subAccountID = &req.SubAccountId
	}

	asset, err := h.assetService.CreateAsset(ctx, service.CreateAssetInput{
		UserID:        req.UserId,
		SubAccountID:  subAccountID,
		Symbol:        req.Symbol,
		Name:          req.Name,
		Type:          models.AssetType(req.Type.String()),
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		Currency:      req.Currency,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create asset: %v", err)
	}

	return assetToProto(asset), nil
}

// GetAsset gets an asset by ID
func (h *GRPCHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.Asset, error) {
	asset, err := h.assetService.GetAsset(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "asset not found: %v", err)
	}

	return assetToProto(asset), nil
}

// ListAssets lists assets for a user
func (h *GRPCHandler) ListAssets(ctx context.Context, req *pb.ListAssetsRequest) (*pb.ListAssetsResponse, error) {
	assets, total, totalValue, totalProfitLoss, err := h.assetService.ListAssets(ctx, req.UserId, req.SubAccountId, models.AssetType(req.Type.String()), int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list assets: %v", err)
	}

	pbAssets := make([]*pb.Asset, len(assets))
	for i, a := range assets {
		pbAssets[i] = assetToProto(&a)
	}

	return &pb.ListAssetsResponse{
		Assets:          pbAssets,
		Total:           int32(total),
		Page:            req.Page,
		PageSize:        req.PageSize,
		TotalValue:      totalValue,
		TotalProfitLoss: totalProfitLoss,
	}, nil
}

// UpdateAsset updates an asset
func (h *GRPCHandler) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.Asset, error) {
	asset, err := h.assetService.UpdateAsset(ctx, service.UpdateAssetInput{
		ID:            req.Id,
		UserID:        req.UserId,
		Quantity:      req.Quantity,
		PurchasePrice: req.PurchasePrice,
		CurrentPrice:  req.CurrentPrice,
		Name:          req.Name,
		Metadata:      nil,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update asset: %v", err)
	}

	return assetToProto(asset), nil
}

// DeleteAsset deletes an asset
func (h *GRPCHandler) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*emptypb.Empty, error) {
	err := h.assetService.DeleteAsset(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete asset: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetAssetPrice gets the current price of an asset
func (h *GRPCHandler) GetAssetPrice(ctx context.Context, req *pb.GetAssetPriceRequest) (*pb.AssetPriceResponse, error) {
	price, err := h.assetService.GetAssetPrice(ctx, req.Symbol, models.AssetType(req.Type.String()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get price: %v", err)
	}

	return &pb.AssetPriceResponse{
		Symbol:            price.Symbol,
		Price:             price.Price,
		Currency:          price.Currency,
		Change_24H:        price.Change24h,
		ChangePercent_24H: price.ChangePercent24h,
		UpdatedAt:         timestamppb.Now(),
	}, nil
}

// GetMultipleAssetPrices gets prices for multiple assets
func (h *GRPCHandler) GetMultipleAssetPrices(ctx context.Context, req *pb.GetMultipleAssetPricesRequest) (*pb.GetMultipleAssetPricesResponse, error) {
	queries := make(map[string]string)
	for _, q := range req.Queries {
		queries[q.Symbol] = q.Type.String()
	}

	prices, failedSymbols, err := h.assetService.GetMultipleAssetPrices(ctx, queries)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get prices: %v", err)
	}

	pbPrices := make(map[string]*pb.AssetPriceResponse)
	for symbol, price := range prices {
		pbPrices[symbol] = &pb.AssetPriceResponse{
			Symbol:            price.Symbol,
			Price:             price.Price,
			Currency:          price.Currency,
			Change_24H:        price.Change24h,
			ChangePercent_24H: price.ChangePercent24h,
			UpdatedAt:         timestamppb.Now(),
		}
	}

	return &pb.GetMultipleAssetPricesResponse{
		Prices:        pbPrices,
		FailedSymbols: failedSymbols,
	}, nil
}

// RefreshAssetPrices refreshes prices for user's assets
func (h *GRPCHandler) RefreshAssetPrices(ctx context.Context, req *pb.RefreshAssetPricesRequest) (*pb.RefreshAssetPricesResponse, error) {
	updatedCount, failedIds, err := h.assetService.RefreshAssetPrices(ctx, req.UserId, req.AssetIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh prices: %v", err)
	}

	return &pb.RefreshAssetPricesResponse{
		UpdatedCount: int32(updatedCount),
		FailedIds:    failedIds,
	}, nil
}

// GetAssetHistory gets price history for an asset
func (h *GRPCHandler) GetAssetHistory(ctx context.Context, req *pb.GetAssetHistoryRequest) (*pb.AssetHistoryResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	history, err := h.assetService.GetAssetHistory(ctx, req.Symbol, models.AssetType(req.Type.String()), startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history: %v", err)
	}

	pbHistory := make([]*pb.PricePoint, len(history))
	for i, hp := range history {
		pbHistory[i] = &pb.PricePoint{
			Timestamp: timestamppb.New(hp.Date),
			Open:      hp.Open,
			High:      hp.High,
			Low:       hp.Low,
			Close:     hp.Close,
			Volume:    hp.Volume,
		}
	}

	return &pb.AssetHistoryResponse{
		Symbol:  req.Symbol,
		History: pbHistory,
	}, nil
}

// GetPortfolioPerformance gets portfolio performance
func (h *GRPCHandler) GetPortfolioPerformance(ctx context.Context, req *pb.GetPortfolioPerformanceRequest) (*pb.PortfolioPerformanceResponse, error) {
	perf, err := h.assetService.GetPortfolioPerformance(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get performance: %v", err)
	}

	var allocations []*pb.AssetAllocation
	for _, a := range perf.Allocation {
		allocations = append(allocations, &pb.AssetAllocation{
			Type:       assetTypeToProto(a.Type),
			Value:      a.Value,
			Percentage: a.Percentage,
		})
	}

	return &pb.PortfolioPerformanceResponse{
		TotalValue:        perf.TotalValue,
		TotalInvested:     perf.TotalInvested,
		TotalProfitLoss:   perf.TotalProfitLoss,
		ProfitLossPercent: perf.ProfitLossPercent,
		Currency:          req.BaseCurrency,
		Allocation:        allocations,
	}, nil
}

// SearchAssets searches for assets
func (h *GRPCHandler) SearchAssets(ctx context.Context, req *pb.SearchAssetsRequest) (*pb.SearchAssetsResponse, error) {
	results, err := h.assetService.SearchAssets(ctx, req.Query, models.AssetType(req.Type.String()), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search assets: %v", err)
	}

	pbResults := make([]*pb.AssetSearchResult, len(results))
	for i, r := range results {
		pbResults[i] = &pb.AssetSearchResult{
			Symbol:   r.Symbol,
			Name:     r.Name,
			Type:     assetTypeToProto(models.AssetType(r.Type)),
			Exchange: r.Exchange,
			Currency: r.Currency,
		}
	}

	return &pb.SearchAssetsResponse{
		Results: pbResults,
	}, nil
}

// Helper functions
func assetToProto(a *models.Asset) *pb.Asset {
	subAccountID := ""
	if a.SubAccountID != nil {
		subAccountID = *a.SubAccountID
	}

	asset := &pb.Asset{
		Id:                a.ID,
		UserId:            a.UserID,
		SubAccountId:      subAccountID,
		Symbol:            a.Symbol,
		Name:              a.Name,
		Type:              assetTypeToProto(a.Type),
		Quantity:          a.Quantity,
		PurchasePrice:     a.PurchasePrice,
		CurrentPrice:      a.CurrentPrice,
		Currency:          a.Currency,
		TotalValue:        a.TotalValue,
		ProfitLoss:        a.ProfitLoss,
		ProfitLossPercent: a.ProfitLossPercent,
		CreatedAt:         timestamppb.New(a.CreatedAt),
		UpdatedAt:         timestamppb.New(a.UpdatedAt),
	}
	if a.PurchaseDate != nil {
		asset.PurchaseDate = timestamppb.New(*a.PurchaseDate)
	}
	if a.PriceUpdatedAt != nil {
		asset.PriceUpdatedAt = timestamppb.New(*a.PriceUpdatedAt)
	}
	return asset
}

func assetTypeToProto(t models.AssetType) pb.AssetType {
	switch t {
	case models.AssetTypeStock:
		return pb.AssetType_ASSET_TYPE_STOCK
	case models.AssetTypeCrypto:
		return pb.AssetType_ASSET_TYPE_CRYPTO
	case models.AssetTypeETF:
		return pb.AssetType_ASSET_TYPE_ETF
	case models.AssetTypeRealEstate:
		return pb.AssetType_ASSET_TYPE_REAL_ESTATE
	case models.AssetTypeCash:
		return pb.AssetType_ASSET_TYPE_CASH
	case models.AssetTypeBond:
		return pb.AssetType_ASSET_TYPE_BOND
	default:
		return pb.AssetType_ASSET_TYPE_OTHER
	}
}
