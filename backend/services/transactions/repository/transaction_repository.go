package repository

import (
	"context"
	"errors"
	"time"

	"github.com/radmickey/money-control/backend/pkg/database"
	"github.com/radmickey/money-control/backend/services/transactions/models"
	"gorm.io/gorm"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

// GetByID finds a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id, userID string) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return &tx, nil
}

// List lists transactions for a user with pagination
func (r *TransactionRepository) List(ctx context.Context, userID string, subAccountID string, page, pageSize int, sortBy string, sortDesc bool) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).Where("user_id = ?", userID)

	if subAccountID != "" {
		query = query.Where("sub_account_id = ?", subAccountID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort
	if sortBy == "" {
		sortBy = "date"
	}
	order := sortBy
	if sortDesc {
		order += " DESC"
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Order(order).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// ListByDateRange lists transactions within a date range
func (r *TransactionRepository) ListByDateRange(ctx context.Context, userID string, startDate, endDate time.Time, subAccountID string, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate)

	if subAccountID != "" {
		query = query.Where("sub_account_id = ?", subAccountID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// ListByCategory lists transactions by category
func (r *TransactionRepository) ListByCategory(ctx context.Context, userID string, category models.TransactionCategory, startDate, endDate time.Time, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND category = ?", userID, category)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

// Delete soft-deletes a transaction
func (r *TransactionRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{})
	if result.RowsAffected == 0 {
		return ErrTransactionNotFound
	}
	return result.Error
}

// GetSummary gets transaction summary for a user
func (r *TransactionRepository) GetSummary(ctx context.Context, userID string, startDate, endDate time.Time) (*TransactionSummary, error) {
	var summary TransactionSummary

	// Get total income
	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?", userID, models.TransactionTypeIncome, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalIncome).Error; err != nil {
		return nil, err
	}

	// Get total expenses
	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?", userID, models.TransactionTypeExpense, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalExpenses).Error; err != nil {
		return nil, err
	}

	summary.NetFlow = summary.TotalIncome - summary.TotalExpenses

	// Get count
	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Count(&summary.TransactionCount).Error; err != nil {
		return nil, err
	}

	// Get by category
	type CategoryResult struct {
		Category string
		Total    float64
	}

	var categoryResults []CategoryResult
	if err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?", userID, models.TransactionTypeExpense, startDate, endDate).
		Select("category, SUM(amount) as total").
		Group("category").
		Scan(&categoryResults).Error; err != nil {
		return nil, err
	}

	summary.ByCategory = make(map[string]float64)
	for _, cr := range categoryResults {
		summary.ByCategory[cr.Category] = cr.Total
	}

	return &summary, nil
}

// BulkUpdateCategory updates category for multiple transactions
func (r *TransactionRepository) BulkUpdateCategory(ctx context.Context, userID string, ids []string, category models.TransactionCategory) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Update("category", category)
	return result.RowsAffected, result.Error
}

// TransactionSummary holds transaction summary data
type TransactionSummary struct {
	TotalIncome      float64
	TotalExpenses    float64
	NetFlow          float64
	TransactionCount int64
	ByCategory       map[string]float64
}

// CategoryRuleRepository handles category rules
type CategoryRuleRepository struct {
	db *gorm.DB
}

// NewCategoryRuleRepository creates a new category rule repository
func NewCategoryRuleRepository(db *gorm.DB) *CategoryRuleRepository {
	return &CategoryRuleRepository{db: db}
}

// Create creates a category rule
func (r *CategoryRuleRepository) Create(ctx context.Context, rule *models.CategoryRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetByUserID gets category rules for a user
func (r *CategoryRuleRepository) GetByUserID(ctx context.Context, userID string) ([]models.CategoryRule, error) {
	var rules []models.CategoryRule
	if err := r.db.WithContext(ctx).
		Where("user_id = ? OR is_global = true", userID).
		Order("priority DESC").
		Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// Delete deletes a category rule
func (r *CategoryRuleRepository) Delete(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.CategoryRule{}).Error
}

