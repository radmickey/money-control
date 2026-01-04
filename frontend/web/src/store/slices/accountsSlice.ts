import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { accountsAPI } from '../../services/api';

export interface SubAccount {
  id: string;
  name: string;
  currency: string;
  balance: number;
  convertedBalance?: number;
  assetType: string;
}

export interface Account {
  id: string;
  name: string;
  type: string;
  currency: string;
  totalBalance: number;
  // Pre-calculated by backend
  convertedTotalBalance?: number;  // In account's display currency
  balanceInBaseCurrency?: number;  // In user's base currency (for net worth)
  displayCurrency?: string;
  isMixedCurrency?: boolean;
  subAccounts: SubAccount[];
  icon?: string;
  color?: string;
}

interface AccountsState {
  accounts: Account[];
  selectedAccount: Account | null;
  loading: boolean;
  error: string | null;
}

const initialState: AccountsState = {
  accounts: [],
  selectedAccount: null,
  loading: false,
  error: null,
};

// Helper to extract data from API response
const extractData = (response: any) => {
  const data = response.data?.data || response.data;
  return data;
};

// Map API account to frontend Account type
const mapApiAccount = (apiAccount: any): Account => ({
  id: apiAccount.id,
  name: apiAccount.name,
  type: apiAccount.type,
  currency: apiAccount.currency,
  totalBalance: apiAccount.total_balance || apiAccount.totalBalance || 0,
  // New fields from backend conversion
  convertedTotalBalance: apiAccount.converted_total_balance ?? apiAccount.convertedTotalBalance ?? apiAccount.total_balance ?? 0,
  balanceInBaseCurrency: apiAccount.balance_in_base_currency ?? apiAccount.balanceInBaseCurrency ?? apiAccount.total_balance ?? 0,
  displayCurrency: apiAccount.display_currency || apiAccount.displayCurrency || apiAccount.currency || 'USD',
  isMixedCurrency: apiAccount.is_mixed_currency ?? apiAccount.isMixedCurrency ?? false,
  subAccounts: (apiAccount.sub_accounts || apiAccount.subAccounts || []).map((sub: any) => ({
    id: sub.id,
    name: sub.name,
    currency: sub.currency,
    balance: sub.balance || 0,
    convertedBalance: sub.converted_balance ?? sub.convertedBalance ?? sub.balance ?? 0,
    assetType: sub.asset_type || sub.assetType || 'cash',
  })),
  icon: apiAccount.icon,
  color: apiAccount.color,
});

export const fetchAccounts = createAsyncThunk(
  'accounts/fetchAccounts',
  async (_, { rejectWithValue }) => {
    try {
      const response = await accountsAPI.list();
      const data = extractData(response);
      const accounts = data?.accounts || data || [];
      return Array.isArray(accounts) ? accounts.map(mapApiAccount) : [];
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch accounts');
    }
  }
);

export const createAccount = createAsyncThunk(
  'accounts/createAccount',
  async (accountData: { name: string; type: string; currency: string }, { rejectWithValue }) => {
    try {
      const response = await accountsAPI.create(accountData);
      const data = extractData(response);
      return mapApiAccount(data);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to create account');
    }
  }
);

export const updateAccount = createAsyncThunk(
  'accounts/updateAccount',
  async ({ id, data }: { id: string; data: Partial<Account> }, { rejectWithValue }) => {
    try {
      const response = await accountsAPI.update(id, data);
      const respData = extractData(response);
      return mapApiAccount(respData);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to update account');
    }
  }
);

export const deleteAccount = createAsyncThunk(
  'accounts/deleteAccount',
  async (id: string, { rejectWithValue }) => {
    try {
      await accountsAPI.delete(id);
      return id;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to delete account');
    }
  }
);

// Sub-account actions
export const createSubAccount = createAsyncThunk(
  'accounts/createSubAccount',
  async ({ accountId, data }: { accountId: string; data: { name: string; balance: number; currency: string } }, { rejectWithValue }) => {
    try {
      const response = await accountsAPI.createSubAccount(accountId, {
        account_id: accountId,
        ...data,
      });
      const respData = extractData(response);
      return {
        accountId,
        subAccount: {
          id: respData.id,
          name: respData.name,
          currency: respData.currency,
          balance: respData.balance || 0,
          assetType: respData.asset_type || 'cash',
        },
      };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to create sub-account');
    }
  }
);

