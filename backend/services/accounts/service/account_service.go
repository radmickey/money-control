package service

import (
	"context"
	"time"

	"github.com/radmickey/money-control/backend/services/accounts/models"
	"github.com/radmickey/money-control/backend/services/accounts/repository"
)

// AccountService handles account business logic
type AccountService struct {
	accountRepo        *repository.AccountRepository
	subAccountRepo     *repository.SubAccountRepository
	balanceHistoryRepo *repository.BalanceHistoryRepository
}

// NewAccountService creates a new account service
func NewAccountService(
	accountRepo *repository.AccountRepository,
	subAccountRepo *repository.SubAccountRepository,
	balanceHistoryRepo *repository.BalanceHistoryRepository,
) *AccountService {
	return &AccountService{
		accountRepo:        accountRepo,
		subAccountRepo:     subAccountRepo,
		balanceHistoryRepo: balanceHistoryRepo,
	}
}

// CreateAccountInput holds input for creating an account
type CreateAccountInput struct {
	UserID      string
	Name        string
	Type        models.AccountType
	Currency    string
	Description string
	Icon        string
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(ctx context.Context, input CreateAccountInput) (*models.Account, error) {
	if input.Currency == "" {
		input.Currency = "USD"
	}

	account := &models.Account{
		UserID:      input.UserID,
		Name:        input.Name,
		Type:        input.Type,
		Currency:    input.Currency,
		Description: input.Description,
		Icon:        input.Icon,
		IsActive:    true,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccount returns an account by ID
func (s *AccountService) GetAccount(ctx context.Context, id, userID string) (*models.Account, error) {
	return s.accountRepo.GetByID(ctx, id, userID)
}

// ListAccounts lists accounts for a user
func (s *AccountService) ListAccounts(ctx context.Context, userID string, accountType models.AccountType, page, pageSize int) ([]models.Account, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.accountRepo.ListByUserID(ctx, userID, accountType, page, pageSize)
}

// UpdateAccountInput holds input for updating an account
type UpdateAccountInput struct {
	ID          string
	UserID      string
	Name        string
	Description string
	Icon        string
	Currency    string
}

// UpdateAccount updates an account
func (s *AccountService) UpdateAccount(ctx context.Context, input UpdateAccountInput) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		account.Name = input.Name
	}
	if input.Description != "" {
		account.Description = input.Description
	}
	if input.Icon != "" {
		account.Icon = input.Icon
	}
	if input.Currency != "" {
		account.Currency = input.Currency
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// DeleteAccount deletes an account
func (s *AccountService) DeleteAccount(ctx context.Context, id, userID string) error {
	return s.accountRepo.Delete(ctx, id, userID)
}

// CreateSubAccountInput holds input for creating a sub-account
type CreateSubAccountInput struct {
	AccountID   string
	UserID      string
	Name        string
	AssetType   models.AssetType
	Currency    string
	Balance     float64
	Symbol      string
	Quantity    float64
	Description string
}

// CreateSubAccount creates a new sub-account
func (s *AccountService) CreateSubAccount(ctx context.Context, input CreateSubAccountInput) (*models.SubAccount, error) {
	// Verify account exists and belongs to user
	_, err := s.accountRepo.GetByID(ctx, input.AccountID, input.UserID)
	if err != nil {
		return nil, err
	}

	if input.Currency == "" {
		input.Currency = "USD"
	}

	subAccount := &models.SubAccount{
		AccountID:   input.AccountID,
		UserID:      input.UserID,
		Name:        input.Name,
		AssetType:   input.AssetType,
		Currency:    input.Currency,
		Balance:     input.Balance,
		Symbol:      input.Symbol,
		Quantity:    input.Quantity,
		Description: input.Description,
		IsActive:    true,
	}

	if err := s.subAccountRepo.Create(ctx, subAccount); err != nil {
		return nil, err
	}

	// Update account total balance
	_ = s.accountRepo.UpdateTotalBalance(ctx, input.AccountID)

	// Record initial balance in history
	if input.Balance > 0 {
		_ = s.balanceHistoryRepo.Create(ctx, &models.BalanceHistory{
			SubAccountID: subAccount.ID,
			UserID:       input.UserID,
			Balance:      input.Balance,
			Date:         time.Now(),
		})
	}

	return subAccount, nil
}

// GetSubAccount returns a sub-account by ID
func (s *AccountService) GetSubAccount(ctx context.Context, id, userID string) (*models.SubAccount, error) {
	return s.subAccountRepo.GetByID(ctx, id, userID)
}

// ListSubAccounts lists sub-accounts for an account
func (s *AccountService) ListSubAccounts(ctx context.Context, accountID, userID string, assetType models.AssetType) ([]models.SubAccount, error) {
	if accountID != "" {
		return s.subAccountRepo.ListByAccountID(ctx, accountID, userID)
	}
	return s.subAccountRepo.ListByUserID(ctx, userID, assetType)
}

// UpdateSubAccountInput holds input for updating a sub-account
type UpdateSubAccountInput struct {
	ID          string
	UserID      string
	Name        string
	Description string
	Balance     float64
	Quantity    float64
}

// UpdateSubAccount updates a sub-account
func (s *AccountService) UpdateSubAccount(ctx context.Context, input UpdateSubAccountInput) (*models.SubAccount, error) {
	subAccount, err := s.subAccountRepo.GetByID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		subAccount.Name = input.Name
	}
	if input.Description != "" {
		subAccount.Description = input.Description
	}

	oldBalance := subAccount.Balance
	if input.Balance != 0 {
		subAccount.Balance = input.Balance
	}
	if input.Quantity != 0 {
		subAccount.Quantity = input.Quantity
	}

	if err := s.subAccountRepo.Update(ctx, subAccount); err != nil {
		return nil, err
	}

	// Update account total balance
	_ = s.accountRepo.UpdateTotalBalance(ctx, subAccount.AccountID)

	// Record balance change in history if changed
	if subAccount.Balance != oldBalance {
		_ = s.balanceHistoryRepo.Create(ctx, &models.BalanceHistory{
			SubAccountID: subAccount.ID,
			UserID:       input.UserID,
			Balance:      subAccount.Balance,
			Date:         time.Now(),
		})
	}

	return subAccount, nil
}

// DeleteSubAccount deletes a sub-account
func (s *AccountService) DeleteSubAccount(ctx context.Context, id, userID string) error {
	subAccount, err := s.subAccountRepo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := s.subAccountRepo.Delete(ctx, id, userID); err != nil {
		return err
	}

	// Update account total balance
	_ = s.accountRepo.UpdateTotalBalance(ctx, subAccount.AccountID)

	return nil
}

// UpdateSubAccountBalance updates balance of a sub-account
func (s *AccountService) UpdateSubAccountBalance(ctx context.Context, id, userID string, balance, quantity float64) (*models.SubAccount, error) {
	subAccount, err := s.subAccountRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	oldBalance := subAccount.Balance
	subAccount.Balance = balance
	if quantity > 0 {
		subAccount.Quantity = quantity
	}

	if err := s.subAccountRepo.Update(ctx, subAccount); err != nil {
		return nil, err
	}

	// Update account total balance
	_ = s.accountRepo.UpdateTotalBalance(ctx, subAccount.AccountID)

	// Record balance change in history if changed
	if balance != oldBalance {
		_ = s.balanceHistoryRepo.Create(ctx, &models.BalanceHistory{
			SubAccountID: subAccount.ID,
			UserID:       userID,
			Balance:      balance,
			Date:         time.Now(),
		})
	}

	return subAccount, nil
}

// GetUserNetWorth calculates total net worth for a user
func (s *AccountService) GetUserNetWorth(ctx context.Context, userID string) (map[string]float64, error) {
	return s.subAccountRepo.GetTotalBalanceByUser(ctx, userID)
}

// GetBalanceByAssetType gets balances grouped by asset type
func (s *AccountService) GetBalanceByAssetType(ctx context.Context, userID string) (map[models.AssetType]float64, error) {
	return s.subAccountRepo.GetBalanceByAssetType(ctx, userID)
}

// GetAccountsSummary gets summary of accounts for a user
func (s *AccountService) GetAccountsSummary(ctx context.Context, userID string) (*AccountsSummary, error) {
	accounts, total, err := s.accountRepo.ListByUserID(ctx, userID, "", 1, 1000)
	if err != nil {
		return nil, err
	}

	summary := &AccountsSummary{
		TotalAccounts: int(total),
		ByType:        make(map[models.AccountType]TypeSummary),
	}

	for _, account := range accounts {
		summary.TotalSubAccounts += len(account.SubAccounts)
		typeSummary := summary.ByType[account.Type]
		typeSummary.Count++
		typeSummary.TotalBalance += account.TotalBalance
		summary.ByType[account.Type] = typeSummary
	}

	return summary, nil
}

// AccountsSummary holds summary information about accounts
type AccountsSummary struct {
	TotalAccounts    int
	TotalSubAccounts int
	ByType           map[models.AccountType]TypeSummary
}

// TypeSummary holds summary for an account type
type TypeSummary struct {
	Count        int
	TotalBalance float64
}

