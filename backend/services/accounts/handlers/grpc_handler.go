package handlers

import (
	"context"

	pb "github.com/radmickey/money-control/backend/proto/accounts"
	"github.com/radmickey/money-control/backend/services/accounts/models"
	"github.com/radmickey/money-control/backend/services/accounts/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the AccountsServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedAccountsServiceServer
	accountService *service.AccountService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(accountService *service.AccountService) *GRPCHandler {
	return &GRPCHandler{
		accountService: accountService,
	}
}

// CreateAccount creates a new account
func (h *GRPCHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	account, err := h.accountService.CreateAccount(ctx, service.CreateAccountInput{
		UserID:      req.UserId,
		Name:        req.Name,
		Type:        protoToAccountType(req.Type),
		Currency:    req.Currency,
		Description: req.Description,
		Icon:        req.Icon,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	return accountToProto(account), nil
}

// GetAccount gets an account by ID
func (h *GRPCHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	account, err := h.accountService.GetAccount(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
	}

	return accountToProto(account), nil
}

// ListAccounts lists accounts for a user
func (h *GRPCHandler) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	// Don't filter by type if UNSPECIFIED
	accountType := models.AccountType("")
	if req.Type != pb.AccountType_ACCOUNT_TYPE_UNSPECIFIED {
		accountType = protoToAccountType(req.Type)
	}
	accounts, total, err := h.accountService.ListAccounts(ctx, req.UserId, accountType, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list accounts: %v", err)
	}

	pbAccounts := make([]*pb.Account, len(accounts))
	for i, acc := range accounts {
		pbAccounts[i] = accountToProto(&acc)
	}

	return &pb.ListAccountsResponse{
		Accounts: pbAccounts,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateAccount updates an account
func (h *GRPCHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.Account, error) {
	account, err := h.accountService.UpdateAccount(ctx, service.UpdateAccountInput{
		ID:          req.Id,
		UserID:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Currency:    req.Currency,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update account: %v", err)
	}

	return accountToProto(account), nil
}

// DeleteAccount deletes an account
func (h *GRPCHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	err := h.accountService.DeleteAccount(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete account: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// CreateSubAccount creates a new sub-account
func (h *GRPCHandler) CreateSubAccount(ctx context.Context, req *pb.CreateSubAccountRequest) (*pb.SubAccount, error) {
	subAccount, err := h.accountService.CreateSubAccount(ctx, service.CreateSubAccountInput{
		AccountID:   req.AccountId,
		UserID:      req.UserId,
		Name:        req.Name,
		AssetType:   protoToAssetType(req.AssetType),
		Currency:    req.Currency,
		Balance:     req.Balance,
		Symbol:      req.Symbol,
		Quantity:    req.Quantity,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sub-account: %v", err)
	}

	return subAccountToProto(subAccount), nil
}

// GetSubAccount gets a sub-account by ID
func (h *GRPCHandler) GetSubAccount(ctx context.Context, req *pb.GetSubAccountRequest) (*pb.SubAccount, error) {
	subAccount, err := h.accountService.GetSubAccount(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sub-account not found: %v", err)
	}

	return subAccountToProto(subAccount), nil
}

// ListSubAccounts lists sub-accounts
func (h *GRPCHandler) ListSubAccounts(ctx context.Context, req *pb.ListSubAccountsRequest) (*pb.ListSubAccountsResponse, error) {
	assetType := protoToAssetType(req.AssetType)
	subAccounts, err := h.accountService.ListSubAccounts(ctx, req.AccountId, req.UserId, assetType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list sub-accounts: %v", err)
	}

	pbSubAccounts := make([]*pb.SubAccount, len(subAccounts))
	for i, sa := range subAccounts {
		pbSubAccounts[i] = subAccountToProto(&sa)
	}

	return &pb.ListSubAccountsResponse{
		SubAccounts: pbSubAccounts,
	}, nil
}

// UpdateSubAccount updates a sub-account
func (h *GRPCHandler) UpdateSubAccount(ctx context.Context, req *pb.UpdateSubAccountRequest) (*pb.SubAccount, error) {
	subAccount, err := h.accountService.UpdateSubAccount(ctx, service.UpdateSubAccountInput{
		ID:          req.Id,
		UserID:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		Balance:     req.Balance,
		Quantity:    req.Quantity,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sub-account: %v", err)
	}

	return subAccountToProto(subAccount), nil
}

// DeleteSubAccount deletes a sub-account
func (h *GRPCHandler) DeleteSubAccount(ctx context.Context, req *pb.DeleteSubAccountRequest) (*emptypb.Empty, error) {
	err := h.accountService.DeleteSubAccount(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete sub-account: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// UpdateSubAccountBalance updates balance of a sub-account
func (h *GRPCHandler) UpdateSubAccountBalance(ctx context.Context, req *pb.UpdateSubAccountBalanceRequest) (*pb.SubAccount, error) {
	subAccount, err := h.accountService.UpdateSubAccountBalance(ctx, req.Id, req.UserId, req.Balance, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sub-account balance: %v", err)
	}
	return subAccountToProto(subAccount), nil
}

// GetUserNetWorth gets total net worth for a user
func (h *GRPCHandler) GetUserNetWorth(ctx context.Context, req *pb.GetUserNetWorthRequest) (*pb.NetWorthResponse, error) {
	netWorthByCurrency, err := h.accountService.GetUserNetWorth(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get net worth: %v", err)
	}

	totalNetWorth := netWorthByCurrency[req.BaseCurrency]
	if totalNetWorth == 0 {
		// Sum all if base currency not found
		for _, v := range netWorthByCurrency {
			totalNetWorth += v
		}
	}

	byAssetType, err := h.accountService.GetBalanceByAssetType(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get allocation: %v", err)
	}

	pbByAssetType := make(map[string]float64)
	for k, v := range byAssetType {
		pbByAssetType[string(k)] = v
	}

	return &pb.NetWorthResponse{
		TotalNetWorth: totalNetWorth,
		Currency:      req.BaseCurrency,
		ByAssetType:   pbByAssetType,
		CalculatedAt:  timestamppb.Now(),
	}, nil
}

// GetAccountsSummary gets accounts summary
func (h *GRPCHandler) GetAccountsSummary(ctx context.Context, req *pb.GetAccountsSummaryRequest) (*pb.AccountsSummaryResponse, error) {
	accounts, totalAccounts, err := h.accountService.ListAccounts(ctx, req.UserId, "", 1, 1000)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accounts summary: %v", err)
	}

	// Count by type
	typeCount := make(map[pb.AccountType]int32)
	typeBalance := make(map[pb.AccountType]float64)
	var totalSubAccounts int32

	for _, acc := range accounts {
		at := accountTypeToProto(acc.Type)
		typeCount[at]++
		typeBalance[at] += acc.TotalBalance
		totalSubAccounts += int32(len(acc.SubAccounts))
	}

	var byType []*pb.AccountTypeSummary
	for t, c := range typeCount {
		byType = append(byType, &pb.AccountTypeSummary{
			Type:         t,
			Count:        c,
			TotalBalance: typeBalance[t],
		})
	}

	return &pb.AccountsSummaryResponse{
		TotalAccounts:    int32(totalAccounts),
		TotalSubAccounts: totalSubAccounts,
		ByType:           byType,
	}, nil
}

// Helper functions
func accountToProto(a *models.Account) *pb.Account {
	pbSubAccounts := make([]*pb.SubAccount, len(a.SubAccounts))
	for i, sa := range a.SubAccounts {
		pbSubAccounts[i] = subAccountToProto(&sa)
	}

	return &pb.Account{
		Id:          a.ID,
		UserId:      a.UserID,
		Name:        a.Name,
		Type:        accountTypeToProto(a.Type),
		Currency:    a.Currency,
		TotalBalance: a.TotalBalance,
		Description: a.Description,
		Icon:        a.Icon,
		CreatedAt:   timestamppb.New(a.CreatedAt),
		UpdatedAt:   timestamppb.New(a.UpdatedAt),
		SubAccounts: pbSubAccounts,
	}
}

func subAccountToProto(s *models.SubAccount) *pb.SubAccount {
	return &pb.SubAccount{
		Id:          s.ID,
		AccountId:   s.AccountID,
		Name:        s.Name,
		AssetType:   assetTypeToProto(s.AssetType),
		Currency:    s.Currency,
		Balance:     s.Balance,
		Symbol:      s.Symbol,
		Quantity:    s.Quantity,
		Description: s.Description,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
	}
}

func accountTypeToProto(t models.AccountType) pb.AccountType {
	switch t {
	case models.AccountTypeBank:
		return pb.AccountType_ACCOUNT_TYPE_BANK
	case models.AccountTypeCash:
		return pb.AccountType_ACCOUNT_TYPE_CASH
	case models.AccountTypeInvestment:
		return pb.AccountType_ACCOUNT_TYPE_INVESTMENT
	case models.AccountTypeCrypto:
		return pb.AccountType_ACCOUNT_TYPE_CRYPTO
	case models.AccountTypeRealEstate:
		return pb.AccountType_ACCOUNT_TYPE_REAL_ESTATE
	default:
		return pb.AccountType_ACCOUNT_TYPE_OTHER
	}
}

func assetTypeToProto(t models.AssetType) pb.AssetType {
	switch t {
	case models.AssetTypeBank:
		return pb.AssetType_ASSET_TYPE_BANK
	case models.AssetTypeCash:
		return pb.AssetType_ASSET_TYPE_CASH
	case models.AssetTypeStocks:
		return pb.AssetType_ASSET_TYPE_STOCKS
	case models.AssetTypeCrypto:
		return pb.AssetType_ASSET_TYPE_CRYPTO
	case models.AssetTypeETF:
		return pb.AssetType_ASSET_TYPE_ETF
	case models.AssetTypeRealEstate:
		return pb.AssetType_ASSET_TYPE_REAL_ESTATE
	case models.AssetTypeBonds:
		return pb.AssetType_ASSET_TYPE_BONDS
	default:
		return pb.AssetType_ASSET_TYPE_OTHER
	}
}

// protoToAccountType converts protobuf AccountType to models.AccountType
func protoToAccountType(t pb.AccountType) models.AccountType {
	switch t {
	case pb.AccountType_ACCOUNT_TYPE_BANK:
		return models.AccountTypeBank
	case pb.AccountType_ACCOUNT_TYPE_CASH:
		return models.AccountTypeCash
	case pb.AccountType_ACCOUNT_TYPE_INVESTMENT:
		return models.AccountTypeInvestment
	case pb.AccountType_ACCOUNT_TYPE_CRYPTO:
		return models.AccountTypeCrypto
	case pb.AccountType_ACCOUNT_TYPE_REAL_ESTATE:
		return models.AccountTypeRealEstate
	case pb.AccountType_ACCOUNT_TYPE_OTHER:
		return models.AccountTypeOther
	default:
		return models.AccountTypeOther
	}
}

// protoToAssetType converts protobuf AssetType to models.AssetType
func protoToAssetType(t pb.AssetType) models.AssetType {
	switch t {
	case pb.AssetType_ASSET_TYPE_BANK:
		return models.AssetTypeBank
	case pb.AssetType_ASSET_TYPE_CASH:
		return models.AssetTypeCash
	case pb.AssetType_ASSET_TYPE_STOCKS:
		return models.AssetTypeStocks
	case pb.AssetType_ASSET_TYPE_CRYPTO:
		return models.AssetTypeCrypto
	case pb.AssetType_ASSET_TYPE_ETF:
		return models.AssetTypeETF
	case pb.AssetType_ASSET_TYPE_REAL_ESTATE:
		return models.AssetTypeRealEstate
	case pb.AssetType_ASSET_TYPE_BONDS:
		return models.AssetTypeBonds
	default:
		return models.AssetTypeOther
	}
}