export const updateSubAccount = createAsyncThunk(
  'accounts/updateSubAccount',
  async ({ accountId, subAccountId, data }: { accountId: string; subAccountId: string; data: { name?: string; balance?: number } }, { rejectWithValue }) => {
    try {
      const response = await accountsAPI.updateSubAccount(accountId, subAccountId, data);
      const respData = extractData(response);
      return {
        accountId,
        subAccount: {
          id: respData.id,
          name: respData.name,
          currency: respData.currency,
          balance: respData.balance || 0,
          assetType: respData.asset_type || 'cash',
        },
      };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to update sub-account');
    }
  }
);

export const deleteSubAccount = createAsyncThunk(
  'accounts/deleteSubAccount',
  async ({ accountId, subAccountId }: { accountId: string; subAccountId: string }, { rejectWithValue }) => {
    try {
      await accountsAPI.deleteSubAccount(accountId, subAccountId);
      return { accountId, subAccountId };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to delete sub-account');
    }
  }
);

const accountsSlice = createSlice({
  name: 'accounts',
  initialState,
  reducers: {
    selectAccount: (state, action) => {
      state.selectedAccount = action.payload;
    },
    clearAccountsError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch
      .addCase(fetchAccounts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAccounts.fulfilled, (state, action) => {
        state.loading = false;
        state.accounts = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchAccounts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
        state.accounts = [];
      })
      // Create
      .addCase(createAccount.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createAccount.fulfilled, (state, action) => {
        state.loading = false;
        if (!Array.isArray(state.accounts)) {
          state.accounts = [];
        }
        state.accounts.push(action.payload);
      })
      .addCase(createAccount.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Update
      .addCase(updateAccount.fulfilled, (state, action) => {
        const index = state.accounts.findIndex(a => a.id === action.payload.id);
        if (index !== -1) {
          state.accounts[index] = action.payload;
        }
      })
      // Delete
      .addCase(deleteAccount.fulfilled, (state, action) => {
        state.accounts = state.accounts.filter(a => a.id !== action.payload);
      })
      // Create Sub-Account
      .addCase(createSubAccount.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createSubAccount.fulfilled, (state, action) => {
        state.loading = false;
        const account = state.accounts.find(a => a.id === action.payload.accountId);
        if (account) {
          if (!account.subAccounts) {
            account.subAccounts = [];
          }
          account.subAccounts.push(action.payload.subAccount);
          account.totalBalance = (account.totalBalance || 0) + (action.payload.subAccount.balance || 0);
        }
      })
      .addCase(createSubAccount.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Update Sub-Account
      .addCase(updateSubAccount.fulfilled, (state, action) => {
        const account = state.accounts.find(a => a.id === action.payload.accountId);
        if (account && account.subAccounts) {
          const subIndex = account.subAccounts.findIndex(s => s.id === action.payload.subAccount.id);
          if (subIndex !== -1) {
            const oldBalance = account.subAccounts[subIndex].balance || 0;
            const newBalance = action.payload.subAccount.balance || 0;
            account.subAccounts[subIndex] = action.payload.subAccount;
            account.totalBalance = (account.totalBalance || 0) - oldBalance + newBalance;
          }
        }
      })
      // Delete Sub-Account
      .addCase(deleteSubAccount.fulfilled, (state, action) => {
        const account = state.accounts.find(a => a.id === action.payload.accountId);
        if (account && account.subAccounts) {
          const subAccount = account.subAccounts.find(s => s.id === action.payload.subAccountId);
          if (subAccount) {
            account.totalBalance = (account.totalBalance || 0) - (subAccount.balance || 0);
          }
          account.subAccounts = account.subAccounts.filter(s => s.id !== action.payload.subAccountId);
        }
      });
  },
});

export const { selectAccount, clearAccountsError } = accountsSlice.actions;
export default accountsSlice.reducer;
