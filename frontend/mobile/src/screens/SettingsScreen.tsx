import React from 'react';
import { View, Text, ScrollView, StyleSheet, TouchableOpacity, Switch } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Ionicons } from '@expo/vector-icons';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { logout } from '../store/slices/authSlice';

export default function SettingsScreen() {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const [notifications, setNotifications] = React.useState(true);
  const [biometric, setBiometric] = React.useState(false);

  const handleLogout = () => { dispatch(logout()); };

  const SettingItem = ({ icon, title, subtitle, onPress, right }: any) => (
    <TouchableOpacity style={styles.settingItem} onPress={onPress}>
      <View style={styles.settingIcon}>
        <Ionicons name={icon} size={22} color="#8eafff" />
      </View>
      <View style={styles.settingInfo}>
        <Text style={styles.settingTitle}>{title}</Text>
        {subtitle && <Text style={styles.settingSubtitle}>{subtitle}</Text>}
      </View>
      {right || <Ionicons name="chevron-forward" size={20} color="#8eafff" />}
    </TouchableOpacity>
  );

  return (
    <SafeAreaView style={styles.container} edges={['left', 'right']}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.profileSection}>
          <View style={styles.avatar}>
            <Text style={styles.avatarText}>{user?.name?.[0]?.toUpperCase() || 'U'}</Text>
          </View>
          <Text style={styles.profileName}>{user?.name || 'User'}</Text>
          <Text style={styles.profileEmail}>{user?.email}</Text>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Account</Text>
          <SettingItem icon="person-outline" title="Profile" subtitle="Edit your profile" />
          <SettingItem icon="globe-outline" title="Base Currency" subtitle={user?.baseCurrency || 'USD'} />
          <SettingItem icon="card-outline" title="Subscriptions" subtitle="Manage your plan" />
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Preferences</Text>
          <SettingItem icon="notifications-outline" title="Notifications" right={<Switch value={notifications} onValueChange={setNotifications} trackColor={{ true: '#2d4fff' }} />} />
          <SettingItem icon="finger-print" title="Biometric Login" right={<Switch value={biometric} onValueChange={setBiometric} trackColor={{ true: '#2d4fff' }} />} />
          <SettingItem icon="moon-outline" title="Dark Mode" subtitle="Always on" />
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Security</Text>
          <SettingItem icon="lock-closed-outline" title="Change Password" />
          <SettingItem icon="shield-checkmark-outline" title="Two-Factor Auth" subtitle="Not enabled" />
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Support</Text>
          <SettingItem icon="help-circle-outline" title="Help Center" />
          <SettingItem icon="chatbubble-outline" title="Contact Us" />
          <SettingItem icon="document-text-outline" title="Privacy Policy" />
        </View>

        <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
          <Ionicons name="log-out-outline" size={22} color="#f43f5e" />
          <Text style={styles.logoutText}>Log Out</Text>
        </TouchableOpacity>

        <Text style={styles.version}>Version 1.0.0</Text>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#0c0d1f' },
  scrollView: { flex: 1, padding: 24 },
  profileSection: { alignItems: 'center', marginBottom: 32 },
  avatar: { width: 80, height: 80, borderRadius: 40, backgroundColor: '#2d4fff', justifyContent: 'center', alignItems: 'center', marginBottom: 16 },
  avatarText: { fontSize: 32, fontWeight: 'bold', color: '#fff' },
  profileName: { fontSize: 24, fontWeight: 'bold', color: '#fff', marginBottom: 4 },
  profileEmail: { fontSize: 16, color: '#8eafff' },
  section: { marginBottom: 24 },
  sectionTitle: { fontSize: 14, fontWeight: '600', color: '#8eafff', textTransform: 'uppercase', marginBottom: 12, letterSpacing: 0.5 },
  settingItem: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#161a8f22', borderRadius: 12, padding: 16, marginBottom: 8, borderWidth: 1, borderColor: '#2d4fff22' },
  settingIcon: { width: 40, height: 40, borderRadius: 10, backgroundColor: '#2d4fff22', justifyContent: 'center', alignItems: 'center', marginRight: 12 },
  settingInfo: { flex: 1 },
  settingTitle: { fontSize: 16, fontWeight: '500', color: '#fff' },
  settingSubtitle: { fontSize: 14, color: '#8eafff', marginTop: 2 },
  logoutButton: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', backgroundColor: '#f43f5e22', borderRadius: 12, padding: 16, marginTop: 8, gap: 8 },
  logoutText: { fontSize: 16, fontWeight: '600', color: '#f43f5e' },
  version: { textAlign: 'center', color: '#8eafff88', marginTop: 24, marginBottom: 40 },
});

