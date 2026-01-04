import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { authAPI } from '../../services/api';
import { extractErrorMessage } from '../../utils/api-helpers';

interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  baseCurrency: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

// Try to restore user from localStorage
const getSavedUser = (): User | null => {
  try {
    const saved = localStorage.getItem('user');
    return saved ? JSON.parse(saved) : null;
  } catch {
    return null;
  }
};

const saveUser = (user: User | null) => {
  if (user) {
    localStorage.setItem('user', JSON.stringify(user));
  } else {
    localStorage.removeItem('user');
  }
};

const initialState: AuthState = {
  user: getSavedUser(),
  token: localStorage.getItem('token'),
  refreshToken: localStorage.getItem('refreshToken'),
  isAuthenticated: !!localStorage.getItem('token'),
  loading: !!localStorage.getItem('token') && !getSavedUser(), // Loading if token exists but no user
  error: null,
};

export const login = createAsyncThunk(
  'auth/login',
  async (credentials: { email: string; password: string }, { rejectWithValue }) => {
    try {
      const response = await authAPI.login(credentials);
      // API returns { success: true, data: { access_token, ... } }
      const data = response.data.data || response.data;
      localStorage.setItem('token', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      return {
        accessToken: data.access_token,
        refreshToken: data.refresh_token,
        user: {
          id: data.user.id,
          email: data.user.email,
          firstName: data.user.first_name || '',
          lastName: data.user.last_name || '',
          baseCurrency: data.user.base_currency || 'USD',
        },
      };
    } catch (error: any) {
      const message = error.response?.data?.error?.message || error.response?.data?.error || 'Login failed';
      return rejectWithValue(message);
    }
  }
);

export const register = createAsyncThunk(
  'auth/register',
  async (userData: { email: string; password: string; firstName: string; lastName: string }, { rejectWithValue }) => {
    try {
      const response = await authAPI.register(userData);
      // API returns { success: true, data: { access_token, ... } }
      const data = response.data.data || response.data;
      localStorage.setItem('token', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      return {
        accessToken: data.access_token,
        refreshToken: data.refresh_token,
        user: {
          id: data.user.id,
          email: data.user.email,
          firstName: data.user.first_name || '',
          lastName: data.user.last_name || '',
          baseCurrency: data.user.base_currency || 'USD',
        },
      };
    } catch (error: any) {
      const message = error.response?.data?.error?.message || error.response?.data?.error || 'Registration failed';
      return rejectWithValue(message);
    }
  }
);

export const getProfile = createAsyncThunk(
  'auth/getProfile',
  async (_, { rejectWithValue }) => {
    try {
      const response = await authAPI.getProfile();
      // API returns { success: true, data: { ... } }
      const data = response.data.data || response.data;
      return {
        id: data.id,
        email: data.email,
        firstName: data.first_name || '',
        lastName: data.last_name || '',
        baseCurrency: data.base_currency || 'USD',
      };
    } catch (error: any) {
      return rejectWithValue(extractErrorMessage(error, 'Failed to fetch profile'));
    }
  }
);

// Update user profile
export const updateProfile = createAsyncThunk(
  'auth/updateProfile',
  async (profileData: { firstName?: string; lastName?: string; baseCurrency?: string }, { rejectWithValue }) => {
    try {
      const response = await authAPI.updateProfile(profileData);
      const data = response.data.data || response.data;
      return {
        id: data.id,
        email: data.email,
        firstName: data.first_name || '',
        lastName: data.last_name || '',
        baseCurrency: data.base_currency || 'USD',
      };
    } catch (error: any) {
      return rejectWithValue(extractErrorMessage(error, 'Failed to update profile'));
    }
  }
);

// Telegram login - authenticate using Telegram initData
export const telegramLogin = createAsyncThunk(
  'auth/telegramLogin',
  async (initData: string, { rejectWithValue }) => {
    try {
      const response = await authAPI.telegramAuth(initData);
      const data = response.data.data || response.data;
      localStorage.setItem('token', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      return {
        accessToken: data.access_token,
        refreshToken: data.refresh_token,
        user: {
          id: data.user.id,
          email: data.user.email || '',
          firstName: data.user.first_name || '',
          lastName: data.user.last_name || '',
          baseCurrency: data.user.base_currency || 'USD',
        },
      };
    } catch (error: any) {
      const message = error.response?.data?.error?.message || error.response?.data?.error || 'Telegram login failed';
      return rejectWithValue(message);
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.refreshToken = null;
      state.isAuthenticated = false;
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('user');
    },
    clearError: (state) => {
      state.error = null;
    },
    setCredentials: (state, action: PayloadAction<{ accessToken: string; refreshToken: string }>) => {
      state.token = action.payload.accessToken;
      state.refreshToken = action.payload.refreshToken;
      state.isAuthenticated = true;
      localStorage.setItem('token', action.payload.accessToken);
      localStorage.setItem('refreshToken', action.payload.refreshToken);
    },
  },
  extraReducers: (builder) => {
    builder
      // Login
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.token = action.payload.accessToken;
        state.refreshToken = action.payload.refreshToken;
        state.user = action.payload.user;
        saveUser(action.payload.user);
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Register
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.token = action.payload.accessToken;
        state.refreshToken = action.payload.refreshToken;
        state.user = action.payload.user;
        saveUser(action.payload.user);
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Get Profile
      .addCase(getProfile.pending, (state) => {
        state.loading = true;
      })
      .addCase(getProfile.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        saveUser(action.payload);
      })
      .addCase(getProfile.rejected, (state, action) => {
        state.loading = false;
        const error = action.payload as string;
        // Only logout if it's an auth error (401), not other errors
        if (error?.includes('Unauthorized') || error?.includes('401') || error?.includes('token')) {
          state.isAuthenticated = false;
          state.token = null;
          state.user = null;
          localStorage.removeItem('token');
          localStorage.removeItem('refreshToken');
          localStorage.removeItem('user');
        }
        state.error = error || 'Failed to fetch profile';
      })
      // Update Profile
      .addCase(updateProfile.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateProfile.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
        saveUser(action.payload);
      })
      .addCase(updateProfile.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // Telegram Login
      .addCase(telegramLogin.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(telegramLogin.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.token = action.payload.accessToken;
        state.refreshToken = action.payload.refreshToken;
        state.user = action.payload.user;
        saveUser(action.payload.user);
      })
      .addCase(telegramLogin.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const { logout, clearError, setCredentials } = authSlice.actions;
export default authSlice.reducer;
