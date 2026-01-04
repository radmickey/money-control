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
  const [editingSubAccount, setEditingSubAccount] = useState<{ accountId: string; id: string; name: string; balance: string } | null>(null);
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
    if (confirm('Are you sure you want to delete this account?')) {
      await dispatch(deleteAccount(id));
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
    setEditingSubAccount({ accountId, id: sub.id, name: sub.name, balance: sub.balance.toString() });
    setShowEditSubAccountModal(true);
  };

  const handleUpdateSubAccount = async () => {
    if (editingSubAccount) {
      await dispatch(updateSubAccount({
        accountId: editingSubAccount.accountId,
        subAccountId: editingSubAccount.id,
        data: {
          name: editingSubAccount.name,
          balance: parseFloat(editingSubAccount.balance) || 0,
        },
      }));
      setShowEditSubAccountModal(false);
      setEditingSubAccount(null);
    }
  };

  const handleDeleteSubAccount = async (accountId: string, subAccountId: string) => {
    if (confirm('Are you sure you want to delete this sub-account?')) {
      await dispatch(deleteSubAccount({ accountId, subAccountId }));
    }
  };

  // Simple currency formatting - no conversion logic needed
  const formatCurrency = (amount: number, currency: string): string => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount);
  };

  const getIcon = (type: string | number): React.ElementType => {
    const typeStr = typeof type === 'number' 
      ? ['other', 'bank', 'cash', 'investment', 'crypto', 'real_estate', 'other'][type] || 'other'
      : type;
    return accountTypeIcons[typeStr] || Wallet;
  };

  const formatAccountType = (type: string | number): string => {
    const typeMap: { [key: number]: string } = {
      0: 'Other',
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
            // Use pre-calculated values from backend (already mapped by Redux)
            const displayCurrency = account.displayCurrency || account.currency || 'USD';
            const totalBalance = account.convertedTotalBalance ?? account.totalBalance ?? 0;
            const isMixed = account.isMixedCurrency ?? false;
            
            return (
              <motion.div
                key={account.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: index * 0.05 }}
                className="glass p-6 rounded-2xl relative"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-accent-purple/20 to-accent-blue/20 flex items-center justify-center">
                    <Icon className="w-6 h-6 text-accent-purple" />
                  </div>
                  <div className="relative group">
                    <button className="p-2 rounded-lg hover:bg-midnight-800/50 transition-colors">
                      <MoreVertical className="w-5 h-5 text-midnight-400" />
                    </button>
                    <div className="absolute right-0 top-full mt-1 w-32 bg-midnight-900 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10 border border-midnight-700">
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
                  <div>
                    <p className="text-2xl font-semibold mt-1">
                      {isMixed ? '~' : ''}{formatCurrency(totalBalance, displayCurrency)}
                    </p>
                    {isMixed && (
                      <p className="text-xs text-midnight-500 mt-1">Mixed currencies</p>
                    )}
                  </div>
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
                    <motion.div
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      exit={{ opacity: 0, height: 0 }}
                      className="mt-3 space-y-2"
                    >
                      {(account.subAccounts || []).map((sub: any) => (
                        <div
                          key={sub.id}
                          className="flex items-center justify-between p-3 rounded-lg bg-midnight-800/30"
                        >
                          <div className="text-sm">{sub.name}</div>
                          <div className="flex items-center gap-2">
                            <div className="text-sm font-medium">
                              {formatCurrency(sub.balance, sub.currency)}
                            </div>
                            <div className="flex gap-1">
                              <button
                                onClick={() => handleEditSubAccount(account.id, sub)}
                                className="p-1 rounded hover:bg-midnight-700/50"
                              >
                                <Edit className="w-3 h-3 text-midnight-400" />
                              </button>
                              <button
                                onClick={() => handleDeleteSubAccount(account.id, sub.id)}
                                className="p-1 rounded hover:bg-midnight-700/50"
                              >
                                <Trash2 className="w-3 h-3 text-accent-coral" />
                              </button>
                            </div>
                          </div>
                        </div>
                      ))}
                      <button
                        onClick={() => handleAddSubAccount(account.id)}
                        className="w-full flex items-center justify-center gap-2 p-2 rounded-lg border border-dashed border-midnight-600 text-midnight-400 hover:border-accent-purple hover:text-accent-purple transition-colors text-sm"
                      >
                        <Plus className="w-4 h-4" />
                        Add Sub-account
                      </button>
                    </motion.div>
                  )}
                </div>
              </motion.div>
            );
          })}
        </div>
      )}

      {/* Create Account Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">Create Account</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="p-2 text-midnight-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

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
                <label className="block text-sm text-midnight-400 mb-2">Type</label>
                <select
                  value={newAccount.type}
                  onChange={(e) => setNewAccount({ ...newAccount, type: e.target.value })}
                  className="input-field"
                >
                  <option value="bank">Bank</option>
                  <option value="cash">Cash</option>
                  <option value="investment">Investment</option>
                  <option value="crypto">Crypto</option>
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
                  <option value="RUB">RUB - Russian Ruble</option>
                  <option value="JPY">JPY - Japanese Yen</option>
                  <option value="CHF">CHF - Swiss Franc</option>
                  <option value="CAD">CAD - Canadian Dollar</option>
                  <option value="AUD">AUD - Australian Dollar</option>
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
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">Edit Account</h3>
              <button
                onClick={() => { setShowEditModal(false); setEditingAccount(null); }}
                className="p-2 text-midnight-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm text-midnight-400 mb-2">Account Name</label>
                <input
                  type="text"
                  value={editingAccount.name}
                  onChange={(e) => setEditingAccount({ ...editingAccount, name: e.target.value })}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => { setShowEditModal(false); setEditingAccount(null); }}
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
                  placeholder="e.g., Savings, USD Account"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Initial Balance</label>
                <input
                  type="text"
                  inputMode="decimal"
                  value={newSubAccount.balance}
                  onChange={(e) => {
                    const value = e.target.value;
                    if (value === '' || /^\d*\.?\d*$/.test(value)) {
                      setNewSubAccount({ ...newSubAccount, balance: value });
                    }
                  }}
                  className="input-field"
                  placeholder="Enter amount"
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
                  <option value="RUB">RUB - Russian Ruble</option>
                  <option value="JPY">JPY - Japanese Yen</option>
                  <option value="CHF">CHF - Swiss Franc</option>
                  <option value="CAD">CAD - Canadian Dollar</option>
                  <option value="AUD">AUD - Australian Dollar</option>
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
                Add
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
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Balance</label>
                <input
                  type="text"
                  inputMode="decimal"
                  value={editingSubAccount.balance}
                  onChange={(e) => {
                    const value = e.target.value;
                    if (value === '' || /^\d*\.?\d*$/.test(value)) {
                      setEditingSubAccount({ ...editingSubAccount, balance: value });
                    }
                  }}
                  className="input-field"
                  placeholder="0.00"
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
