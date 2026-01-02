import React, { useState } from 'react';
import { View, Text, ScrollView, StyleSheet, TouchableOpacity, TextInput } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';

const mockTransactions = [
  { id: '1', description: 'Salary', amount: 5000, type: 'income', category: 'Income', date: '2024-01-15' },
  { id: '2', description: 'Groceries', amount: -150, type: 'expense', category: 'Food', date: '2024-01-14' },
  { id: '3', description: 'Netflix', amount: -15.99, type: 'expense', category: 'Entertainment', date: '2024-01-13' },
  { id: '4', description: 'Transfer to Savings', amount: -500, type: 'transfer', category: 'Transfer', date: '2024-01-12' },
];

export default function TransactionsScreen() {
  const [searchQuery, setSearchQuery] = useState('');

  const formatCurrency = (value: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(Math.abs(value));
  const formatDate = (dateStr: string) => new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

  const getTypeIcon = (type: string): keyof typeof Ionicons.glyphMap => {
    if (type === 'income') return 'arrow-down';
    if (type === 'expense') return 'arrow-up';
    return 'swap-horizontal';
  };

  const getTypeColor = (type: string) => {
    if (type === 'income') return '#10b981';
    if (type === 'expense') return '#f43f5e';
    return '#f59e0b';
  };

  return (
    <SafeAreaView style={styles.container} edges={['left', 'right']}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={styles.title}>Transactions</Text>
          <TouchableOpacity style={styles.addButton}>
            <Ionicons name="add" size={24} color="#fff" />
          </TouchableOpacity>
        </View>

        <View style={styles.searchContainer}>
          <Ionicons name="search" size={20} color="#8eafff" />
          <TextInput style={styles.searchInput} placeholder="Search transactions..." placeholderTextColor="#8eafff88" value={searchQuery} onChangeText={setSearchQuery} />
        </View>

        <View style={styles.filterRow}>
          {['All', 'Income', 'Expense', 'Transfer'].map((filter) => (
            <TouchableOpacity key={filter} style={[styles.filterChip, filter === 'All' && styles.filterChipActive]}>
              <Text style={[styles.filterChipText, filter === 'All' && styles.filterChipTextActive]}>{filter}</Text>
            </TouchableOpacity>
          ))}
        </View>

        {mockTransactions.map((tx) => (
          <TouchableOpacity key={tx.id} style={styles.transactionItem}>
            <View style={[styles.txIcon, { backgroundColor: getTypeColor(tx.type) + '22' }]}>
              <Ionicons name={getTypeIcon(tx.type)} size={20} color={getTypeColor(tx.type)} />
            </View>
            <View style={styles.txInfo}>
              <Text style={styles.txDescription}>{tx.description}</Text>
              <Text style={styles.txCategory}>{tx.category}</Text>
            </View>
            <View style={styles.txAmountContainer}>
              <Text style={[styles.txAmount, { color: getTypeColor(tx.type) }]}>
                {tx.type === 'income' ? '+' : '-'}{formatCurrency(tx.amount)}
              </Text>
              <Text style={styles.txDate}>{formatDate(tx.date)}</Text>
            </View>
          </TouchableOpacity>
        ))}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#0c0d1f' },
  scrollView: { flex: 1, padding: 24 },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 },
  title: { fontSize: 24, fontWeight: 'bold', color: '#fff' },
  addButton: { width: 44, height: 44, borderRadius: 12, backgroundColor: '#2d4fff', justifyContent: 'center', alignItems: 'center' },
  searchContainer: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#161a8f22', borderRadius: 12, paddingHorizontal: 16, marginBottom: 16, borderWidth: 1, borderColor: '#2d4fff33' },
  searchInput: { flex: 1, height: 48, color: '#fff', fontSize: 16, marginLeft: 12 },
  filterRow: { flexDirection: 'row', gap: 8, marginBottom: 24 },
  filterChip: { paddingHorizontal: 16, paddingVertical: 8, borderRadius: 20, backgroundColor: '#161a8f22', borderWidth: 1, borderColor: '#2d4fff33' },
  filterChipActive: { backgroundColor: '#2d4fff', borderColor: '#2d4fff' },
  filterChipText: { color: '#8eafff', fontSize: 14 },
  filterChipTextActive: { color: '#fff' },
  transactionItem: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#161a8f22', borderRadius: 16, padding: 16, marginBottom: 12, borderWidth: 1, borderColor: '#2d4fff22' },
  txIcon: { width: 44, height: 44, borderRadius: 12, justifyContent: 'center', alignItems: 'center', marginRight: 12 },
  txInfo: { flex: 1 },
  txDescription: { fontSize: 16, fontWeight: '500', color: '#fff', marginBottom: 2 },
  txCategory: { fontSize: 14, color: '#8eafff' },
  txAmountContainer: { alignItems: 'flex-end' },
  txAmount: { fontSize: 16, fontWeight: '600' },
  txDate: { fontSize: 12, color: '#8eafff', marginTop: 2 },
});

