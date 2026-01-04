// Currency formatting
export const formatCurrency = (
  value: number | undefined | null,
  currency: string = 'USD',
  showDecimals: boolean = true
): string => {
  const safeValue = value ?? 0;
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
      minimumFractionDigits: showDecimals ? 2 : 0,
      maximumFractionDigits: showDecimals ? 2 : 0,
    }).format(safeValue);
  } catch {
    return `${currency} ${safeValue.toFixed(showDecimals ? 2 : 0)}`;
  }
};

// Percent formatting
export const formatPercent = (value: number | undefined | null): string => {
  if (value === undefined || value === null) return '+0.00%';
  const sign = value >= 0 ? '+' : '';
  return `${sign}${value.toFixed(2)}%`;
};

// Date formatting
export const formatDate = (
  dateValue: string | { seconds?: number; nanos?: number } | Date
): string => {
  let date: Date;
  if (dateValue instanceof Date) {
    date = dateValue;
  } else if (typeof dateValue === 'object' && dateValue.seconds) {
    date = new Date(dateValue.seconds * 1000);
  } else if (typeof dateValue === 'string') {
    date = new Date(dateValue);
  } else {
    return 'N/A';
  }
  if (isNaN(date.getTime())) return 'N/A';
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
};

// Account type enum to string
const ACCOUNT_TYPE_MAP: Record<number, string> = {
  0: 'Other',
  1: 'Bank',
  2: 'Cash',
  3: 'Investment',
  4: 'Crypto',
  5: 'Real Estate',
  6: 'Other',
};

export const formatAccountType = (type: string | number): string => {
  if (typeof type === 'number') {
    return ACCOUNT_TYPE_MAP[type] || 'Other';
  }
  return type.charAt(0).toUpperCase() + type.slice(1).replace('_', ' ');
};

export const getAccountTypeKey = (type: string | number): string => {
  if (typeof type === 'number') {
    return ['other', 'bank', 'cash', 'investment', 'crypto', 'real_estate', 'other'][type] || 'other';
  }
  return type.toLowerCase();
};

// Transaction type enum to string
const TRANSACTION_TYPE_MAP: Record<number, string> = {
  0: 'unspecified',
  1: 'income',
  2: 'expense',
  3: 'transfer',
};

export const getTransactionType = (type: string | number): string => {
  if (typeof type === 'number') {
    return TRANSACTION_TYPE_MAP[type] || 'expense';
  }
  return type;
};

// Transaction category enum to string
const CATEGORY_MAP: Record<number, string> = {
  0: 'Other',
  1: 'Salary',
  2: 'Food & Dining',
  3: 'Transportation',
  4: 'Shopping',
  5: 'Entertainment',
  6: 'Healthcare',
  7: 'Bills & Utilities',
  8: 'Transfer',
  9: 'Other',
};

export const getTransactionCategory = (category: string | number): string => {
  if (typeof category === 'number') {
    return CATEGORY_MAP[category] || 'Other';
  }
  return category;
};

