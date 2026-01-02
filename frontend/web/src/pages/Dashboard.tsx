import React, { useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  TrendingUp,
  TrendingDown,
  Wallet,
  ArrowUpRight,
  ArrowDownRight,
  RefreshCw,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchNetWorth, fetchAllocation, fetchTrends } from '../store/slices/insightsSlice';
import { fetchAccounts } from '../store/slices/accountsSlice';
import AllocationChart from '../components/charts/AllocationChart';
import TrendChart from '../components/charts/TrendChart';

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  const { netWorth, allocation, trends, loading } = useAppSelector((state) => state.insights);
  const { accounts } = useAppSelector((state) => state.accounts);
  const { user } = useAppSelector((state) => state.auth);

  const baseCurrency = user?.baseCurrency || 'USD';

  useEffect(() => {
    dispatch(fetchNetWorth(baseCurrency));
    dispatch(fetchAllocation(baseCurrency));
    dispatch(fetchTrends({ baseCurrency, period: '30d' }));
    dispatch(fetchAccounts());
  }, [dispatch, baseCurrency]);

  const formatCurrency = (value: number | undefined | null) => {
    const safeValue = value ?? 0;
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: baseCurrency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(safeValue);
  };

  const formatPercent = (value: number | undefined | null) => {
    if (value === undefined || value === null) return '+0.00%';
    const sign = value >= 0 ? '+' : '';
    return `${sign}${value.toFixed(2)}%`;
  };

  const formatCurrencySafe = (value: number | undefined | null) => {
    return formatCurrency(value || 0);
  };

  const formatAccountType = (type: string | number): string => {
    const typeMap: { [key: number]: string } = {
      0: 'Unspecified',
      1: 'Bank',
      2: 'Cash',
      3: 'Investment',
      4: 'Crypto',
      5: 'Real Estate',
      6: 'Other',
    };
    if (typeof type === 'number') {
      return typeMap[type] || 'Other';
    }
    return type.charAt(0).toUpperCase() + type.slice(1).replace('_', ' ');
  };

  const handleRefresh = () => {
    dispatch(fetchNetWorth(baseCurrency));
    dispatch(fetchAllocation(baseCurrency));
    dispatch(fetchTrends({ baseCurrency, period: '30d' }));
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-display font-bold text-gradient">Dashboard</h1>
          <p className="text-midnight-400 mt-1">Your financial overview</p>
        </div>
        <button
          onClick={handleRefresh}
          disabled={loading}
          className="btn-primary flex items-center gap-2"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* Net Worth Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="glass rounded-2xl p-8 card-glow"
      >
        <div className="flex items-start justify-between">
          <div>
            <p className="text-midnight-400 text-sm font-medium uppercase tracking-wider">
              Total Net Worth
            </p>
            <h2 className="text-5xl font-display font-bold mt-2">
              {netWorth?.total !== undefined ? formatCurrency(netWorth.total) : '$0'}
            </h2>
            <div className="flex items-center gap-4 mt-4">
              {netWorth && (
                <>
                  <div className={`flex items-center gap-1 ${netWorth.change24h >= 0 ? 'text-accent-emerald' : 'text-accent-coral'}`}>
                    {netWorth.change24h >= 0 ? (
                      <TrendingUp className="w-4 h-4" />
                    ) : (
                      <TrendingDown className="w-4 h-4" />
                    )}
                    <span className="font-medium">
                      {formatPercent(netWorth.changePercent24h)}
                    </span>
                    <span className="text-midnight-400 text-sm">24h</span>
                  </div>
                  <div className={`flex items-center gap-1 ${netWorth.change7d >= 0 ? 'text-accent-emerald' : 'text-accent-coral'}`}>
                    {netWorth.change7d >= 0 ? (
                      <ArrowUpRight className="w-4 h-4" />
                    ) : (
                      <ArrowDownRight className="w-4 h-4" />
                    )}
                    <span className="font-medium">
                      {formatPercent(netWorth.changePercent7d)}
                    </span>
                    <span className="text-midnight-400 text-sm">7d</span>
                  </div>
                </>
              )}
            </div>
          </div>
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-midnight-500 to-accent-violet flex items-center justify-center">
            <Wallet className="w-8 h-8 text-white" />
          </div>
        </div>
      </motion.div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Performance Chart */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="glass rounded-2xl p-6"
        >
          <h3 className="text-lg font-semibold mb-4">Net Worth Trend</h3>
          <div className="h-64">
            <TrendChart data={trends} />
          </div>
        </motion.div>

        {/* Allocation Chart */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="glass rounded-2xl p-6"
        >
          <h3 className="text-lg font-semibold mb-4">Asset Allocation</h3>
          <div className="h-64">
            <AllocationChart data={allocation} />
          </div>
        </motion.div>
      </div>

      {/* Accounts Overview */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="glass rounded-2xl p-6"
      >
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold">Accounts</h3>
          <a
            href="/accounts"
            className="text-sm text-midnight-400 hover:text-white transition-colors"
          >
            View all →
          </a>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {(Array.isArray(accounts) ? accounts : []).slice(0, 6).map((account, index) => (
            <motion.div
              key={account.id}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.1 * index }}
              className="p-4 rounded-xl bg-midnight-950/50 border border-midnight-800/50 hover:border-midnight-500/30 transition-all duration-200"
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-midnight-600 to-midnight-800 flex items-center justify-center">
                  <Wallet className="w-5 h-5 text-midnight-300" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{account.name}</p>
                  <p className="text-sm text-midnight-400">{formatAccountType(account.type)}</p>
                </div>
              </div>
              <div className="mt-3 flex items-baseline justify-between">
                <span className="text-lg font-semibold">
                  {formatCurrency(account.totalBalance)}
                </span>
                <span className="text-xs text-midnight-400">{account.currency}</span>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </div>
  );
};

export default Dashboard;

