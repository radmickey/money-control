import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  User,
  Bell,
  Shield,
  CreditCard,
  Globe,
  Palette,
  Smartphone,
  Key,
  Save,
} from 'lucide-react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { updateProfile, getProfile } from '../store/slices/authSlice';

const currencies = [
  { code: 'USD', name: 'US Dollar', symbol: '$' },
  { code: 'EUR', name: 'Euro', symbol: '€' },
  { code: 'GBP', name: 'British Pound', symbol: '£' },
  { code: 'JPY', name: 'Japanese Yen', symbol: '¥' },
  { code: 'RUB', name: 'Russian Ruble', symbol: '₽' },
  { code: 'CHF', name: 'Swiss Franc', symbol: 'CHF' },
  { code: 'CAD', name: 'Canadian Dollar', symbol: 'C$' },
  { code: 'AUD', name: 'Australian Dollar', symbol: 'A$' },
];

const Settings: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user, loading } = useAppSelector((state) => state.auth);
  const [activeSection, setActiveSection] = useState('profile');
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    baseCurrency: 'USD',
    notifications: {
      email: true,
      push: true,
      priceAlerts: true,
      weeklyReport: true,
    },
    security: {
      twoFactor: false,
      biometric: false,
    },
  });
  const [saved, setSaved] = useState(false);

  // Load user profile on mount (only if not already loaded)
  useEffect(() => {
    if (!user) {
      dispatch(getProfile());
    }
  }, [dispatch, user]);

  // Sync formData with user when user changes
  useEffect(() => {
    if (user) {
      setFormData(prev => ({
        ...prev,
        name: user.firstName || '',
        email: user.email || '',
        baseCurrency: user.baseCurrency || 'USD',
      }));
    }
  }, [user]);

  const sections = [
    { id: 'profile', icon: User, label: 'Profile' },
    { id: 'currency', icon: Globe, label: 'Currency' },
    { id: 'notifications', icon: Bell, label: 'Notifications' },
    { id: 'security', icon: Shield, label: 'Security' },
    { id: 'integrations', icon: CreditCard, label: 'Integrations' },
  ];

  const handleSave = async () => {
    try {
      await dispatch(updateProfile({
        firstName: formData.name,
        baseCurrency: formData.baseCurrency,
      })).unwrap();
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (error) {
      console.error('Failed to save settings:', error);
    }
  };

  const renderSection = () => {
    switch (activeSection) {
      case 'profile':
        return (
          <div className="space-y-6">
            <div>
              <label className="block text-sm text-midnight-400 mb-2">Display Name</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-sm text-midnight-400 mb-2">Email</label>
              <input
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-sm text-midnight-400 mb-2">Profile Picture</label>
              <div className="flex items-center gap-4">
                <div className="w-16 h-16 rounded-full bg-gradient-to-br from-accent-emerald to-accent-violet flex items-center justify-center">
                  <User className="w-8 h-8 text-white" />
                </div>
                <button className="px-4 py-2 rounded-lg border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors">
                  Change
                </button>
              </div>
            </div>
          </div>
        );

      case 'currency':
        return (
          <div className="space-y-6">
            <div>
              <label className="block text-sm text-midnight-400 mb-2">Base Currency</label>
              <p className="text-sm text-midnight-500 mb-4">
                All your assets and net worth will be converted to this currency
              </p>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                {currencies.map((currency) => (
                  <button
                    key={currency.code}
                    onClick={() => setFormData({ ...formData, baseCurrency: currency.code })}
                    className={`p-4 rounded-xl border transition-all ${
                      formData.baseCurrency === currency.code
                        ? 'border-midnight-500 bg-midnight-500/20'
                        : 'border-midnight-800 hover:border-midnight-600'
                    }`}
                  >
                    <p className="text-2xl mb-1">{currency.symbol}</p>
                    <p className="text-sm font-medium">{currency.code}</p>
                    <p className="text-xs text-midnight-400">{currency.name}</p>
                  </button>
                ))}
              </div>
            </div>
          </div>
        );

      case 'notifications':
        return (
          <div className="space-y-6">
            {[
              { key: 'email', label: 'Email Notifications', desc: 'Receive updates via email' },
              { key: 'push', label: 'Push Notifications', desc: 'Receive mobile push notifications' },
              { key: 'priceAlerts', label: 'Price Alerts', desc: 'Get notified on significant price changes' },
              { key: 'weeklyReport', label: 'Weekly Report', desc: 'Receive a weekly summary of your finances' },
            ].map((item) => (
              <div
                key={item.key}
                className="flex items-center justify-between p-4 rounded-xl bg-midnight-950/50"
              >
                <div>
                  <p className="font-medium">{item.label}</p>
                  <p className="text-sm text-midnight-400">{item.desc}</p>
                </div>
                <button
                  onClick={() =>
                    setFormData({
                      ...formData,
                      notifications: {
                        ...formData.notifications,
                        [item.key]: !formData.notifications[item.key as keyof typeof formData.notifications],
                      },
                    })
                  }
                  className={`w-12 h-6 rounded-full transition-colors relative ${
                    formData.notifications[item.key as keyof typeof formData.notifications]
                      ? 'bg-accent-emerald'
                      : 'bg-midnight-700'
                  }`}
                >
                  <span
                    className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-all ${
                      formData.notifications[item.key as keyof typeof formData.notifications]
                        ? 'right-1'
                        : 'left-1'
                    }`}
                  />
                </button>
              </div>
            ))}
          </div>
        );

      case 'security':
        return (
          <div className="space-y-6">
            <div className="p-4 rounded-xl bg-midnight-950/50">
              <div className="flex items-center gap-3 mb-4">
                <Key className="w-5 h-5 text-midnight-400" />
                <div>
                  <p className="font-medium">Change Password</p>
                  <p className="text-sm text-midnight-400">Update your password regularly</p>
                </div>
              </div>
              <button className="px-4 py-2 rounded-lg border border-midnight-700 text-midnight-300 hover:bg-midnight-800/50 transition-colors">
                Change Password
              </button>
            </div>

            <div className="flex items-center justify-between p-4 rounded-xl bg-midnight-950/50">
              <div className="flex items-center gap-3">
                <Shield className="w-5 h-5 text-midnight-400" />
                <div>
                  <p className="font-medium">Two-Factor Authentication</p>
                  <p className="text-sm text-midnight-400">Add an extra layer of security</p>
                </div>
              </div>
              <button
                onClick={() =>
                  setFormData({
                    ...formData,
                    security: { ...formData.security, twoFactor: !formData.security.twoFactor },
                  })
                }
                className={`w-12 h-6 rounded-full transition-colors relative ${
                  formData.security.twoFactor ? 'bg-accent-emerald' : 'bg-midnight-700'
                }`}
              >
                <span
                  className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-all ${
                    formData.security.twoFactor ? 'right-1' : 'left-1'
                  }`}
                />
              </button>
            </div>

            <div className="flex items-center justify-between p-4 rounded-xl bg-midnight-950/50">
              <div className="flex items-center gap-3">
                <Smartphone className="w-5 h-5 text-midnight-400" />
                <div>
                  <p className="font-medium">Biometric Authentication</p>
                  <p className="text-sm text-midnight-400">Use Face ID or fingerprint</p>
                </div>
              </div>
              <button
                onClick={() =>
                  setFormData({
                    ...formData,
                    security: { ...formData.security, biometric: !formData.security.biometric },
                  })
                }
                className={`w-12 h-6 rounded-full transition-colors relative ${
                  formData.security.biometric ? 'bg-accent-emerald' : 'bg-midnight-700'
                }`}
              >
                <span
                  className={`absolute top-1 w-4 h-4 rounded-full bg-white transition-all ${
                    formData.security.biometric ? 'right-1' : 'left-1'
                  }`}
                />
              </button>
            </div>
          </div>
        );

      case 'integrations':
        return (
          <div className="space-y-6">
            <p className="text-midnight-400">
              Connect your accounts to automatically sync transactions and balances.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[
                { name: 'Bank Connection', desc: 'Connect via Plaid', connected: false },
                { name: 'Google Account', desc: 'Sign in with Google', connected: true },
                { name: 'Telegram', desc: 'Link your Telegram account', connected: false },
                { name: 'Crypto Wallets', desc: 'Connect MetaMask, etc.', connected: false },
              ].map((integration) => (
                <div
                  key={integration.name}
                  className="p-4 rounded-xl border border-midnight-800 hover:border-midnight-600 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">{integration.name}</p>
                      <p className="text-sm text-midnight-400">{integration.desc}</p>
                    </div>
                    <button
                      className={`px-3 py-1 rounded-lg text-sm ${
                        integration.connected
                          ? 'bg-accent-emerald/20 text-accent-emerald'
                          : 'bg-midnight-700 text-midnight-300 hover:bg-midnight-600'
                      }`}
                    >
                      {integration.connected ? 'Connected' : 'Connect'}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-display font-bold text-gradient">Settings</h1>
        <p className="text-midnight-400 mt-1">Manage your account preferences</p>
      </div>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Sidebar */}
        <div className="lg:w-64 flex-shrink-0">
          <nav className="glass rounded-2xl p-2 space-y-1">
            {sections.map((section) => (
              <button
                key={section.id}
                onClick={() => setActiveSection(section.id)}
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                  activeSection === section.id
                    ? 'bg-midnight-500/20 text-white'
                    : 'text-midnight-400 hover:text-white hover:bg-midnight-800/30'
                }`}
              >
                <section.icon className="w-5 h-5" />
                <span>{section.label}</span>
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1">
          <motion.div
            key={activeSection}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="glass rounded-2xl p-6"
          >
            <h2 className="text-xl font-semibold mb-6 capitalize">{activeSection}</h2>
            {renderSection()}

            <div className="mt-8 pt-6 border-t border-midnight-800/50 flex justify-end">
              <button
                onClick={handleSave}
                className={`btn-primary flex items-center gap-2 ${
                  saved ? 'bg-accent-emerald' : ''
                }`}
              >
                <Save className="w-4 h-4" />
                {saved ? 'Saved!' : 'Save Changes'}
              </button>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
};

export default Settings;

