import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Plus,
  TrendingUp,
  TrendingDown,
  Bitcoin,
  BarChart3,
  Home,
  RefreshCw,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchAssets, createAsset, fetchAssetPrice } from '../store/slices/assetsSlice';
import { formatCurrency, formatPercent } from '../utils/formatters';

const assetTypeIcons: { [key: string]: React.ElementType } = {
  stock: BarChart3,
  crypto: Bitcoin,
  etf: TrendingUp,
  real_estate: Home,
  other: BarChart3,
};

const Assets: React.FC = () => {
  const dispatch = useAppDispatch();
  const { assets, prices, loading, pricesLoading } = useAppSelector((state) => state.assets);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newAsset, setNewAsset] = useState({
    type: 'stock',
    symbol: '',
    name: '',
    quantity: '',
    purchasePrice: '',
  });

  useEffect(() => {
    dispatch(fetchAssets());
  }, [dispatch]);

  useEffect(() => {
    // Fetch prices for all assets
    if (assets && assets.length > 0) {
      assets.forEach((asset) => {
        if (asset.type !== 'real_estate' && asset.type !== 'other') {
          dispatch(fetchAssetPrice({ symbol: asset.symbol, type: asset.type }));
        }
      });
    }
  }, [dispatch, assets]);

  const handleCreateAsset = async () => {
    await dispatch(
      createAsset({
        ...newAsset,
        quantity: parseFloat(newAsset.quantity),
        purchasePrice: parseFloat(newAsset.purchasePrice),
      })
    );
    setShowCreateModal(false);
    setNewAsset({
      type: 'stock',
      symbol: '',
      name: '',
      quantity: '',
      purchasePrice: '',
    });
  };

  const handleRefreshPrices = () => {
    assets.forEach((asset) => {
      if (asset.type !== 'real_estate' && asset.type !== 'other') {
        dispatch(fetchAssetPrice({ symbol: asset.symbol, type: asset.type }));
      }
    });
  };

  const calculateGainLoss = (asset: any) => {
    const currentPrice = prices[asset.symbol] || asset.currentPrice;
    const currentValue = currentPrice * asset.quantity;
    const costBasis = asset.purchasePrice * asset.quantity;
    const gainLoss = currentValue - costBasis;
    const gainLossPercent = (gainLoss / costBasis) * 100;
    return { currentValue, gainLoss, gainLossPercent };
  };

  const totalValue = (assets || []).reduce((sum, asset) => {
    const { currentValue } = calculateGainLoss(asset);
    return sum + currentValue;
  }, 0);

  const totalGainLoss = (assets || []).reduce((sum, asset) => {
    const { gainLoss } = calculateGainLoss(asset);
    return sum + gainLoss;
  }, 0);

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-display font-bold text-gradient">Assets</h1>
          <p className="text-midnight-400 mt-1">Track stocks, crypto, and more</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={handleRefreshPrices}
            disabled={pricesLoading}
            className="px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors flex items-center gap-2"
          >
            <RefreshCw className={`w-4 h-4 ${pricesLoading ? 'animate-spin' : ''}`} />
            Refresh
          </button>
          <button
            onClick={() => setShowCreateModal(true)}
            className="btn-primary flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Add Asset
          </button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="glass rounded-2xl p-6"
        >
          <p className="text-midnight-400 text-sm">Total Value</p>
          <p className="text-3xl font-bold mt-2">{formatCurrency(totalValue)}</p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="glass rounded-2xl p-6"
        >
          <p className="text-midnight-400 text-sm">Total Gain/Loss</p>
          <p
            className={`text-3xl font-bold mt-2 ${
              totalGainLoss >= 0 ? 'text-accent-emerald' : 'text-accent-coral'
            }`}
          >
            {formatCurrency(totalGainLoss)}
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="glass rounded-2xl p-6"
        >
          <p className="text-midnight-400 text-sm">Assets Tracked</p>
          <p className="text-3xl font-bold mt-2">{(assets || []).length}</p>
        </motion.div>
      </div>

      {/* Assets List */}
      <div className="glass rounded-2xl overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="w-12 h-12 border-4 border-midnight-500 border-t-transparent rounded-full animate-spin" />
          </div>
        ) : !assets || assets.length === 0 ? (
          <div className="p-12 text-center">
            <BarChart3 className="w-16 h-16 text-midnight-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold mb-2">No assets yet</h3>
            <p className="text-midnight-400 mb-6">
              Add stocks, crypto, or other assets to track your portfolio
            </p>
            <button onClick={() => setShowCreateModal(true)} className="btn-primary">
              Add Asset
            </button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-midnight-800/50">
                  <th className="text-left p-4 text-midnight-400 font-medium">Asset</th>
                  <th className="text-right p-4 text-midnight-400 font-medium">Price</th>
                  <th className="text-right p-4 text-midnight-400 font-medium">Holdings</th>
                  <th className="text-right p-4 text-midnight-400 font-medium">Value</th>
                  <th className="text-right p-4 text-midnight-400 font-medium">Gain/Loss</th>
                </tr>
              </thead>
              <tbody>
                {(Array.isArray(assets) ? assets : []).map((asset, index) => {
                  const Icon = assetTypeIcons[asset.type] || BarChart3;
                  const { currentValue, gainLoss, gainLossPercent } = calculateGainLoss(asset);
                  const currentPrice = prices[asset.symbol] || asset.currentPrice;

                  return (
                    <motion.tr
                      key={asset.id}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: 0.02 * index }}
                      className="border-b border-midnight-800/30 hover:bg-midnight-800/20 transition-colors"
                    >
                      <td className="p-4">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-lg bg-midnight-800/50 flex items-center justify-center">
                            <Icon className="w-5 h-5 text-midnight-300" />
                          </div>
                          <div>
                            <p className="font-medium">{asset.symbol}</p>
                            <p className="text-sm text-midnight-400">{asset.name}</p>
                          </div>
                        </div>
                      </td>
                      <td className="p-4 text-right">{formatCurrency(currentPrice)}</td>
                      <td className="p-4 text-right font-mono">{asset.quantity}</td>
                      <td className="p-4 text-right font-semibold">{formatCurrency(currentValue)}</td>
                      <td className="p-4 text-right">
                        <div
                          className={`flex items-center justify-end gap-1 ${
                            gainLoss >= 0 ? 'text-accent-emerald' : 'text-accent-coral'
                          }`}
                        >
                          {gainLoss >= 0 ? (
                            <TrendingUp className="w-4 h-4" />
                          ) : (
                            <TrendingDown className="w-4 h-4" />
                          )}
                          <span>{formatPercent(gainLossPercent)}</span>
                        </div>
                        <p className="text-sm text-midnight-400 mt-1">
                          {formatCurrency(gainLoss)}
                        </p>
                      </td>
                    </motion.tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <h3 className="text-xl font-semibold mb-6">Add Asset</h3>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">Asset Type</label>
                <select
                  value={newAsset.type}
                  onChange={(e) => setNewAsset({ ...newAsset, type: e.target.value })}
                  className="input-field"
                >
                  <option value="stock">Stock</option>
                  <option value="crypto">Cryptocurrency</option>
                  <option value="etf">ETF</option>
                  <option value="real_estate">Real Estate</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Symbol</label>
                <input
                  type="text"
                  value={newAsset.symbol}
                  onChange={(e) =>
                    setNewAsset({ ...newAsset, symbol: e.target.value.toUpperCase() })
                  }
                  className="input-field"
                  placeholder="e.g., AAPL, BTC"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Name</label>
                <input
                  type="text"
                  value={newAsset.name}
                  onChange={(e) => setNewAsset({ ...newAsset, name: e.target.value })}
                  className="input-field"
                  placeholder="e.g., Apple Inc."
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Quantity</label>
                  <input
                    type="number"
                    value={newAsset.quantity}
                    onChange={(e) => setNewAsset({ ...newAsset, quantity: e.target.value })}
                    className="input-field"
                    placeholder="0"
                    step="any"
                  />
                </div>
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Purchase Price</label>
                  <input
                    type="number"
                    value={newAsset.purchasePrice}
                    onChange={(e) =>
                      setNewAsset({ ...newAsset, purchasePrice: e.target.value })
                    }
                    className="input-field"
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowCreateModal(false)}
                className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateAsset}
                disabled={!newAsset.symbol || !newAsset.quantity}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Add
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
};

export default Assets;

