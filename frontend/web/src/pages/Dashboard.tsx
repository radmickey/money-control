import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  TrendingUp,
  TrendingDown,
  Wallet,
  Building2,
  Banknote,
  Bitcoin,
  Home,
  BarChart3,
  ArrowUpRight,
  ArrowDownRight,
  ArrowLeftRight,
  RefreshCw,
  PieChart,
  DollarSign,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchNetWorth, fetchAllocation, fetchTrends } from '../store/slices/insightsSlice';
import { fetchAccounts } from '../store/slices/accountsSlice';
import { fetchAssets } from '../store/slices/assetsSlice';
import { fetchTransactions } from '../store/slices/transactionsSlice';
import AllocationChart from '../components/charts/AllocationChart';
import TrendChart from '../components/charts/TrendChart';

const accountTypeIcons: { [key: string]: React.ElementType } = {
  bank: Building2,
  cash: Banknote,
  investment: BarChart3,
  crypto: Bitcoin,
  real_estate: Home,
  other: Wallet,
};

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  const { netWorth, allocation, trends, loading } = useAppSelector((state) => state.insights);
  const { accounts } = useAppSelector((state) => state.accounts);
  const { assets } = useAppSelector((state) => state.assets);
  const { transactions } = useAppSelector((state) => state.transactions);
  const { user } = useAppSelector((state) => state.auth);

  const baseCurrency = user?.baseCurrency || 'USD';

  useEffect(() => {
    dispatch(fetchNetWorth(baseCurrency));
    dispatch(fetchAllocation(baseCurrency));
    dispatch(fetchTrends({ baseCurrency, period: '30d' }));
    dispatch(fetchAccounts());
    dispatch(fetchAssets());
    dispatch(fetchTransactions({}));
  }, [dispatch, baseCurrency]);

  const formatCurrency = (value: number | undefined | null, currency?: string) => {
    const safeValue = value ?? 0;
    const curr = currency || baseCurrency;
    try {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: curr,
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
      }).format(safeValue);
    } catch {
      return `${curr} ${safeValue.toFixed(0)}`;
    }
  };

  const formatPercent = (value: number | undefined | null) => {
    if (value === undefined || value === null) return '+0.00%';
    const sign = value >= 0 ? '+' : '';
    return `${sign}${value.toFixed(2)}%`;
  };

  const getAccountTypeIcon = (type: string | number) => {
    const typeMap: { [key: number]: string } = {
      0: 'other', 1: 'bank', 2: 'cash', 3: 'investment', 4: 'crypto', 5: 'real_estate', 6: 'other',
    };
    const typeString = typeof type === 'number' ? (typeMap[type] || 'other') : type;
    return accountTypeIcons[typeString.toLowerCase()] || Wallet;
  };

  const formatAccountType = (type: string | number): string => {
    const typeMap: { [key: number]: string } = {
      0: 'Other', 1: 'Bank', 2: 'Cash', 3: 'Investment', 4: 'Crypto', 5: 'Real Estate', 6: 'Other',
    };
    if (typeof type === 'number') return typeMap[type] || 'Other';
    return type.charAt(0).toUpperCase() + type.slice(1).replace('_', ' ');
  };

  const getTransactionType = (type: string | number): string => {
    const typeMap: { [key: number]: string } = { 0: 'expense', 1: 'income', 2: 'expense', 3: 'transfer' };
    return typeof type === 'number' ? (typeMap[type] || 'expense') : type;
  };

  const handleRefresh = () => {
    dispatch(fetchNetWorth(baseCurrency));
    dispatch(fetchAllocation(baseCurrency));
    dispatch(fetchTrends({ baseCurrency, period: '30d' }));
    dispatch(fetchAccounts());
    dispatch(fetchAssets());
    dispatch(fetchTransactions({}));
  };

  // Calculate totals by account type
  const accountsByType = (Array.isArray(accounts) ? accounts : []).reduce((acc, account) => {
    const type = formatAccountType(account.type);
    acc[type] = (acc[type] || 0) + (account.totalBalance || 0);
    return acc;
  }, {} as { [key: string]: number });

  // Calculate total assets value
  const totalAssetsValue = (Array.isArray(assets) ? assets : []).reduce(
    (sum, asset) => sum + (asset.currentValue || asset.purchasePrice || 0) * (asset.quantity || 1),
    0
  );

  // Recent transactions
  const recentTransactions = (Array.isArray(transactions) ? transactions : []).slice(0, 5);

  return (
    <div className="space-y-6">
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

      {/* Net Worth Hero */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="glass rounded-2xl p-8 card-glow relative overflow-hidden"
      >
        <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-accent-violet/20 to-transparent rounded-full blur-3xl" />
        <div className="relative z-10">
          <p className="text-midnight-400 text-sm font-medium uppercase tracking-wider">
            Total Net Worth
          </p>
          <h2 className="text-5xl md:text-6xl font-display font-bold mt-2 bg-gradient-to-r from-white to-zinc-400 bg-clip-text text-transparent">
            {formatCurrency(netWorth?.total)}
          </h2>
          <div className="flex items-center gap-6 mt-4">
            <div className={`flex items-center gap-2 ${(netWorth?.change24h || 0) >= 0 ? 'text-accent-emerald' : 'text-accent-coral'}`}>
              {(netWorth?.change24h || 0) >= 0 ? <TrendingUp className="w-5 h-5" /> : <TrendingDown className="w-5 h-5" />}
              <span className="text-lg font-semibold">{formatPercent(netWorth?.changePercent24h)}</span>
              <span className="text-midnight-400">today</span>
            </div>
            <div className={`flex items-center gap-2 ${(netWorth?.change7d || 0) >= 0 ? 'text-accent-emerald' : 'text-accent-coral'}`}>
              {(netWorth?.change7d || 0) >= 0 ? <ArrowUpRight className="w-5 h-5" /> : <ArrowDownRight className="w-5 h-5" />}
              <span className="text-lg font-semibold">{formatPercent(netWorth?.changePercent7d)}</span>
              <span className="text-midnight-400">this week</span>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Quick Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {Object.entries(accountsByType).slice(0, 4).map(([type, value], index) => {
          const Icon = accountTypeIcons[type.toLowerCase()] || Wallet;
          return (
            <motion.div
              key={type}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * index }}
              className="glass rounded-xl p-4 hover:border-midnight-500/30 transition-all"
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-midnight-600 to-midnight-800 flex items-center justify-center">
                  <Icon className="w-5 h-5 text-midnight-300" />
                </div>
                <div>
                  <p className="text-sm text-midnight-400">{type}</p>
                  <p className="text-lg font-semibold">{formatCurrency(value)}</p>
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Net Worth Trend */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="glass rounded-2xl p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-accent-emerald" />
              Net Worth Trend
            </h3>
            <select className="bg-midnight-900/50 border border-midnight-700 rounded-lg px-3 py-1 text-sm">
              <option>30 Days</option>
              <option>90 Days</option>
              <option>1 Year</option>
            </select>
          </div>
          <div className="h-64">
            <TrendChart data={trends} />
          </div>
        </motion.div>

        {/* Allocation */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="glass rounded-2xl p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold flex items-center gap-2">
              <PieChart className="w-5 h-5 text-accent-violet" />
              Asset Allocation
            </h3>
          </div>
          <div className="h-64">
            <AllocationChart data={allocation} />
          </div>
        </motion.div>
      </div>

      {/* Bottom Grid: Accounts, Assets, Transactions */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Accounts */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="glass rounded-2xl p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Accounts</h3>
            <Link to="/accounts" className="text-sm text-accent-violet hover:text-accent-violet/80">
              View all →
            </Link>
          </div>
          <div className="space-y-3">
            {(Array.isArray(accounts) ? accounts : []).slice(0, 4).map((account) => {
              const Icon = getAccountTypeIcon(account.type);
              return (
                <div key={account.id} className="flex items-center justify-between p-3 rounded-lg bg-midnight-900/30 hover:bg-midnight-800/30 transition-colors">
                  <div className="flex items-center gap-3">
                    <Icon className="w-5 h-5 text-midnight-400" />
                    <span className="font-medium truncate max-w-[120px]">{account.name}</span>
                  </div>
                  <span className="font-semibold">{formatCurrency(account.totalBalance, account.currency)}</span>
                </div>
              );
            })}
            {(!accounts || accounts.length === 0) && (
              <p className="text-midnight-400 text-center py-4">No accounts yet</p>
            )}
          </div>
        </motion.div>

        {/* Assets */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="glass rounded-2xl p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Assets</h3>
            <Link to="/assets" className="text-sm text-accent-violet hover:text-accent-violet/80">
              View all →
            </Link>
          </div>
          <div className="space-y-3">
            {(Array.isArray(assets) ? assets : []).slice(0, 4).map((asset) => (
              <div key={asset.id} className="flex items-center justify-between p-3 rounded-lg bg-midnight-900/30 hover:bg-midnight-800/30 transition-colors">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-gradient-to-br from-accent-amber to-accent-coral flex items-center justify-center text-xs font-bold">
                    {asset.symbol?.slice(0, 2) || '??'}
                  </div>
                  <div>
                    <p className="font-medium">{asset.symbol || asset.name}</p>
                    <p className="text-xs text-midnight-400">{asset.quantity} units</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-semibold">{formatCurrency((asset.currentValue || asset.purchasePrice || 0) * (asset.quantity || 1))}</p>
                  {asset.currentValue && asset.purchasePrice && (
                    <p className={`text-xs ${asset.currentValue > asset.purchasePrice ? 'text-accent-emerald' : 'text-accent-coral'}`}>
                      {asset.currentValue > asset.purchasePrice ? '+' : ''}{(((asset.currentValue - asset.purchasePrice) / asset.purchasePrice) * 100).toFixed(1)}%
                    </p>
                  )}
                </div>
              </div>
            ))}
            {(!assets || assets.length === 0) && (
              <div className="text-center py-4">
                <DollarSign className="w-8 h-8 text-midnight-500 mx-auto mb-2" />
                <p className="text-midnight-400">No assets yet</p>
              </div>
            )}
            {totalAssetsValue > 0 && (
              <div className="pt-3 border-t border-midnight-800/50 flex justify-between">
                <span className="text-midnight-400">Total Value</span>
                <span className="font-semibold">{formatCurrency(totalAssetsValue)}</span>
              </div>
            )}
          </div>
        </motion.div>

        {/* Recent Transactions */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="glass rounded-2xl p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Recent Transactions</h3>
            <Link to="/transactions" className="text-sm text-accent-violet hover:text-accent-violet/80">
              View all →
            </Link>
          </div>
          <div className="space-y-3">
            {recentTransactions.map((tx) => {
              const txType = getTransactionType(tx.type);
              return (
                <div key={tx.id} className="flex items-center justify-between p-3 rounded-lg bg-midnight-900/30">
                  <div className="flex items-center gap-3">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                      txType === 'income' ? 'bg-accent-emerald/20' :
                      txType === 'transfer' ? 'bg-accent-amber/20' : 'bg-accent-coral/20'
                    }`}>
                      {txType === 'income' ? (
                        <ArrowDownRight className="w-4 h-4 text-accent-emerald" />
                      ) : txType === 'transfer' ? (
                        <ArrowLeftRight className="w-4 h-4 text-accent-amber" />
                      ) : (
                        <ArrowUpRight className="w-4 h-4 text-accent-coral" />
                      )}
                    </div>
                    <p className="font-medium truncate max-w-[100px]">{tx.description || 'No description'}</p>
                  </div>
                  <p className={`font-semibold ${txType === 'income' ? 'text-accent-emerald' : 'text-accent-coral'}`}>
                    {txType === 'income' ? '+' : '-'}{formatCurrency(Math.abs(tx.amount), tx.currency)}
                  </p>
                </div>
              );
            })}
            {recentTransactions.length === 0 && (
              <p className="text-midnight-400 text-center py-4">No transactions yet</p>
            )}
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default Dashboard;
