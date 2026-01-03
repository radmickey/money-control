import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import {
  Plus,
  ArrowUpRight,
  ArrowDownRight,
  ArrowLeftRight,
  Search,
  Edit,
  Trash2,
  X,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchTransactions, createTransaction, updateTransaction, deleteTransaction, setFilters } from '../store/slices/transactionsSlice';

const categories = [
  'Food & Dining',
  'Transportation',
  'Shopping',
  'Entertainment',
  'Bills & Utilities',
  'Healthcare',
  'Travel',
  'Education',
  'Income',
  'Investment',
  'Transfer',
  'Other',
];

const typeMap: { [key: number]: string } = {
  0: 'unspecified',
  1: 'income',
  2: 'expense',
  3: 'transfer',
};

const categoryMap: { [key: number]: string } = {
  0: 'Other',
  1: 'Salary',
  2: 'Food & Dining',
  3: 'Transportation',
  4: 'Shopping',
  5: 'Entertainment',
  6: 'Healthcare',
  7: 'Bills & Utilities',
  8: 'Transfer',
  9: 'Other',
};

interface TransactionForm {
  amount: string;
  type: string;
  category: string;
  description: string;
  date: string;
  currency: string;
}

const defaultForm: TransactionForm = {
  amount: '',
  type: 'expense',
  category: 'Other',
  description: '',
  date: new Date().toISOString().split('T')[0],
  currency: 'USD',
};

