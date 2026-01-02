import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../../services/api';

export interface Account {
  id: string;
  name: string;
  type: string;
  currency: string;
  totalBalance: number;
}

interface AccountsState {
  accounts: Account[];
  loading: boolean;
  error: string | null;
}

const initialState: AccountsState = {
  accounts: [],
  loading: false,
  error: null,
};

export const fetchAccounts = createAsyncThunk('accounts/fetchAccounts', async (_, { rejectWithValue }) => {
  try {
    const response = await api.get('/accounts');
    return response.data.accounts;
  } catch (error: any) {
    return rejectWithValue(error.response?.data?.error || 'Failed to fetch accounts');
  }
});

const accountsSlice = createSlice({
  name: 'accounts',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchAccounts.pending, (state) => { state.loading = true; })
      .addCase(fetchAccounts.fulfilled, (state, action) => {
        state.loading = false;
        state.accounts = action.payload;
      })
      .addCase(fetchAccounts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export default accountsSlice.reducer;

