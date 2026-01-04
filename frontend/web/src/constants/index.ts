import type { ElementType } from 'react';
import {
  Wallet,
  Building2,
  Banknote,
  TrendingUp,
  Bitcoin,
  Home,
  BarChart3,
} from 'lucide-react';

// Account type icons
export const ACCOUNT_TYPE_ICONS: Record<string, ElementType> = {
  bank: Building2,
  cash: Banknote,
  investment: TrendingUp,
  crypto: Bitcoin,
  real_estate: Home,
  other: Wallet,
};

// Dashboard uses BarChart3 for investment
export const DASHBOARD_ACCOUNT_ICONS: Record<string, ElementType> = {
  bank: Building2,
  cash: Banknote,
  investment: BarChart3,
  crypto: Bitcoin,
  real_estate: Home,
  other: Wallet,
};

// Currencies list
export const CURRENCIES = [
  { code: 'USD', name: 'US Dollar', symbol: '$' },
  { code: 'EUR', name: 'Euro', symbol: '€' },
  { code: 'GBP', name: 'British Pound', symbol: '£' },
  { code: 'RUB', name: 'Russian Ruble', symbol: '₽' },
  { code: 'JPY', name: 'Japanese Yen', symbol: '¥' },
  { code: 'CHF', name: 'Swiss Franc', symbol: 'CHF' },
  { code: 'CAD', name: 'Canadian Dollar', symbol: 'C$' },
  { code: 'AUD', name: 'Australian Dollar', symbol: 'A$' },
] as const;

// Account types for forms
export const ACCOUNT_TYPES = [
  { value: 'bank', label: 'Bank' },
  { value: 'cash', label: 'Cash' },
  { value: 'investment', label: 'Investment' },
  { value: 'crypto', label: 'Crypto' },
  { value: 'real_estate', label: 'Real Estate' },
  { value: 'other', label: 'Other' },
] as const;

// Transaction categories
export const TRANSACTION_CATEGORIES = [
  'Food & Dining',
  'Transportation',
  'Shopping',
  'Entertainment',
  'Bills & Utilities',
  'Healthcare',
  'Travel',
  'Education',
  'Income',
  'Investment',
  'Transfer',
  'Other',
] as const;

// Transaction types
export const TRANSACTION_TYPES = ['expense', 'income', 'transfer'] as const;

