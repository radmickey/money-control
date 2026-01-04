import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Plus,
  Wallet,
  MoreVertical,
  Edit,
  Trash2,
  ChevronDown,
  ChevronUp,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchAccounts, createAccount, updateAccount, deleteAccount, createSubAccount, updateSubAccount, deleteSubAccount } from '../store/slices/accountsSlice';
import { formatCurrency, formatAccountType, getAccountTypeKey } from '../utils/formatters';
import { ACCOUNT_TYPE_ICONS, ACCOUNT_TYPES } from '../constants';
import Modal, { FormField, CancelButton, SubmitButton, CurrencySelect } from '../components/common/Modal';

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

  const getIcon = (type: string | number): React.ElementType => {
    return ACCOUNT_TYPE_ICONS[getAccountTypeKey(type)] || Wallet;
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
      <Modal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        title="Create Account"
        footer={
          <>
            <CancelButton onClick={() => setShowCreateModal(false)} />
            <SubmitButton onClick={handleCreateAccount} disabled={!newAccount.name}>
              Create
            </SubmitButton>
          </>
        }
      >
        <FormField label="Account Name">
          <input
            type="text"
            value={newAccount.name}
            onChange={(e) => setNewAccount({ ...newAccount, name: e.target.value })}
            className="input-field"
            placeholder="e.g., Main Bank Account"
          />
        </FormField>
        <FormField label="Type">
          <select
            value={newAccount.type}
            onChange={(e) => setNewAccount({ ...newAccount, type: e.target.value })}
            className="input-field"
          >
            {ACCOUNT_TYPES.map(t => (
              <option key={t.value} value={t.value}>{t.label}</option>
            ))}
          </select>
        </FormField>
        <FormField label="Currency">
          <CurrencySelect
            value={newAccount.currency}
            onChange={(v) => setNewAccount({ ...newAccount, currency: v })}
          />
        </FormField>
      </Modal>

      {/* Edit Account Modal */}
      <Modal
        isOpen={showEditModal && !!editingAccount}
        onClose={() => { setShowEditModal(false); setEditingAccount(null); }}
        title="Edit Account"
        footer={
          <>
            <CancelButton onClick={() => { setShowEditModal(false); setEditingAccount(null); }} />
            <SubmitButton onClick={handleEditAccount} disabled={!editingAccount?.name}>
              Save Changes
            </SubmitButton>
          </>
        }
      >
        {editingAccount && (
          <FormField label="Account Name">
            <input
              type="text"
              value={editingAccount.name}
              onChange={(e) => setEditingAccount({ ...editingAccount, name: e.target.value })}
              className="input-field"
            />
          </FormField>
        )}
      </Modal>

      {/* Create Sub-Account Modal */}
      <Modal
        isOpen={showSubAccountModal}
        onClose={() => setShowSubAccountModal(false)}
        title="Add Sub-Account"
        footer={
          <>
            <CancelButton onClick={() => setShowSubAccountModal(false)} />
            <SubmitButton onClick={handleCreateSubAccount} disabled={!newSubAccount.name}>
              Add
            </SubmitButton>
          </>
        }
      >
        <FormField label="Name">
          <input
            type="text"
            value={newSubAccount.name}
            onChange={(e) => setNewSubAccount({ ...newSubAccount, name: e.target.value })}
            className="input-field"
            placeholder="e.g., Savings, USD Account"
          />
        </FormField>
        <FormField label="Initial Balance">
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
        </FormField>
        <FormField label="Currency">
          <CurrencySelect
            value={newSubAccount.currency}
            onChange={(v) => setNewSubAccount({ ...newSubAccount, currency: v })}
          />
        </FormField>
      </Modal>

      {/* Edit Sub-Account Modal */}
      <Modal
        isOpen={showEditSubAccountModal && !!editingSubAccount}
        onClose={() => { setShowEditSubAccountModal(false); setEditingSubAccount(null); }}
        title="Edit Sub-Account"
        footer={
          <>
            <CancelButton onClick={() => { setShowEditSubAccountModal(false); setEditingSubAccount(null); }} />
            <SubmitButton onClick={handleUpdateSubAccount} disabled={!editingSubAccount?.name}>
              Save Changes
            </SubmitButton>
          </>
        }
      >
        {editingSubAccount && (
          <>
            <FormField label="Name">
              <input
                type="text"
                value={editingSubAccount.name}
                onChange={(e) => setEditingSubAccount({ ...editingSubAccount, name: e.target.value })}
                className="input-field"
              />
            </FormField>
            <FormField label="Balance">
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
            </FormField>
          </>
        )}
      </Modal>
    </div>
  );
};

export default Accounts;
