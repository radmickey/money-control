package service

import (
	"context"
	"strings"
	"time"

	"github.com/radmickey/money-control/backend/services/transactions/models"
	"github.com/radmickey/money-control/backend/services/transactions/repository"
)

// TransactionService handles transaction business logic
type TransactionService struct {
	txRepo       *repository.TransactionRepository
	ruleRepo     *repository.CategoryRuleRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	txRepo *repository.TransactionRepository,
	ruleRepo *repository.CategoryRuleRepository,
) *TransactionService {
	return &TransactionService{
		txRepo:   txRepo,
		ruleRepo: ruleRepo,
	}
}

// CreateTransactionInput holds input for creating a transaction
type CreateTransactionInput struct {
	UserID               string
	SubAccountID         *string
	Amount               float64
	Currency             string
	Type                 models.TransactionType
	Category             models.TransactionCategory
	CustomCategory       string
	Description          string
	Merchant             string
	Date                 time.Time
	TransferToSubAccount *string
	Metadata             *string
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*models.Transaction, error) {
	// Auto-categorize if no category provided
	if input.Category == "" {
		input.Category = s.autoCategorize(ctx, input.UserID, input.Description, input.Merchant)
	}

	if input.Currency == "" {
		input.Currency = "USD"
	}

	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	tx := &models.Transaction{
		UserID:               input.UserID,
		SubAccountID:         input.SubAccountID,
		Amount:               input.Amount,
		Currency:             input.Currency,
		Type:                 input.Type,
		Category:             input.Category,
		CustomCategory:       input.CustomCategory,
		Description:          input.Description,
		Merchant:             input.Merchant,
		Date:                 input.Date,
		TransferToSubAccount: input.TransferToSubAccount,
		Metadata:             input.Metadata,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransaction returns a transaction by ID
func (s *TransactionService) GetTransaction(ctx context.Context, id, userID string) (*models.Transaction, error) {
	return s.txRepo.GetByID(ctx, id, userID)
}

// ListTransactions lists transactions for a user
func (s *TransactionService) ListTransactions(ctx context.Context, userID, subAccountID string, page, pageSize int, sortBy string, sortDesc bool) ([]models.Transaction, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.txRepo.List(ctx, userID, subAccountID, page, pageSize, sortBy, sortDesc)
}

// ListTransactionsByDateRange lists transactions in a date range
func (s *TransactionService) ListTransactionsByDateRange(ctx context.Context, userID string, startDate, endDate time.Time, subAccountID string, page, pageSize int) ([]models.Transaction, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.txRepo.ListByDateRange(ctx, userID, startDate, endDate, subAccountID, page, pageSize)
}

// ListTransactionsByCategory lists transactions by category
func (s *TransactionService) ListTransactionsByCategory(ctx context.Context, userID string, category models.TransactionCategory, startDate, endDate time.Time, page, pageSize int) ([]models.Transaction, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.txRepo.ListByCategory(ctx, userID, category, startDate, endDate, page, pageSize)
}

// UpdateTransactionInput holds input for updating a transaction
type UpdateTransactionInput struct {
	ID             string
	UserID         string
	Amount         float64
	Type           models.TransactionType
	Category       models.TransactionCategory
	CustomCategory string
	Description    string
	Merchant       string
	Date           time.Time
	Currency       string
}

// UpdateTransaction updates a transaction
func (s *TransactionService) UpdateTransaction(ctx context.Context, input UpdateTransactionInput) (*models.Transaction, error) {
	tx, err := s.txRepo.GetByID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	if input.Amount != 0 {
		tx.Amount = input.Amount
	}
	if input.Type != "" {
		tx.Type = input.Type
	}
	if input.Category != "" {
		tx.Category = input.Category
	}
	if input.CustomCategory != "" {
		tx.CustomCategory = input.CustomCategory
	}
	if input.Description != "" {
		tx.Description = input.Description
	}
	if input.Merchant != "" {
		tx.Merchant = input.Merchant
	}
	if !input.Date.IsZero() {
		tx.Date = input.Date
	}
	if input.Currency != "" {
		tx.Currency = input.Currency
	}

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// DeleteTransaction deletes a transaction
func (s *TransactionService) DeleteTransaction(ctx context.Context, id, userID string) error {
	return s.txRepo.Delete(ctx, id, userID)
}

// GetTransactionsSummary gets transaction summary
func (s *TransactionService) GetTransactionsSummary(ctx context.Context, userID string, startDate, endDate time.Time) (*repository.TransactionSummary, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	return s.txRepo.GetSummary(ctx, userID, startDate, endDate)
}

// CategorizeTransaction updates the category of a transaction
func (s *TransactionService) CategorizeTransaction(ctx context.Context, id, userID string, category models.TransactionCategory, customCategory string) (*models.Transaction, error) {
	tx, err := s.txRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	tx.Category = category
	if customCategory != "" {
		tx.CustomCategory = customCategory
	}

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// BulkCategorize updates category for multiple transactions
func (s *TransactionService) BulkCategorize(ctx context.Context, userID string, ids []string, category models.TransactionCategory) (int64, []string, error) {
	updated, err := s.txRepo.BulkUpdateCategory(ctx, userID, ids, category)
	if err != nil {
		return 0, nil, err
	}
	return updated, nil, nil
}

// autoCategorize attempts to categorize a transaction based on rules
// This is a stub implementation - can be extended with ML/AI later
func (s *TransactionService) autoCategorize(ctx context.Context, userID, description, merchant string) models.TransactionCategory {
	// Get user rules
	rules, err := s.ruleRepo.GetByUserID(ctx, userID)
	if err != nil {
		return models.CategoryOther
	}

	text := strings.ToLower(description + " " + merchant)

	// Check rules
	for _, rule := range rules {
		if strings.Contains(text, strings.ToLower(rule.Pattern)) {
			return rule.Category
		}
	}

	// Default keyword matching
	keywordCategories := map[string]models.TransactionCategory{
		"grocery":     models.CategoryFood,
		"restaurant":  models.CategoryFood,
		"uber eats":   models.CategoryFood,
		"doordash":    models.CategoryFood,
		"uber":        models.CategoryTransport,
		"lyft":        models.CategoryTransport,
		"gas":         models.CategoryTransport,
		"netflix":     models.CategoryEntertainment,
		"spotify":     models.CategoryEntertainment,
		"amazon":      models.CategoryShopping,
		"walmart":     models.CategoryShopping,
		"pharmacy":    models.CategoryHealthcare,
		"doctor":      models.CategoryHealthcare,
		"hospital":    models.CategoryHealthcare,
		"electric":    models.CategoryUtilities,
		"water":       models.CategoryUtilities,
		"internet":    models.CategoryUtilities,
		"rent":        models.CategoryHousing,
		"mortgage":    models.CategoryHousing,
		"insurance":   models.CategoryInsurance,
		"salary":      models.CategorySalary,
		"payroll":     models.CategorySalary,
		"dividend":    models.CategoryInvestmentIncome,
		"interest":    models.CategoryInvestmentIncome,
		"tax":         models.CategoryTaxes,
	}

	for keyword, category := range keywordCategories {
		if strings.Contains(text, keyword) {
			return category
		}
	}

	return models.CategoryOther
}

// CreateCategoryRule creates a categorization rule
func (s *TransactionService) CreateCategoryRule(ctx context.Context, userID, pattern string, category models.TransactionCategory, priority int) (*models.CategoryRule, error) {
	rule := &models.CategoryRule{
		UserID:   userID,
		Pattern:  pattern,
		Category: category,
		Priority: priority,
	}

	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// GetCategoryRules gets category rules for a user
func (s *TransactionService) GetCategoryRules(ctx context.Context, userID string) ([]models.CategoryRule, error) {
	return s.ruleRepo.GetByUserID(ctx, userID)
}

// DeleteCategoryRule deletes a category rule
func (s *TransactionService) DeleteCategoryRule(ctx context.Context, id, userID string) error {
	return s.ruleRepo.Delete(ctx, id, userID)
}

