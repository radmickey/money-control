import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { assetsAPI } from '../../services/api';

export interface Asset {
  id: string;
  userId: string;
  accountId: string;
  type: 'stock' | 'crypto' | 'etf' | 'real_estate' | 'other';
  symbol: string;
  name: string;
  quantity: number;
  purchasePrice: number;
  currentPrice: number;
  currency: string;
  lastUpdated: string;
}

interface AssetsState {
  assets: Asset[];
  prices: { [symbol: string]: number };
  loading: boolean;
  pricesLoading: boolean;
  error: string | null;
}

const initialState: AssetsState = {
  assets: [],
  prices: {},
  loading: false,
  pricesLoading: false,
  error: null,
};

// Helper to extract data from API response
const extractData = (response: any) => {
  return response.data?.data || response.data;
};

// Map API asset to frontend Asset
const mapAsset = (data: any): Asset => ({
  id: data.id,
  userId: data.user_id || '',
  accountId: data.account_id || '',
  type: data.type || 'other',
  symbol: data.symbol || '',
  name: data.name || '',
  quantity: data.quantity || 0,
  purchasePrice: data.purchase_price || 0,
  currentPrice: data.current_price || 0,
  currency: data.currency || 'USD',
  lastUpdated: data.last_updated || new Date().toISOString(),
});

export const fetchAssets = createAsyncThunk(
  'assets/fetchAssets',
  async (_, { rejectWithValue }) => {
    try {
      const response = await assetsAPI.list();
      const data = extractData(response);
      const assets = data?.assets || data || [];
      return Array.isArray(assets) ? assets.map(mapAsset) : [];
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch assets');
    }
  }
);

export const createAsset = createAsyncThunk(
  'assets/createAsset',
  async (assetData: Partial<Asset>, { rejectWithValue }) => {
    try {
      const response = await assetsAPI.create({
        type: assetData.type,
        symbol: assetData.symbol,
        name: assetData.name,
        quantity: assetData.quantity,
        purchase_price: assetData.purchasePrice,
        currency: assetData.currency,
        account_id: assetData.accountId,
      });
      const data = extractData(response);
      return mapAsset(data);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to create asset');
    }
  }
);

export const updateAsset = createAsyncThunk(
  'assets/updateAsset',
  async ({ id, data }: { id: string; data: Partial<Asset> }, { rejectWithValue }) => {
    try {
      const response = await assetsAPI.update(id, {
        symbol: data.symbol,
        name: data.name,
        quantity: data.quantity,
        purchase_price: data.purchasePrice,
        current_price: data.currentPrice,
      });
      const respData = extractData(response);
      return mapAsset(respData);
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to update asset');
    }
  }
);

export const deleteAsset = createAsyncThunk(
  'assets/deleteAsset',
  async (id: string, { rejectWithValue }) => {
    try {
      await assetsAPI.delete(id);
      return id;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to delete asset');
    }
  }
);

export const fetchAssetPrice = createAsyncThunk(
  'assets/fetchAssetPrice',
  async ({ symbol, type }: { symbol: string; type: string }, { rejectWithValue }) => {
    try {
      const response = await assetsAPI.getPrice(symbol, type);
      const data = extractData(response);
      return { symbol, price: data?.price || 0 };
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error?.message || error.response?.data?.error || 'Failed to fetch price');
    }
  }
);

const assetsSlice = createSlice({
  name: 'assets',
  initialState,
  reducers: {
    updatePrice: (state, action) => {
      state.prices[action.payload.symbol] = action.payload.price;
    },
    clearAssetsError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchAssets.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchAssets.fulfilled, (state, action) => {
        state.loading = false;
        state.assets = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchAssets.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
        state.assets = [];
      })
      .addCase(createAsset.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createAsset.fulfilled, (state, action) => {
        state.loading = false;
        if (!Array.isArray(state.assets)) {
          state.assets = [];
        }
        state.assets.push(action.payload);
      })
      .addCase(createAsset.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(updateAsset.fulfilled, (state, action) => {
        const index = state.assets.findIndex(a => a.id === action.payload.id);
        if (index !== -1) {
          state.assets[index] = action.payload;
        }
      })
      .addCase(deleteAsset.fulfilled, (state, action) => {
        state.assets = state.assets.filter(a => a.id !== action.payload);
      })
      .addCase(fetchAssetPrice.pending, (state) => {
        state.pricesLoading = true;
      })
      .addCase(fetchAssetPrice.fulfilled, (state, action) => {
        state.pricesLoading = false;
        state.prices[action.payload.symbol] = action.payload.price;
      })
      .addCase(fetchAssetPrice.rejected, (state) => {
        state.pricesLoading = false;
      });
  },
});

export const { updatePrice, clearAssetsError } = assetsSlice.actions;
export default assetsSlice.reducer;
