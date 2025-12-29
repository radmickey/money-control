package repository

import (
	"context"
	"errors"

	"github.com/radmickey/money-control/backend/pkg/database"
	"github.com/radmickey/money-control/backend/services/accounts/models"
	"gorm.io/gorm"
)

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrSubAccountNotFound = errors.New("sub-account not found")
	ErrUnauthorized       = errors.New("unauthorized access to resource")
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetByID finds an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id, userID string) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).Preload("SubAccounts").Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// ListByUserID lists all accounts for a user
func (r *AccountRepository) ListByUserID(ctx context.Context, userID string, accountType models.AccountType, page, pageSize int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Account{}).Where("user_id = ?", userID)

	if accountType != "" {
		query = query.Where("type = ?", accountType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Preload("SubAccounts").Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// Update updates an account
func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// Delete soft-deletes an account
func (r *AccountRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{})
	if result.RowsAffected == 0 {
		return ErrAccountNotFound
	}
	return result.Error
}

// UpdateTotalBalance updates the total balance of an account
func (r *AccountRepository) UpdateTotalBalance(ctx context.Context, accountID string) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE accounts
		SET total_balance = (
			SELECT COALESCE(SUM(balance), 0)
			FROM sub_accounts
			WHERE account_id = ? AND deleted_at IS NULL
		)
		WHERE id = ?
	`, accountID, accountID).Error
}

// SubAccountRepository handles database operations for sub-accounts
type SubAccountRepository struct {
	db *gorm.DB
}

// NewSubAccountRepository creates a new sub-account repository
func NewSubAccountRepository(db *gorm.DB) *SubAccountRepository {
	return &SubAccountRepository{db: db}
}

// Create creates a new sub-account
func (r *SubAccountRepository) Create(ctx context.Context, subAccount *models.SubAccount) error {
	return r.db.WithContext(ctx).Create(subAccount).Error
}

// GetByID finds a sub-account by ID
func (r *SubAccountRepository) GetByID(ctx context.Context, id, userID string) (*models.SubAccount, error) {
	var subAccount models.SubAccount
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&subAccount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubAccountNotFound
		}
		return nil, err
	}
	return &subAccount, nil
}

// ListByAccountID lists all sub-accounts for an account
func (r *SubAccountRepository) ListByAccountID(ctx context.Context, accountID, userID string) ([]models.SubAccount, error) {
	var subAccounts []models.SubAccount
	if err := r.db.WithContext(ctx).Where("account_id = ? AND user_id = ?", accountID, userID).Order("created_at DESC").Find(&subAccounts).Error; err != nil {
		return nil, err
	}
	return subAccounts, nil
}

// ListByUserID lists all sub-accounts for a user
func (r *SubAccountRepository) ListByUserID(ctx context.Context, userID string, assetType models.AssetType) ([]models.SubAccount, error) {
	var subAccounts []models.SubAccount
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}

	if err := query.Order("created_at DESC").Find(&subAccounts).Error; err != nil {
		return nil, err
	}
	return subAccounts, nil
}

// Update updates a sub-account
func (r *SubAccountRepository) Update(ctx context.Context, subAccount *models.SubAccount) error {
	return r.db.WithContext(ctx).Save(subAccount).Error
}

// Delete soft-deletes a sub-account
func (r *SubAccountRepository) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.SubAccount{})
	if result.RowsAffected == 0 {
		return ErrSubAccountNotFound
	}
	return result.Error
}

// UpdateBalance updates the balance of a sub-account
func (r *SubAccountRepository) UpdateBalance(ctx context.Context, id string, balance, quantity float64) error {
	updates := map[string]interface{}{
		"balance": balance,
	}
	if quantity > 0 {
		updates["quantity"] = quantity
	}
	return r.db.WithContext(ctx).Model(&models.SubAccount{}).Where("id = ?", id).Updates(updates).Error
}

// GetTotalBalanceByUser gets total balance across all sub-accounts for a user
func (r *SubAccountRepository) GetTotalBalanceByUser(ctx context.Context, userID string) (map[string]float64, error) {
	type Result struct {
		Currency string
		Total    float64
	}

	var results []Result
	if err := r.db.WithContext(ctx).Model(&models.SubAccount{}).
		Select("currency, SUM(balance) as total").
		Where("user_id = ?", userID).
		Group("currency").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	balances := make(map[string]float64)
	for _, r := range results {
		balances[r.Currency] = r.Total
	}
	return balances, nil
}

// GetBalanceByAssetType gets total balance grouped by asset type
func (r *SubAccountRepository) GetBalanceByAssetType(ctx context.Context, userID string) (map[models.AssetType]float64, error) {
	type Result struct {
		AssetType models.AssetType
		Total     float64
	}

	var results []Result
	if err := r.db.WithContext(ctx).Model(&models.SubAccount{}).
		Select("asset_type, SUM(balance) as total").
		Where("user_id = ?", userID).
		Group("asset_type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	balances := make(map[models.AssetType]float64)
	for _, r := range results {
		balances[r.AssetType] = r.Total
	}
	return balances, nil
}

// BalanceHistoryRepository handles balance history operations
type BalanceHistoryRepository struct {
	db *gorm.DB
}

// NewBalanceHistoryRepository creates a new balance history repository
func NewBalanceHistoryRepository(db *gorm.DB) *BalanceHistoryRepository {
	return &BalanceHistoryRepository{db: db}
}

// Create creates a balance history entry
func (r *BalanceHistoryRepository) Create(ctx context.Context, history *models.BalanceHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetBySubAccountID gets balance history for a sub-account
func (r *BalanceHistoryRepository) GetBySubAccountID(ctx context.Context, subAccountID, userID string, limit int) ([]models.BalanceHistory, error) {
	var history []models.BalanceHistory
	query := r.db.WithContext(ctx).Where("sub_account_id = ? AND user_id = ?", subAccountID, userID).Order("date DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

