import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { insightsAPI } from '../../services/api';

export interface NetWorth {
  total: number;
  change24h: number;
  changePercent24h: number;
  change7d: number;
  changePercent7d: number;
  change30d: number;
  changePercent30d: number;
}

export interface AllocationItem {
  category: string;
  amount: number;
  percentage: number;
  color: string;
}

export interface TrendPoint {
  date: string;
  value: number;
}

interface InsightsState {
  netWorth: NetWorth | null;
  allocation: AllocationItem[];
  trends: TrendPoint[];
  period: '7d' | '30d' | '90d' | '1y' | 'all';
  loading: boolean;
  error: string | null;
}

const initialState: InsightsState = {
  netWorth: null,
  allocation: [],
  trends: [],
  period: '30d',
  loading: false,
  error: null,
};

// Helper to extract data from API response
const extractData = (response: any) => {
  return response.data?.data || response.data;
};

export const fetchNetWorth = createAsyncThunk(
  'insights/fetchNetWorth',
  async (baseCurrency: string, { rejectWithValue }) => {
    try {
      const response = await insightsAPI.getNetWorth(baseCurrency);
      const data = extractData(response);
      return {
        total: data?.total ?? data?.total_net_worth ?? data?.totalNetWorth ?? 0,
        change24h: data?.change24h ?? data?.change_24h ?? 0,
        changePercent24h: data?.changePercent24h ?? data?.change_percent_24h ?? 0,
        change7d: data?.change7d ?? data?.change_7d ?? 0,
        changePercent7d: data?.changePercent7d ?? data?.change_percent_7d ?? 0,
        change30d: data?.change30d ?? data?.change_30d ?? 0,
        changePercent30d: data?.changePercent30d ?? data?.change_percent_30d ?? 0,
      };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch net worth');
    }
  }
);

export const fetchAllocation = createAsyncThunk(
  'insights/fetchAllocation',
  async (baseCurrency: string, { rejectWithValue }) => {
    try {
      const response = await insightsAPI.getAllocation(baseCurrency);
      const data = extractData(response);
      const allocations = data?.allocations || data || [];
      return Array.isArray(allocations) ? allocations : [];
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch allocation');
    }
  }
);

export const fetchTrends = createAsyncThunk(
  'insights/fetchTrends',
  async ({ baseCurrency, period }: { baseCurrency: string; period: string }, { rejectWithValue }) => {
    try {
      const response = await insightsAPI.getTrends(baseCurrency, period);
      const data = extractData(response);
      const trends = data?.trends || data || [];
      return Array.isArray(trends) ? trends : [];
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch trends');
    }
  }
);

const insightsSlice = createSlice({
  name: 'insights',
  initialState,
  reducers: {
    setPeriod: (state, action) => {
      state.period = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      // Net Worth
      .addCase(fetchNetWorth.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchNetWorth.fulfilled, (state, action) => {
        state.loading = false;
        state.netWorth = action.payload;
      })
      .addCase(fetchNetWorth.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
        state.netWorth = null;
      })
      // Allocation
      .addCase(fetchAllocation.fulfilled, (state, action) => {
        state.allocation = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchAllocation.rejected, (state) => {
        state.allocation = [];
      })
      // Trends
      .addCase(fetchTrends.fulfilled, (state, action) => {
        state.trends = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchTrends.rejected, (state) => {
        state.trends = [];
      });
  },
});

export const { setPeriod } = insightsSlice.actions;
export default insightsSlice.reducer;
