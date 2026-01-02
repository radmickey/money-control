import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../../services/api';

interface NetWorth {
  total: number;
  change24h: number;
  changePercent24h: number;
}

interface AllocationItem {
  category: string;
  amount: number;
  percentage: number;
  color: string;
}

interface InsightsState {
  netWorth: NetWorth | null;
  allocation: AllocationItem[];
  loading: boolean;
  error: string | null;
}

const initialState: InsightsState = {
  netWorth: null,
  allocation: [],
  loading: false,
  error: null,
};

export const fetchNetWorth = createAsyncThunk('insights/fetchNetWorth', async (baseCurrency: string, { rejectWithValue }) => {
  try {
    const response = await api.get('/insights/net-worth', { params: { baseCurrency } });
    return response.data;
  } catch (error: any) {
    return rejectWithValue(error.response?.data?.error || 'Failed to fetch net worth');
  }
});

export const fetchAllocation = createAsyncThunk('insights/fetchAllocation', async (baseCurrency: string, { rejectWithValue }) => {
  try {
    const response = await api.get('/insights/allocation', { params: { baseCurrency } });
    return response.data.allocations;
  } catch (error: any) {
    return rejectWithValue(error.response?.data?.error || 'Failed to fetch allocation');
  }
});

const insightsSlice = createSlice({
  name: 'insights',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchNetWorth.pending, (state) => { state.loading = true; })
      .addCase(fetchNetWorth.fulfilled, (state, action) => {
        state.loading = false;
        state.netWorth = {
          total: action.payload.totalNetWorth,
          change24h: action.payload.change24h,
          changePercent24h: action.payload.changePercent24h,
        };
      })
      .addCase(fetchNetWorth.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(fetchAllocation.fulfilled, (state, action) => {
        state.allocation = action.payload;
      });
  },
});

export default insightsSlice.reducer;

