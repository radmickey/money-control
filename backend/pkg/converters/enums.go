package converters

import (
	"strings"
	"time"

	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	transactionspb "github.com/radmickey/money-control/backend/proto/transactions"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Account type conversions
var accountTypeMap = map[string]accountspb.AccountType{
	"bank":        accountspb.AccountType_ACCOUNT_TYPE_BANK,
	"cash":        accountspb.AccountType_ACCOUNT_TYPE_CASH,
	"investment":  accountspb.AccountType_ACCOUNT_TYPE_INVESTMENT,
	"crypto":      accountspb.AccountType_ACCOUNT_TYPE_CRYPTO,
	"real_estate": accountspb.AccountType_ACCOUNT_TYPE_REAL_ESTATE,
	"other":       accountspb.AccountType_ACCOUNT_TYPE_OTHER,
}

var accountTypeReverseMap = map[accountspb.AccountType]string{
	accountspb.AccountType_ACCOUNT_TYPE_BANK:        "bank",
	accountspb.AccountType_ACCOUNT_TYPE_CASH:        "cash",
	accountspb.AccountType_ACCOUNT_TYPE_INVESTMENT:  "investment",
	accountspb.AccountType_ACCOUNT_TYPE_CRYPTO:      "crypto",
	accountspb.AccountType_ACCOUNT_TYPE_REAL_ESTATE: "real_estate",
	accountspb.AccountType_ACCOUNT_TYPE_OTHER:       "other",
}

func StringToAccountType(s string) accountspb.AccountType {
	if t, ok := accountTypeMap[strings.ToLower(s)]; ok {
		return t
	}
	return accountspb.AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

func AccountTypeToString(t accountspb.AccountType) string {
	if s, ok := accountTypeReverseMap[t]; ok {
		return s
	}
	return "other"
}

// Asset type conversions
var assetTypeMap = map[string]accountspb.AssetType{
	"bank":        accountspb.AssetType_ASSET_TYPE_BANK,
	"cash":        accountspb.AssetType_ASSET_TYPE_CASH,
	"stocks":      accountspb.AssetType_ASSET_TYPE_STOCKS,
	"crypto":      accountspb.AssetType_ASSET_TYPE_CRYPTO,
	"etf":         accountspb.AssetType_ASSET_TYPE_ETF,
	"real_estate": accountspb.AssetType_ASSET_TYPE_REAL_ESTATE,
	"bonds":       accountspb.AssetType_ASSET_TYPE_BONDS,
	"other":       accountspb.AssetType_ASSET_TYPE_OTHER,
}

func StringToAssetType(s string) accountspb.AssetType {
	if t, ok := assetTypeMap[strings.ToLower(s)]; ok {
		return t
	}
	return accountspb.AssetType_ASSET_TYPE_OTHER
}

// Transaction type conversions
var transactionTypeMap = map[string]transactionspb.TransactionType{
	"income":   transactionspb.TransactionType_TRANSACTION_TYPE_INCOME,
	"expense":  transactionspb.TransactionType_TRANSACTION_TYPE_EXPENSE,
	"transfer": transactionspb.TransactionType_TRANSACTION_TYPE_TRANSFER,
}

func StringToTransactionType(s string) transactionspb.TransactionType {
	if t, ok := transactionTypeMap[strings.ToLower(s)]; ok {
		return t
	}
	return transactionspb.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
}

// Transaction category conversions
var transactionCategoryMap = map[string]transactionspb.TransactionCategory{
	"salary":        transactionspb.TransactionCategory_TRANSACTION_CATEGORY_SALARY,
	"food":          transactionspb.TransactionCategory_TRANSACTION_CATEGORY_FOOD,
	"transport":     transactionspb.TransactionCategory_TRANSACTION_CATEGORY_TRANSPORT,
	"shopping":      transactionspb.TransactionCategory_TRANSACTION_CATEGORY_SHOPPING,
	"entertainment": transactionspb.TransactionCategory_TRANSACTION_CATEGORY_ENTERTAINMENT,
	"utilities":     transactionspb.TransactionCategory_TRANSACTION_CATEGORY_UTILITIES,
	"healthcare":    transactionspb.TransactionCategory_TRANSACTION_CATEGORY_HEALTHCARE,
	"transfer":      transactionspb.TransactionCategory_TRANSACTION_CATEGORY_TRANSFER,
	"other":         transactionspb.TransactionCategory_TRANSACTION_CATEGORY_OTHER,
}

func StringToTransactionCategory(s string) transactionspb.TransactionCategory {
	if t, ok := transactionCategoryMap[strings.ToLower(s)]; ok {
		return t
	}
	return transactionspb.TransactionCategory_TRANSACTION_CATEGORY_OTHER
}

// Date parsing helper
func ParseDate(s string) *timestamppb.Timestamp {
	if s == "" {
		return nil
	}

	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return timestamppb.New(t)
	}

	// Try YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return timestamppb.New(t)
	}

	return nil
}

// FormatTime formats timestamp to ISO string
func FormatTime(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format("2006-01-02T15:04:05Z")
}

// DefaultCurrency returns default currency if empty
func DefaultCurrency(currency string) string {
	if currency == "" {
		return "USD"
	}
	return currency
}

