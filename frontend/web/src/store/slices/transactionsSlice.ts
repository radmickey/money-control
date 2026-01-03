import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { transactionsAPI } from '../../services/api';

export interface Transaction {
  id: string;
  accountId: string;
  subAccountId: string;
  amount: number;
  currency: string;
  type: 'income' | 'expense' | 'transfer';
  category: string;
  description: string;
  date: string;
  createdAt: string;
}

interface TransactionsState {
  transactions: Transaction[];
  filters: {
    accountId?: string;
    category?: string;
    type?: string;
    startDate?: string;
    endDate?: string;
  };
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };
  loading: boolean;
  error: string | null;
}

const initialState: TransactionsState = {
  transactions: [],
  filters: {},
  pagination: {
    page: 1,
    pageSize: 20,
    total: 0,
  },
  loading: false,
  error: null,
};

// Helper to extract data from API response
const extractData = (response: any) => {
  return response.data?.data || response.data;
};

// Map API transaction to frontend Transaction
const mapTransaction = (data: any): Transaction => ({
  id: data.id,
  accountId: data.account_id || '',
  subAccountId: data.sub_account_id || '',
  amount: data.amount || 0,
  currency: data.currency || 'USD',
  type: data.type || 'expense',
  category: data.category || 'other',
  description: data.description || '',
  date: data.date || new Date().toISOString(),
  createdAt: data.created_at || new Date().toISOString(),
});

export const fetchTransactions = createAsyncThunk(
  'transactions/fetchTransactions',
  async (params: any, { rejectWithValue }) => {
    try {
      const response = await transactionsAPI.list(params);
      const data = extractData(response);
      const transactions = data?.transactions || data || [];
      return {
        transactions: Array.isArray(transactions) ? transactions.map(mapTransaction) : [],
        total: data?.total || 0,
      };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch transactions');
    }
  }
);

export const createTransaction = createAsyncThunk(
  'transactions/createTransaction',
  async (transactionData: Partial<Transaction>, { rejectWithValue }) => {
    try {
      const response = await transactionsAPI.create({
        account_id: transactionData.accountId,
        sub_account_id: transactionData.subAccountId,
        amount: transactionData.amount,
        currency: transactionData.currency,
        type: transactionData.type,
        category: transactionData.category,
        description: transactionData.description,
        date: transactionData.date,
      });
      const data = extractData(response);
      return mapTransaction(data);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to create transaction');
    }
  }
);

export const updateTransaction = createAsyncThunk(
  'transactions/updateTransaction',
  async ({ id, data }: { id: string; data: Partial<Transaction> }, { rejectWithValue }) => {
    try {
      const response = await transactionsAPI.update(id, {
        amount: data.amount,
        type: data.type,
        category: data.category,
        description: data.description,
        date: data.date,
        currency: data.currency,
      });
      const respData = extractData(response);
      return mapTransaction(respData);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to update transaction');
    }
  }
);

export const deleteTransaction = createAsyncThunk(
  'transactions/deleteTransaction',
  async (id: string, { rejectWithValue }) => {
    try {
      await transactionsAPI.delete(id);
      return id;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to delete transaction');
    }
  }
);

const transactionsSlice = createSlice({
  name: 'transactions',
  initialState,
  reducers: {
    setFilters: (state, action) => {
      state.filters = { ...state.filters, ...action.payload };
    },
    clearFilters: (state) => {
      state.filters = {};
    },
    setPage: (state, action) => {
      state.pagination.page = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTransactions.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTransactions.fulfilled, (state, action) => {
        state.loading = false;
        state.transactions = Array.isArray(action.payload?.transactions) ? action.payload.transactions : [];
        state.pagination.total = action.payload?.total || 0;
      })
      .addCase(fetchTransactions.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
        state.transactions = [];
      })
      .addCase(createTransaction.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createTransaction.fulfilled, (state, action) => {
        state.loading = false;
        if (!Array.isArray(state.transactions)) {
          state.transactions = [];
        }
        state.transactions.unshift(action.payload);
      })
      .addCase(createTransaction.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(updateTransaction.fulfilled, (state, action) => {
        const index = state.transactions.findIndex(t => t.id === action.payload.id);
        if (index !== -1) {
          state.transactions[index] = action.payload;
        }
      })
      .addCase(deleteTransaction.fulfilled, (state, action) => {
        state.transactions = state.transactions.filter(t => t.id !== action.payload);
      });
  },
});

export const { setFilters, clearFilters, setPage } = transactionsSlice.actions;
export default transactionsSlice.reducer;
