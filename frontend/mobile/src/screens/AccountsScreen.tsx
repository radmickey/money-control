import React, { useEffect } from 'react';
import { View, Text, ScrollView, StyleSheet, TouchableOpacity, RefreshControl } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchAccounts } from '../store/slices/accountsSlice';

const accountTypeIcons: { [key: string]: keyof typeof Ionicons.glyphMap } = {
  bank: 'business', cash: 'cash', investment: 'trending-up', crypto: 'logo-bitcoin', real_estate: 'home', other: 'wallet',
};

export default function AccountsScreen() {
  const dispatch = useAppDispatch();
  const { accounts, loading } = useAppSelector((state) => state.accounts);

  useEffect(() => { dispatch(fetchAccounts()); }, []);

  const formatCurrency = (value: number, currency: string) => new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(value);

  return (
    <SafeAreaView style={styles.container} edges={['left', 'right']}>
      <ScrollView style={styles.scrollView} refreshControl={<RefreshControl refreshing={loading} onRefresh={() => dispatch(fetchAccounts())} tintColor="#2d4fff" />}>
        <View style={styles.header}>
          <Text style={styles.title}>Your Accounts</Text>
          <TouchableOpacity style={styles.addButton}>
            <Ionicons name="add" size={24} color="#fff" />
          </TouchableOpacity>
        </View>

        {accounts.length === 0 ? (
          <View style={styles.emptyState}>
            <Ionicons name="wallet-outline" size={64} color="#8eafff" />
            <Text style={styles.emptyTitle}>No accounts yet</Text>
            <Text style={styles.emptyText}>Add your first account to start tracking</Text>
          </View>
        ) : (
          accounts.map((account) => (
            <TouchableOpacity key={account.id} style={styles.accountCard}>
              <View style={styles.accountIcon}>
                <Ionicons name={accountTypeIcons[account.type.toLowerCase()] || 'wallet'} size={24} color="#8eafff" />
              </View>
              <View style={styles.accountInfo}>
                <Text style={styles.accountName}>{account.name}</Text>
                <Text style={styles.accountType}>{account.type}</Text>
              </View>
              <View style={styles.accountBalance}>
                <Text style={styles.balanceAmount}>{formatCurrency(account.totalBalance, account.currency)}</Text>
                <Text style={styles.balanceCurrency}>{account.currency}</Text>
              </View>
              <Ionicons name="chevron-forward" size={20} color="#8eafff" />
            </TouchableOpacity>
          ))
        )}
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
  emptyState: { alignItems: 'center', paddingVertical: 60 },
  emptyTitle: { fontSize: 20, fontWeight: '600', color: '#fff', marginTop: 16 },
  emptyText: { fontSize: 16, color: '#8eafff', marginTop: 8 },
  accountCard: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#161a8f22', borderRadius: 16, padding: 16, marginBottom: 12, borderWidth: 1, borderColor: '#2d4fff22' },
  accountIcon: { width: 48, height: 48, borderRadius: 12, backgroundColor: '#2d4fff22', justifyContent: 'center', alignItems: 'center', marginRight: 16 },
  accountInfo: { flex: 1 },
  accountName: { fontSize: 16, fontWeight: '600', color: '#fff', marginBottom: 4 },
  accountType: { fontSize: 14, color: '#8eafff', textTransform: 'capitalize' },
  accountBalance: { alignItems: 'flex-end', marginRight: 8 },
  balanceAmount: { fontSize: 16, fontWeight: '600', color: '#fff' },
  balanceCurrency: { fontSize: 12, color: '#8eafff' },
});