const Transactions: React.FC = () => {
  const dispatch = useAppDispatch();
  const { transactions, filters, pagination, loading } = useAppSelector(
    (state) => state.transactions
  );
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<{ id: string; form: TransactionForm } | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [newTransaction, setNewTransaction] = useState<TransactionForm>(defaultForm);

  useEffect(() => {
    dispatch(fetchTransactions({ ...filters, page: pagination.page }));
  }, [dispatch, filters, pagination.page]);

  const handleCreateTransaction = async () => {
    await dispatch(
      createTransaction({
        ...newTransaction,
        amount: parseFloat(newTransaction.amount),
      })
    );
    setShowCreateModal(false);
    setNewTransaction(defaultForm);
  };

  const handleEditClick = (transaction: any) => {
    const txType = getTransactionType(transaction.type);
    const txCategory = getTransactionCategory(transaction.category);
    let dateStr = '';
    if (typeof transaction.date === 'object' && transaction.date.seconds) {
      dateStr = new Date(transaction.date.seconds * 1000).toISOString().split('T')[0];
    } else if (typeof transaction.date === 'string') {
      dateStr = transaction.date.split('T')[0];
    }
    
    setEditingTransaction({
      id: transaction.id,
      form: {
        amount: String(transaction.amount || 0),
        type: txType,
        category: txCategory,
        description: transaction.description || '',
        date: dateStr,
        currency: transaction.currency || 'USD',
      },
    });
    setShowEditModal(true);
  };

  const handleUpdateTransaction = async () => {
    if (editingTransaction) {
      await dispatch(
        updateTransaction({
          id: editingTransaction.id,
          data: {
            ...editingTransaction.form,
            amount: parseFloat(editingTransaction.form.amount),
          },
        })
      );
      setShowEditModal(false);
      setEditingTransaction(null);
    }
  };

  const handleDeleteTransaction = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this transaction?')) {
      await dispatch(deleteTransaction(id));
    }
  };

  const formatCurrency = (value: number, currency: string = 'USD') => {
    try {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: currency || 'USD',
      }).format(Math.abs(value));
    } catch {
      return `${currency} ${Math.abs(value).toFixed(2)}`;
    }
  };

  const formatDate = (dateValue: string | { seconds?: number; nanos?: number }) => {
    let date: Date;
    if (typeof dateValue === 'object' && dateValue.seconds) {
      date = new Date(dateValue.seconds * 1000);
    } else if (typeof dateValue === 'string') {
      date = new Date(dateValue);
    } else {
      return 'N/A';
    }
    if (isNaN(date.getTime())) return 'N/A';
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getTransactionType = (type: string | number): string => {
    if (typeof type === 'number') {
      return typeMap[type] || 'expense';
    }
    return type;
  };

  const getTransactionCategory = (category: string | number): string => {
    if (typeof category === 'number') {
      return categoryMap[category] || 'Other';
    }
    return category;
  };

  const getTypeIcon = (type: string | number) => {
    const typeStr = getTransactionType(type);
    switch (typeStr) {
      case 'income':
        return <ArrowDownRight className="w-5 h-5 text-accent-emerald" />;
      case 'expense':
        return <ArrowUpRight className="w-5 h-5 text-accent-coral" />;
      case 'transfer':
        return <ArrowLeftRight className="w-5 h-5 text-accent-amber" />;
      default:
        return <ArrowLeftRight className="w-5 h-5 text-midnight-400" />;
    }
  };

  const filteredTransactions = (transactions || []).filter((t) => {
    // Search filter
    const matchesSearch =
      (t.description || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      getTransactionCategory(t.category).toLowerCase().includes(searchTerm.toLowerCase());

    // Type filter
    const transactionType = getTransactionType(t.type);
    const matchesType = !filters.type || transactionType === filters.type;

    // Category filter
    const transactionCategory = getTransactionCategory(t.category);
    const matchesCategory = !filters.category || transactionCategory === filters.category;

    return matchesSearch && matchesType && matchesCategory;
  });

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-display font-bold text-gradient">Transactions</h1>
          <p className="text-midnight-400 mt-1">Track your income and expenses</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          Add Transaction
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-4">
        <div className="flex-1 min-w-[200px]">
          <div className="relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-midnight-400 pointer-events-none" />
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search transactions..."
              className="w-full px-4 py-3 pl-12 rounded-xl bg-zinc-900/80 border border-zinc-700/50 text-white placeholder-zinc-500 focus:outline-none focus:border-midnight-500 focus:ring-2 focus:ring-midnight-500/20 transition-all duration-200"
            />
          </div>
        </div>

        <select
          value={filters.type || ''}
          onChange={(e) => dispatch(setFilters({ type: e.target.value || undefined }))}
          className="input-field w-auto min-w-[140px]"
        >
          <option value="">All Types</option>
          <option value="income">Income</option>
          <option value="expense">Expense</option>
          <option value="transfer">Transfer</option>
        </select>

        <select
          value={filters.category || ''}
          onChange={(e) => dispatch(setFilters({ category: e.target.value || undefined }))}
          className="input-field w-auto min-w-[160px]"
        >
          <option value="">All Categories</option>
          {categories.map((cat) => (
            <option key={cat} value={cat}>
              {cat}
            </option>
          ))}
        </select>
      </div>

      {/* Transactions List */}
      <div className="glass rounded-2xl overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="w-12 h-12 border-4 border-midnight-500 border-t-transparent rounded-full animate-spin" />
          </div>
        ) : filteredTransactions.length === 0 ? (
          <div className="p-12 text-center">
            <ArrowLeftRight className="w-16 h-16 text-midnight-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold mb-2">No transactions found</h3>
            <p className="text-midnight-400">
              {searchTerm
                ? 'Try a different search term'
                : 'Add your first transaction to get started'}
            </p>
          </div>
        ) : (
          <div className="divide-y divide-midnight-800/50">
            {filteredTransactions.map((transaction, index) => (
              <motion.div
                key={transaction.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.02 * index }}
                className="p-4 hover:bg-midnight-800/20 transition-colors flex items-center gap-4 group"
              >
                <div className="w-10 h-10 rounded-xl bg-midnight-800/50 flex items-center justify-center">
                  {getTypeIcon(transaction.type)}
                </div>

                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{transaction.description || 'No description'}</p>
                  <p className="text-sm text-midnight-400">{getTransactionCategory(transaction.category)}</p>
                </div>

                <div className="text-right mr-2">
                  <p
                    className={`font-semibold ${
                      getTransactionType(transaction.type) === 'income'
                        ? 'text-accent-emerald'
                        : getTransactionType(transaction.type) === 'expense'
                        ? 'text-accent-coral'
                        : 'text-white'
                    }`}
                  >
                    {getTransactionType(transaction.type) === 'income' ? '+' : '-'}
                    {formatCurrency(transaction.amount, transaction.currency)}
                  </p>
                  <p className="text-sm text-midnight-400">{formatDate(transaction.date)}</p>
                </div>

                <div className="opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                  <button
                    onClick={() => handleEditClick(transaction)}
                    className="p-2 text-midnight-400 hover:text-white hover:bg-midnight-700/50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <Edit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDeleteTransaction(transaction.id)}
                    className="p-2 text-midnight-400 hover:text-accent-coral hover:bg-midnight-700/50 rounded-lg transition-colors"
                    title="Delete"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </motion.div>
            ))}
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
            <h3 className="text-xl font-semibold mb-6">Add Transaction</h3>

            <div className="space-y-4">
              <div className="flex gap-2">
                {['expense', 'income', 'transfer'].map((type) => (
                  <button
                    key={type}
                    onClick={() => setNewTransaction({ ...newTransaction, type })}
                    className={`flex-1 py-2 rounded-lg capitalize transition-colors ${
                      newTransaction.type === type
                        ? 'bg-midnight-500 text-white'
                        : 'bg-midnight-800/50 text-midnight-400 hover:text-white'
                    }`}
                  >
                    {type}
                  </button>
                ))}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Amount</label>
                  <input
                    type="number"
                    value={newTransaction.amount}
                    onChange={(e) =>
                      setNewTransaction({ ...newTransaction, amount: e.target.value })
                    }
                    className="input-field"
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Currency</label>
                  <select
                    value={newTransaction.currency}
                    onChange={(e) =>
                      setNewTransaction({ ...newTransaction, currency: e.target.value })
                    }
                    className="input-field"
                  >
                    <option value="USD">USD</option>
                    <option value="EUR">EUR</option>
                    <option value="GBP">GBP</option>
                    <option value="JPY">JPY</option>
                    <option value="RUB">RUB</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Description</label>
                <input
                  type="text"
                  value={newTransaction.description}
                  onChange={(e) =>
                    setNewTransaction({ ...newTransaction, description: e.target.value })
                  }
                  className="input-field"
                  placeholder="e.g., Coffee at Starbucks"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Category</label>
                <select
                  value={newTransaction.category}
                  onChange={(e) =>
                    setNewTransaction({ ...newTransaction, category: e.target.value })
                  }
                  className="input-field"
                >
                  {categories.map((cat) => (
                    <option key={cat} value={cat}>
                      {cat}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Date</label>
                <input
                  type="date"
                  value={newTransaction.date}
                  onChange={(e) =>
                    setNewTransaction({ ...newTransaction, date: e.target.value })
                  }
                  className="input-field"
                />
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
                onClick={handleCreateTransaction}
                disabled={!newTransaction.amount || !newTransaction.description}
                className="flex-1 btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Add
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* Edit Modal */}
      {showEditModal && editingTransaction && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="glass rounded-2xl p-6 w-full max-w-md"
          >
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold">Edit Transaction</h3>
              <button
                onClick={() => { setShowEditModal(false); setEditingTransaction(null); }}
                className="p-2 text-midnight-400 hover:text-white"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="flex gap-2">
                {['expense', 'income', 'transfer'].map((type) => (
                  <button
                    key={type}
                    onClick={() => setEditingTransaction({
                      ...editingTransaction,
                      form: { ...editingTransaction.form, type }
                    })}
                    className={`flex-1 py-2 rounded-lg capitalize transition-colors ${
                      editingTransaction.form.type === type
                        ? 'bg-midnight-500 text-white'
                        : 'bg-midnight-800/50 text-midnight-400 hover:text-white'
                    }`}
                  >
                    {type}
                  </button>
                ))}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Amount</label>
                  <input
                    type="number"
                    value={editingTransaction.form.amount}
                    onChange={(e) => setEditingTransaction({
                      ...editingTransaction,
                      form: { ...editingTransaction.form, amount: e.target.value }
                    })}
                    className="input-field"
                    placeholder="0.00"
                    step="0.01"
                  />
                </div>
                <div>
                  <label className="block text-sm text-midnight-400 mb-2">Currency</label>
                  <select
                    value={editingTransaction.form.currency}
                    onChange={(e) => setEditingTransaction({
                      ...editingTransaction,
                      form: { ...editingTransaction.form, currency: e.target.value }
                    })}
                    className="input-field"
                  >
                    <option value="USD">USD</option>
                    <option value="EUR">EUR</option>
                    <option value="GBP">GBP</option>
                    <option value="JPY">JPY</option>
                    <option value="RUB">RUB</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Description</label>
                <input
                  type="text"
                  value={editingTransaction.form.description}
                  onChange={(e) => setEditingTransaction({
                    ...editingTransaction,
                    form: { ...editingTransaction.form, description: e.target.value }
                  })}
                  className="input-field"
                  placeholder="e.g., Coffee at Starbucks"
                />
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Category</label>
                <select
                  value={editingTransaction.form.category}
                  onChange={(e) => setEditingTransaction({
                    ...editingTransaction,
                    form: { ...editingTransaction.form, category: e.target.value }
                  })}
                  className="input-field"
                >
                  {categories.map((cat) => (
                    <option key={cat} value={cat}>
                      {cat}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm text-midnight-400 mb-2">Date</label>
                <input
                  type="date"
                  value={editingTransaction.form.date}
                  onChange={(e) => setEditingTransaction({
                    ...editingTransaction,
                    form: { ...editingTransaction.form, date: e.target.value }
                  })}
                  className="input-field"
                />
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => { setShowEditModal(false); setEditingTransaction(null); }}
                className="flex-1 px-4 py-3 rounded-xl border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleUpdateTransaction}
                disabled={!editingTransaction.form.amount || !editingTransaction.form.description}
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

export default Transactions;

