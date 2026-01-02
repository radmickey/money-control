import React, { useEffect } from 'react';
import { View, Text, ScrollView, StyleSheet, RefreshControl, Dimensions } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';
import { PieChart } from 'react-native-chart-kit';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { fetchNetWorth, fetchAllocation } from '../store/slices/insightsSlice';
import { fetchAccounts } from '../store/slices/accountsSlice';

const screenWidth = Dimensions.get('window').width;

export default function DashboardScreen() {
  const dispatch = useAppDispatch();
  const { netWorth, allocation, loading } = useAppSelector((state) => state.insights);
  const { accounts } = useAppSelector((state) => state.accounts);
  const { user } = useAppSelector((state) => state.auth);

  const baseCurrency = user?.baseCurrency || 'USD';

  useEffect(() => {
    loadData();
  }, []);

  const loadData = () => {
    dispatch(fetchNetWorth(baseCurrency));
    dispatch(fetchAllocation(baseCurrency));
    dispatch(fetchAccounts());
  };

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: baseCurrency, minimumFractionDigits: 0 }).format(value);
  };

  const chartData = allocation.map((item, index) => ({
    name: item.category,
    population: item.percentage,
    color: ['#2d4fff', '#10b981', '#8b5cf6', '#f59e0b', '#f43f5e'][index % 5],
    legendFontColor: '#8eafff',
    legendFontSize: 12,
  }));

  return (
    <SafeAreaView style={styles.container} edges={['left', 'right']}>
      <ScrollView style={styles.scrollView} refreshControl={<RefreshControl refreshing={loading} onRefresh={loadData} tintColor="#2d4fff" />}>
        <View style={styles.header}>
          <Text style={styles.greeting}>Hello, {user?.name?.split(' ')[0] || 'there'} 👋</Text>
          <Text style={styles.subtitle}>Your financial overview</Text>
        </View>

        <View style={styles.netWorthCard}>
          <Text style={styles.cardLabel}>Total Net Worth</Text>
          <Text style={styles.netWorthValue}>{netWorth ? formatCurrency(netWorth.total) : '$0'}</Text>
          {netWorth && (
            <View style={[styles.changeContainer, { backgroundColor: netWorth.change24h >= 0 ? '#10b98122' : '#f43f5e22' }]}>
              <Ionicons name={netWorth.change24h >= 0 ? 'trending-up' : 'trending-down'} size={16} color={netWorth.change24h >= 0 ? '#10b981' : '#f43f5e'} />
              <Text style={[styles.changeText, { color: netWorth.change24h >= 0 ? '#10b981' : '#f43f5e' }]}>
                {netWorth.changePercent24h >= 0 ? '+' : ''}{netWorth.changePercent24h.toFixed(2)}% today
              </Text>
            </View>
          )}
        </View>

        {chartData.length > 0 && (
          <View style={styles.chartCard}>
            <Text style={styles.chartTitle}>Asset Allocation</Text>
            <PieChart
              data={chartData}
              width={screenWidth - 48}
              height={200}
              chartConfig={{ color: () => '#fff' }}
              accessor="population"
              backgroundColor="transparent"
              paddingLeft="0"
              absolute
            />
          </View>
        )}

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Accounts</Text>
          {accounts.slice(0, 4).map((account) => (
            <View key={account.id} style={styles.accountItem}>
              <View style={styles.accountIcon}>
                <Ionicons name="wallet" size={20} color="#8eafff" />
              </View>
              <View style={styles.accountInfo}>
                <Text style={styles.accountName}>{account.name}</Text>
                <Text style={styles.accountType}>{account.type}</Text>
              </View>
              <Text style={styles.accountBalance}>{formatCurrency(account.totalBalance)}</Text>
            </View>
          ))}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#0c0d1f' },
  scrollView: { flex: 1, padding: 24 },
  header: { marginBottom: 24 },
  greeting: { fontSize: 28, fontWeight: 'bold', color: '#fff', marginBottom: 4 },
  subtitle: { fontSize: 16, color: '#8eafff' },
  netWorthCard: { backgroundColor: '#161a8f22', borderRadius: 20, padding: 24, marginBottom: 24, borderWidth: 1, borderColor: '#2d4fff33' },
  cardLabel: { fontSize: 14, color: '#8eafff', marginBottom: 8 },
  netWorthValue: { fontSize: 36, fontWeight: 'bold', color: '#fff', marginBottom: 12 },
  changeContainer: { flexDirection: 'row', alignItems: 'center', alignSelf: 'flex-start', paddingHorizontal: 12, paddingVertical: 6, borderRadius: 20, gap: 4 },
  changeText: { fontSize: 14, fontWeight: '500' },
  chartCard: { backgroundColor: '#161a8f22', borderRadius: 20, padding: 20, marginBottom: 24, borderWidth: 1, borderColor: '#2d4fff33' },
  chartTitle: { fontSize: 18, fontWeight: '600', color: '#fff', marginBottom: 16 },
  section: { marginBottom: 24 },
  sectionTitle: { fontSize: 18, fontWeight: '600', color: '#fff', marginBottom: 16 },
  accountItem: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#161a8f22', borderRadius: 16, padding: 16, marginBottom: 12, borderWidth: 1, borderColor: '#2d4fff22' },
  accountIcon: { width: 44, height: 44, borderRadius: 12, backgroundColor: '#2d4fff22', justifyContent: 'center', alignItems: 'center', marginRight: 12 },
  accountInfo: { flex: 1 },
  accountName: { fontSize: 16, fontWeight: '500', color: '#fff', marginBottom: 2 },
  accountType: { fontSize: 14, color: '#8eafff' },
  accountBalance: { fontSize: 16, fontWeight: '600', color: '#fff' },
});

