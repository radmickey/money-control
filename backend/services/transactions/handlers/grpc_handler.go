package handlers

import (
	"context"
	"time"

	pb "github.com/radmickey/money-control/backend/proto/transactions"
	"github.com/radmickey/money-control/backend/services/transactions/models"
	"github.com/radmickey/money-control/backend/services/transactions/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the TransactionsServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedTransactionsServiceServer
	transactionService *service.TransactionService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(transactionService *service.TransactionService) *GRPCHandler {
	return &GRPCHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction creates a new transaction
func (h *GRPCHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.Transaction, error) {
	var date time.Time
	if req.Date != nil {
		date = req.Date.AsTime()
	} else {
		date = time.Now()
	}

	var subAccountID *string
	if req.SubAccountId != "" {
		subAccountID = &req.SubAccountId
	}

	var transferToSubAccount *string
	if req.TransferToSubAccountId != "" {
		transferToSubAccount = &req.TransferToSubAccountId
	}

	tx, err := h.transactionService.CreateTransaction(ctx, service.CreateTransactionInput{
		UserID:               req.UserId,
		SubAccountID:         subAccountID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		Type:                 protoToTransactionType(req.Type),
		Category:             protoToTransactionCategory(req.Category),
		CustomCategory:       req.CustomCategory,
		Description:          req.Description,
		Merchant:             req.Merchant,
		Date:                 date,
		TransferToSubAccount: transferToSubAccount,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	return transactionToProto(tx), nil
}

// GetTransaction gets a transaction by ID
func (h *GRPCHandler) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.Transaction, error) {
	tx, err := h.transactionService.GetTransaction(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transaction not found: %v", err)
	}

	return transactionToProto(tx), nil
}

// ListTransactions lists transactions for a user
func (h *GRPCHandler) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	transactions, total, err := h.transactionService.ListTransactions(ctx, req.UserId, req.SubAccountId, int(req.Page), int(req.PageSize), req.SortBy, req.SortDesc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list transactions: %v", err)
	}

	pbTransactions := make([]*pb.Transaction, len(transactions))
	for i, tx := range transactions {
		pbTransactions[i] = transactionToProto(&tx)
	}

	return &pb.ListTransactionsResponse{
		Transactions: pbTransactions,
		Total:        int32(total),
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

// UpdateTransaction updates a transaction
func (h *GRPCHandler) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.Transaction, error) {
	var date time.Time
	if req.Date != nil {
		date = req.Date.AsTime()
	}

	tx, err := h.transactionService.UpdateTransaction(ctx, service.UpdateTransactionInput{
		ID:             req.Id,
		UserID:         req.UserId,
		Amount:         req.Amount,
		Type:           protoToTransactionType(req.Type),
		Category:       protoToTransactionCategory(req.Category),
		CustomCategory: req.CustomCategory,
		Description:    req.Description,
		Merchant:       req.Merchant,
		Date:           date,
		Currency:       req.Currency,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update transaction: %v", err)
	}

	return transactionToProto(tx), nil
}

// DeleteTransaction deletes a transaction
func (h *GRPCHandler) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*emptypb.Empty, error) {
	err := h.transactionService.DeleteTransaction(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete transaction: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetTransactionsByDateRange gets transactions by date range
func (h *GRPCHandler) GetTransactionsByDateRange(ctx context.Context, req *pb.GetTransactionsByDateRangeRequest) (*pb.ListTransactionsResponse, error) {
	// Use ListTransactions with filters
	transactions, total, err := h.transactionService.ListTransactions(ctx, req.UserId, req.SubAccountId, int(req.Page), int(req.PageSize), "date", true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	pbTransactions := make([]*pb.Transaction, len(transactions))
	for i, tx := range transactions {
		pbTransactions[i] = transactionToProto(&tx)
	}

	return &pb.ListTransactionsResponse{
		Transactions: pbTransactions,
		Total:        int32(total),
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

// GetTransactionsByCategory gets transactions by category
func (h *GRPCHandler) GetTransactionsByCategory(ctx context.Context, req *pb.GetTransactionsByCategoryRequest) (*pb.ListTransactionsResponse, error) {
	// Use ListTransactions - category filtering would need to be added to service
	transactions, total, err := h.transactionService.ListTransactions(ctx, req.UserId, "", int(req.Page), int(req.PageSize), "date", true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	pbTransactions := make([]*pb.Transaction, len(transactions))
	for i, tx := range transactions {
		pbTransactions[i] = transactionToProto(&tx)
	}

	return &pb.ListTransactionsResponse{
		Transactions: pbTransactions,
		Total:        int32(total),
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

// GetTransactionsSummary gets a summary of transactions
func (h *GRPCHandler) GetTransactionsSummary(ctx context.Context, req *pb.GetTransactionsSummaryRequest) (*pb.TransactionsSummaryResponse, error) {
	var startDate, endDate time.Time
	if req.StartDate != nil {
		startDate = req.StartDate.AsTime()
	}
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
	}

	summary, err := h.transactionService.GetTransactionsSummary(ctx, req.UserId, startDate, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get summary: %v", err)
	}

	return &pb.TransactionsSummaryResponse{
		TotalIncome:      summary.TotalIncome,
		TotalExpenses:    summary.TotalExpenses,
		NetFlow:          summary.TotalIncome - summary.TotalExpenses,
		Currency:         req.BaseCurrency,
		ByCategory:       summary.ByCategory,
		TransactionCount: int32(summary.TransactionCount),
	}, nil
}

// CategorizeTransaction categorizes a transaction
func (h *GRPCHandler) CategorizeTransaction(ctx context.Context, req *pb.CategorizeTransactionRequest) (*pb.Transaction, error) {
	tx, err := h.transactionService.CategorizeTransaction(ctx, req.Id, req.UserId, models.TransactionCategory(req.Category.String()), req.CustomCategory)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to categorize transaction: %v", err)
	}

	return transactionToProto(tx), nil
}

// BulkCategorize categorizes multiple transactions
func (h *GRPCHandler) BulkCategorize(ctx context.Context, req *pb.BulkCategorizeRequest) (*pb.BulkCategorizeResponse, error) {
	updated, failedIds, err := h.transactionService.BulkCategorize(ctx, req.UserId, req.TransactionIds, models.TransactionCategory(req.Category.String()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to bulk categorize: %v", err)
	}

	return &pb.BulkCategorizeResponse{
		UpdatedCount: int32(updated),
		FailedIds:    failedIds,
	}, nil
}

// Helper functions
func transactionToProto(t *models.Transaction) *pb.Transaction {
	subAccountID := ""
	if t.SubAccountID != nil {
		subAccountID = *t.SubAccountID
	}

	transferToSubAccount := ""
	if t.TransferToSubAccount != nil {
		transferToSubAccount = *t.TransferToSubAccount
	}

	return &pb.Transaction{
		Id:                     t.ID,
		UserId:                 t.UserID,
		SubAccountId:           subAccountID,
		Amount:                 t.Amount,
		Currency:               t.Currency,
		Type:                   transactionTypeToProto(t.Type),
		Category:               transactionCategoryToProto(t.Category),
		CustomCategory:         t.CustomCategory,
		Description:            t.Description,
		Merchant:               t.Merchant,
		Date:                   timestamppb.New(t.Date),
		CreatedAt:              timestamppb.New(t.CreatedAt),
		UpdatedAt:              timestamppb.New(t.UpdatedAt),
		TransferToSubAccountId: transferToSubAccount,
	}
}

func transactionTypeToProto(t models.TransactionType) pb.TransactionType {
	switch t {
	case models.TransactionTypeIncome:
		return pb.TransactionType_TRANSACTION_TYPE_INCOME
	case models.TransactionTypeExpense:
		return pb.TransactionType_TRANSACTION_TYPE_EXPENSE
	case models.TransactionTypeTransfer:
		return pb.TransactionType_TRANSACTION_TYPE_TRANSFER
	default:
		return pb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}
}

func transactionCategoryToProto(c models.TransactionCategory) pb.TransactionCategory {
	switch c {
	case models.CategorySalary:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_SALARY
	case models.CategoryFood:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_FOOD
	case models.CategoryTransport:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_TRANSPORT
	case models.CategoryShopping:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_SHOPPING
	case models.CategoryEntertainment:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_ENTERTAINMENT
	case models.CategoryHealthcare:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_HEALTHCARE
	case models.CategoryUtilities:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_UTILITIES
	case models.CategoryTransfer:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_TRANSFER
	default:
		return pb.TransactionCategory_TRANSACTION_CATEGORY_OTHER
	}
}

func protoToTransactionType(t pb.TransactionType) models.TransactionType {
	switch t {
	case pb.TransactionType_TRANSACTION_TYPE_INCOME:
		return models.TransactionTypeIncome
	case pb.TransactionType_TRANSACTION_TYPE_EXPENSE:
		return models.TransactionTypeExpense
	case pb.TransactionType_TRANSACTION_TYPE_TRANSFER:
		return models.TransactionTypeTransfer
	default:
		return models.TransactionTypeExpense
	}
}

func protoToTransactionCategory(c pb.TransactionCategory) models.TransactionCategory {
	switch c {
	case pb.TransactionCategory_TRANSACTION_CATEGORY_SALARY:
		return models.CategorySalary
	case pb.TransactionCategory_TRANSACTION_CATEGORY_FOOD:
		return models.CategoryFood
	case pb.TransactionCategory_TRANSACTION_CATEGORY_TRANSPORT:
		return models.CategoryTransport
	case pb.TransactionCategory_TRANSACTION_CATEGORY_SHOPPING:
		return models.CategoryShopping
	case pb.TransactionCategory_TRANSACTION_CATEGORY_ENTERTAINMENT:
		return models.CategoryEntertainment
	case pb.TransactionCategory_TRANSACTION_CATEGORY_HEALTHCARE:
		return models.CategoryHealthcare
	case pb.TransactionCategory_TRANSACTION_CATEGORY_UTILITIES:
		return models.CategoryUtilities
	case pb.TransactionCategory_TRANSACTION_CATEGORY_TRANSFER:
		return models.CategoryTransfer
	default:
		return models.CategoryOther
	}
}
