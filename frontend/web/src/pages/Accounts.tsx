import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Plus,
  Wallet,
  Building2,
  Banknote,
  TrendingUp,
  Bitcoin,
  Home,
  MoreVertical,
  Edit,
  Trash2,
  ChevronDown,
  ChevronUp,
  X,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchAccounts, createAccount, updateAccount, deleteAccount, createSubAccount, updateSubAccount, deleteSubAccount } from '../store/slices/accountsSlice';

const accountTypeIcons: { [key: string]: React.ElementType } = {
  bank: Building2,
  cash: Banknote,
  investment: TrendingUp,
  crypto: Bitcoin,
  real_estate: Home,
  other: Wallet,
};

const Accounts: React.FC = () => {
  const dispatch = useAppDispatch();
  const { accounts, loading } = useAppSelector((state) => state.accounts);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showSubAccountModal, setShowSubAccountModal] = useState(false);
  const [showEditSubAccountModal, setShowEditSubAccountModal] = useState(false);
  const [editingAccount, setEditingAccount] = useState<{ id: string; name: string } | null>(null);
  const [expandedAccounts, setExpandedAccounts] = useState<Set<string>>(new Set());
  const [selectedAccountId, setSelectedAccountId] = useState<string | null>(null);
  const [editingSubAccount, setEditingSubAccount] = useState<{ accountId: string; id: string; name: string; balance: number } | null>(null);
  const [newAccount, setNewAccount] = useState({
    name: '',
    type: 'bank',
    currency: 'USD',
  });
  const [newSubAccount, setNewSubAccount] = useState({
    name: '',
    balance: '',
    currency: 'USD',
  });

  useEffect(() => {
    dispatch(fetchAccounts());
  }, [dispatch]);

  const toggleExpand = (accountId: string) => {
    setExpandedAccounts(prev => {
      const next = new Set(prev);
      if (next.has(accountId)) {
        next.delete(accountId);
      } else {
        next.add(accountId);
      }
      return next;
    });
  };

  const handleCreateAccount = async () => {
    await dispatch(createAccount(newAccount));
    setShowCreateModal(false);
    setNewAccount({ name: '', type: 'bank', currency: 'USD' });
  };

  const handleEditClick = (account: { id: string; name: string }) => {
    setEditingAccount({ id: account.id, name: account.name });
    setShowEditModal(true);
  };

  const handleEditAccount = async () => {
    if (editingAccount) {
      await dispatch(updateAccount({ id: editingAccount.id, data: { name: editingAccount.name } }));
      setShowEditModal(false);
      setEditingAccount(null);
    }
  };

  const handleDeleteAccount = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this account?')) {
      dispatch(deleteAccount(id));
    }
  };

  const handleAddSubAccount = (accountId: string) => {
    setSelectedAccountId(accountId);
    const account = (accounts || []).find(a => a.id === accountId);
    setNewSubAccount({ name: '', balance: '', currency: account?.currency || 'USD' });
    setShowSubAccountModal(true);
  };

  const handleCreateSubAccount = async () => {
    if (selectedAccountId && newSubAccount.name) {
      await dispatch(createSubAccount({
        accountId: selectedAccountId,
        data: {
          name: newSubAccount.name,
          balance: parseFloat(newSubAccount.balance) || 0,
          currency: newSubAccount.currency,
        },
      }));
      setShowSubAccountModal(false);
      setNewSubAccount({ name: '', balance: '', currency: 'USD' });
      setSelectedAccountId(null);
    }
  };

  const handleEditSubAccount = (accountId: string, sub: { id: string; name: string; balance: number }) => {
    setEditingSubAccount({ accountId, id: sub.id, name: sub.name, balance: sub.balance });
    setShowEditSubAccountModal(true);
  };

  const handleUpdateSubAccount = async () => {
    if (editingSubAccount) {
      await dispatch(updateSubAccount({
        accountId: editingSubAccount.accountId,
        subAccountId: editingSubAccount.id,
        data: {
          name: editingSubAccount.name,
          balance: editingSubAccount.balance,
        },
      }));
      setShowEditSubAccountModal(false);
      setEditingSubAccount(null);
    }
  };

  const handleDeleteSubAccount = async (accountId: string, subAccountId: string) => {
    if (window.confirm('Are you sure you want to delete this sub-account?')) {
      dispatch(deleteSubAccount({ accountId, subAccountId }));
    }
  };

  const formatCurrency = (value: number | undefined | null, currency: string) => {
    const safeValue = value ?? 0;
    const safeCurrency = currency || 'USD';
    try {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: safeCurrency,
        minimumFractionDigits: 0,
        maximumFractionDigits: 2,
      }).format(safeValue);
    } catch {
      return `${safeCurrency} ${safeValue.toFixed(2)}`;
    }
  };

  const getIcon = (type: string | number) => {
    // Map numeric type to string
    const typeMap: { [key: number]: string } = {
      0: 'unspecified',
      1: 'bank',
      2: 'cash',
      3: 'investment',
      4: 'crypto',
      5: 'real_estate',
      6: 'other',
    };
    const typeString = typeof type === 'number' ? (typeMap[type] || 'other') : type;
    const Icon = accountTypeIcons[typeString.toLowerCase()] || Wallet;
    return Icon;
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

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-display font-bold text-gradient">Accounts</h1>
          <p className="text-midnight-400 mt-1">Manage your accounts and sub-accounts</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          Add Account
        </button>
      </div>

      {/* Accounts Grid */}
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <div className="w-12 h-12 border-4 border-midnight-500 border-t-transparent rounded-full animate-spin" />
        </div>
      ) : !accounts || accounts.length === 0 ? (
        <div className="glass rounded-2xl p-12 text-center">
          <Wallet className="w-16 h-16 text-midnight-400 mx-auto mb-4" />
          <h3 className="text-xl font-semibold mb-2">No accounts yet</h3>
          <p className="text-midnight-400 mb-6">
            Create your first account to start tracking your finances
          </p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="btn-primary"
          >
            Create Account
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {(Array.isArray(accounts) ? accounts : []).map((account, index) => {
            const Icon = getIcon(account.type);
            return (
              <motion.div
                key={account.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.05 * index }}
                className="glass rounded-2xl p-6 hover:border-midnight-500/30 transition-all duration-200 group"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-midnight-600 to-midnight-800 flex items-center justify-center">
                    <Icon className="w-6 h-6 text-midnight-300" />
                  </div>
                  <div className="relative">
                    <button className="p-2 text-midnight-400 hover:text-white opacity-0 group-hover:opacity-100 transition-opacity">
                      <MoreVertical className="w-5 h-5" />
                    </button>
                    <div className="absolute right-0 top-full mt-1 w-36 glass rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10">
                      <button
                        onClick={() => handleEditClick(account)}
                        className="w-full flex items-center gap-2 px-4 py-2 text-sm text-left hover:bg-midnight-800/50 rounded-t-lg"
                      >
                        <Edit className="w-4 h-4" />
                        Edit
                      </button>
                      <button
                        onClick={() => handleDeleteAccount(account.id)}
                        className="w-full flex items-center gap-2 px-4 py-2 text-sm text-left text-accent-coral hover:bg-midnight-800/50 rounded-b-lg"
                      >
                        <Trash2 className="w-4 h-4" />
                        Delete
                      </button>
                    </div>
                  </div>
                </div>

                <h3 className="text-lg font-semibold">{account.name}</h3>
                <p className="text-midnight-400 text-sm mb-4">{formatAccountType(account.type)}</p>

                <div className="pt-4 border-t border-midnight-800/50">
                  <p className="text-sm text-midnight-400">Total Balance</p>
                  <p className="text-2xl font-semibold mt-1">
                    {formatCurrency(account.totalBalance, account.currency)}
                  </p>
                </div>

                {/* Sub-accounts section */}
                <div className="mt-4 pt-4 border-t border-midnight-800/50">
                  <div
                    className="flex items-center justify-between cursor-pointer"
                    onClick={() => toggleExpand(account.id)}
                  >
                    <p className="text-xs text-midnight-400 uppercase tracking-wider">
                      Sub-accounts ({(account.subAccounts || []).length})
                    </p>
                    {expandedAccounts.has(account.id) ? (
                      <ChevronUp className="w-4 h-4 text-midnight-400" />
                    ) : (
                      <ChevronDown className="w-4 h-4 text-midnight-400" />
                    )}
                  </div>

                  {expandedAccounts.has(account.id) && (
                    <div className="mt-3 space-y-2">
                      {(account.subAccounts || []).map((sub) => (
                        <div key={sub.id} className="flex items-center justify-between text-sm p-2 rounded-lg bg-midnight-800/30 group/sub">
                          <span className="text-midnight-300">{sub.name}</span>
                          <div className="flex items-center gap-2">
                            <span className="font-medium">{formatCurrency(sub.balance, sub.currency)}</span>
                            <div className="opacity-0 group-hover/sub:opacity-100 transition-opacity flex gap-1">
                              <button
                                onClick={(e) => { e.stopPropagation(); handleEditSubAccount(account.id, sub); }}
                                className="p-1 text-midnight-400 hover:text-white"
                              >
                                <Edit className="w-3 h-3" />
                              </button>
                              <button
                                onClick={(e) => { e.stopPropagation(); handleDeleteSubAccount(account.id, sub.id); }}
                                className="p-1 text-accent-coral hover:text-red-400"
                              >
                                <Trash2 className="w-3 h-3" />
                              </button>
                            </div>
                          </div>
                        </div>
                      ))}

                      <button
                        onClick={(e) => { e.stopPropagation(); handleAddSubAccount(account.id); }}
                        className="w-full py-2 px-3 rounded-lg border border-dashed border-midnight-600 text-midnight-400 hover:text-white hover:border-midnight-500 transition-colors text-sm flex items-center justify-center gap-2"
                      >
                        <Plus className="w-4 h-4" />
                        Add Sub-account
                      </button>
                    </div>
                  )}
                </div>
              </motion.div>
            );
          })}
        </div>
      )}

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <h3 className="text-xl font-semibold mb-6">Create New Account</h3>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">Account Name</label>
                <input
                  type="text"
                  value={newAccount.name}
                  onChange={(e) => setNewAccount({ ...newAccount, name: e.target.value })}
                  className="input-field"
                  placeholder="e.g., Main Bank Account"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Account Type</label>
                <select
                  value={newAccount.type}
                  onChange={(e) => setNewAccount({ ...newAccount, type: e.target.value })}
                  className="input-field"
                >
                  <option value="bank">Bank Account</option>
                  <option value="cash">Cash</option>
                  <option value="investment">Investment</option>
                  <option value="crypto">Cryptocurrency</option>
                  <option value="real_estate">Real Estate</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Currency</label>
                <select
                  value={newAccount.currency}
                  onChange={(e) => setNewAccount({ ...newAccount, currency: e.target.value })}
                  className="input-field"
                >
                  <option value="USD">USD - US Dollar</option>
                  <option value="EUR">EUR - Euro</option>
                  <option value="GBP">GBP - British Pound</option>
                  <option value="JPY">JPY - Japanese Yen</option>
                  <option value="RUB">RUB - Russian Ruble</option>
                </select>
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
                onClick={handleCreateAccount}
                disabled={!newAccount.name}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Create
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Edit Account Modal */}
      {showEditModal && editingAccount && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <h3 className="text-xl font-semibold mb-6">Edit Account</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">
                  Account Name
                </label>
                <input
                  type="text"
                  value={editingAccount.name}
                  onChange={(e) => setEditingAccount({ ...editingAccount, name: e.target.value })}
                  className="w-full px-4 py-3 rounded-xl bg-midnight-900/50 border border-midnight-700 focus:border-midnight-500 focus:outline-none"
                  placeholder="Enter account name"
                />
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => {
                  setShowEditModal(false);
                  setEditingAccount(null);
                }}
                className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleEditAccount}
                disabled={!editingAccount.name}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Save Changes
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Create Sub-Account Modal */}
      {showSubAccountModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">Add Sub-Account</h3>
              <button
                onClick={() => setShowSubAccountModal(false)}
                className="p-2 text-midnight-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">Name</label>
                <input
                  type="text"
                  value={newSubAccount.name}
                  onChange={(e) => setNewSubAccount({ ...newSubAccount, name: e.target.value })}
                  className="input-field"
                  placeholder="e.g., USD Balance, Savings"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Initial Balance</label>
                <input
                  type="number"
                  value={newSubAccount.balance}
                  onChange={(e) => setNewSubAccount({ ...newSubAccount, balance: e.target.value })}
                  className="input-field"
                  placeholder="0.00"
                  step="0.01"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Currency</label>
                <select
                  value={newSubAccount.currency}
                  onChange={(e) => setNewSubAccount({ ...newSubAccount, currency: e.target.value })}
                  className="input-field"
                >
                  <option value="USD">USD - US Dollar</option>
                  <option value="EUR">EUR - Euro</option>
                  <option value="GBP">GBP - British Pound</option>
                  <option value="JPY">JPY - Japanese Yen</option>
                  <option value="RUB">RUB - Russian Ruble</option>
                </select>
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowSubAccountModal(false)}
                className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateSubAccount}
                disabled={!newSubAccount.name}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Create
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Edit Sub-Account Modal */}
      {showEditSubAccountModal && editingSubAccount && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">Edit Sub-Account</h3>
              <button
                onClick={() => { setShowEditSubAccountModal(false); setEditingSubAccount(null); }}
                className="p-2 text-midnight-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">Name</label>
                <input
                  type="text"
                  value={editingSubAccount.name}
                  onChange={(e) => setEditingSubAccount({ ...editingSubAccount, name: e.target.value })}
                  className="input-field"
                  placeholder="Sub-account name"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Balance</label>
                <input
                  type="number"
                  value={editingSubAccount.balance}
                  onChange={(e) => setEditingSubAccount({ ...editingSubAccount, balance: parseFloat(e.target.value) || 0 })}
                  className="input-field"
                  placeholder="0.00"
                  step="0.01"
                />
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => { setShowEditSubAccountModal(false); setEditingSubAccount(null); }}
                className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleUpdateSubAccount}
                disabled={!editingSubAccount.name}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Save Changes
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
};

export default Accounts;

