import React from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  TrendingUp,
  Shield,
  Zap,
  PieChart,
  Wallet,
  BarChart3,
  ArrowRight,
  Sparkles,
  Globe,
  LineChart,
  Building2,
  Coins,
} from 'lucide-react';

const Landing: React.FC = () => {
  const features = [
    {
      icon: Wallet,
      title: 'All Assets in One Place',
      description: 'Track stocks, crypto, real estate, bank accounts, and cash — unified view of your entire portfolio.',
      gradient: 'from-violet-500 to-purple-600',
    },
    {
      icon: LineChart,
      title: 'Real-time Updates',
      description: 'Live price feeds from Alpha Vantage and CoinGecko. Always know your current net worth.',
      gradient: 'from-emerald-500 to-teal-600',
    },
    {
      icon: Globe,
      title: 'Multi-Currency Support',
      description: 'Track assets in any currency. Automatic conversion to your base currency for unified insights.',
      gradient: 'from-blue-500 to-cyan-600',
    },
    {
      icon: PieChart,
      title: 'Smart Allocation',
      description: 'Visual breakdown of your portfolio. Understand your asset distribution at a glance.',
      gradient: 'from-amber-500 to-orange-600',
    },
    {
      icon: Shield,
      title: 'Bank-Level Security',
      description: 'End-to-end encryption, secure authentication, and your data never leaves your control.',
      gradient: 'from-rose-500 to-pink-600',
    },
    {
      icon: Zap,
      title: 'Lightning Fast',
      description: 'Built with Go microservices for blazing performance. Sub-second response times.',
      gradient: 'from-indigo-500 to-violet-600',
    },
  ];

  const assetTypes = [
    { icon: Building2, label: 'Banks', color: 'text-blue-400' },
    { icon: Coins, label: 'Crypto', color: 'text-amber-400' },
    { icon: BarChart3, label: 'Stocks', color: 'text-emerald-400' },
    { icon: TrendingUp, label: 'ETFs', color: 'text-purple-400' },
  ];

  return (
    <div className="min-h-screen overflow-hidden">
      {/* Animated background */}
      <div className="fixed inset-0 -z-10">
        <div className="absolute inset-0 bg-[#0a0b14]" />
        <div className="absolute top-0 left-1/4 w-[600px] h-[600px] bg-violet-600/20 rounded-full blur-[120px] animate-pulse" />
        <div className="absolute bottom-0 right-1/4 w-[500px] h-[500px] bg-blue-600/15 rounded-full blur-[100px] animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-indigo-600/10 rounded-full blur-[150px]" />

        {/* Grid pattern */}
        <div
          className="absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage: `linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)`,
            backgroundSize: '60px 60px',
          }}
        />
      </div>

      {/* Navigation */}
      <nav className="relative z-10 flex items-center justify-between px-8 py-6 max-w-7xl mx-auto">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="flex items-center gap-3"
        >
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-violet-500/30">
            <Sparkles className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-bold tracking-tight">Money Control</span>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          className="flex items-center gap-4"
        >
          <Link
            to="/login"
            className="text-zinc-400 hover:text-white transition-colors px-4 py-2"
          >
            Sign In
          </Link>
          <Link
            to="/register"
            className="bg-white text-zinc-900 px-5 py-2.5 rounded-full font-medium hover:bg-zinc-100 transition-colors"
          >
            Get Started
          </Link>
        </motion.div>
      </nav>

      {/* Hero Section */}
      <section className="relative z-10 px-8 pt-16 pb-20 max-w-7xl mx-auto">
        <div className="text-center max-w-4xl mx-auto">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-violet-500/10 border border-violet-500/20 text-violet-300 text-sm mb-8"
          >
            <Sparkles className="w-4 h-4" />
            <span>Track your entire net worth in one place</span>
          </motion.div>

          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-5xl md:text-7xl font-bold tracking-tight leading-[1.1] mb-6"
          >
            Your wealth,{' '}
            <span className="bg-gradient-to-r from-violet-400 via-fuchsia-400 to-indigo-400 bg-clip-text text-transparent">
              visualized
            </span>
          </motion.h1>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="text-xl text-zinc-400 max-w-2xl mx-auto mb-10"
          >
            Aggregate all your assets — crypto, stocks, real estate, bank accounts —
            into a beautiful dashboard. Know your net worth, anytime.
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.4 }}
            className="flex flex-col sm:flex-row items-center justify-center gap-4"
          >
            <Link
              to="/register"
              className="group flex items-center gap-2 bg-gradient-to-r from-violet-600 to-indigo-600 text-white px-8 py-4 rounded-full font-semibold text-lg hover:shadow-xl hover:shadow-violet-500/25 transition-all duration-300 hover:-translate-y-0.5"
            >
              Start for Free
              <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </Link>
            <a
              href="#features"
              className="text-zinc-400 hover:text-white transition-colors px-6 py-4"
            >
              Learn more ↓
            </a>
          </motion.div>

          {/* Asset types showcase */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5 }}
            className="mt-16 flex items-center justify-center gap-8 flex-wrap"
          >
            {assetTypes.map((asset, i) => (
              <motion.div
                key={asset.label}
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.6 + i * 0.1 }}
                className="flex items-center gap-2 text-zinc-500"
              >
                <asset.icon className={`w-5 h-5 ${asset.color}`} />
                <span>{asset.label}</span>
              </motion.div>
            ))}
          </motion.div>
        </div>

        {/* Hero Dashboard Preview */}
        <motion.div
          initial={{ opacity: 0, y: 60 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7, duration: 0.8 }}
          className="mt-20 relative"
        >
          <div className="absolute inset-0 bg-gradient-to-t from-[#0a0b14] via-transparent to-transparent z-10 pointer-events-none" />

          <div className="relative mx-auto max-w-5xl rounded-2xl overflow-hidden border border-zinc-800/50 shadow-2xl shadow-violet-500/10">
            {/* Mock Dashboard */}
            <div className="bg-gradient-to-br from-zinc-900 via-zinc-900 to-zinc-950 p-6">
              {/* Top bar */}
              <div className="flex items-center gap-2 mb-6">
                <div className="w-3 h-3 rounded-full bg-rose-500" />
                <div className="w-3 h-3 rounded-full bg-amber-500" />
                <div className="w-3 h-3 rounded-full bg-emerald-500" />
              </div>

              {/* Dashboard content */}
              <div className="grid grid-cols-3 gap-4">
                {/* Net Worth Card */}
                <div className="col-span-2 bg-gradient-to-br from-violet-600/20 to-indigo-600/10 rounded-xl p-6 border border-violet-500/20">
                  <p className="text-zinc-500 text-sm">Total Net Worth</p>
                  <p className="text-4xl font-bold mt-2">$847,293</p>
                  <div className="flex items-center gap-2 mt-2 text-emerald-400 text-sm">
                    <TrendingUp className="w-4 h-4" />
                    <span>+12.4% this month</span>
                  </div>
                </div>

                {/* Mini cards */}
                <div className="space-y-4">
                  <div className="bg-zinc-800/50 rounded-xl p-4 border border-zinc-700/50">
                    <p className="text-zinc-500 text-xs">Stocks</p>
                    <p className="text-lg font-semibold">$324,100</p>
                  </div>
                  <div className="bg-zinc-800/50 rounded-xl p-4 border border-zinc-700/50">
                    <p className="text-zinc-500 text-xs">Crypto</p>
                    <p className="text-lg font-semibold">$156,820</p>
                  </div>
                </div>
              </div>

              {/* Chart placeholder */}
              <div className="mt-6 h-32 bg-zinc-800/30 rounded-xl border border-zinc-700/30 flex items-end justify-around px-4 pb-4">
                {[40, 65, 45, 70, 55, 80, 60, 75, 85, 70, 90, 85].map((h, i) => (
                  <div
                    key={i}
                    className="w-6 rounded-t bg-gradient-to-t from-violet-600 to-violet-400"
                    style={{ height: `${h}%` }}
                  />
                ))}
              </div>
            </div>
          </div>
        </motion.div>
      </section>

      {/* Features Section */}
      <section id="features" className="relative z-10 px-8 py-24 max-w-7xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8 }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl font-bold mb-4">Everything you need</h2>
          <p className="text-xl text-zinc-400">Powerful features to take control of your finances</p>
        </motion.div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, i) => (
            <motion.div
              key={feature.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.9 + i * 0.1 }}
              className="group relative bg-zinc-900/50 backdrop-blur-sm rounded-2xl p-6 border border-zinc-800/50 hover:border-zinc-700/50 transition-all duration-300 hover:-translate-y-1"
            >
              <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${feature.gradient} flex items-center justify-center mb-4 shadow-lg`}>
                <feature.icon className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
              <p className="text-zinc-400">{feature.description}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* CTA Section */}
      <section className="relative z-10 px-8 py-24">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 1.5 }}
          className="max-w-4xl mx-auto text-center bg-gradient-to-br from-violet-600/20 via-indigo-600/10 to-fuchsia-600/10 rounded-3xl p-12 border border-violet-500/20"
        >
          <h2 className="text-4xl font-bold mb-4">Ready to take control?</h2>
          <p className="text-xl text-zinc-400 mb-8">
            Join thousands of users tracking their net worth with Money Control
          </p>
          <Link
            to="/register"
            className="inline-flex items-center gap-2 bg-white text-zinc-900 px-8 py-4 rounded-full font-semibold text-lg hover:bg-zinc-100 transition-colors"
          >
            Get Started Free
            <ArrowRight className="w-5 h-5" />
          </Link>
        </motion.div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 px-8 py-12 border-t border-zinc-800/50">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-violet-400" />
            <span className="font-semibold">Money Control</span>
          </div>
          <p className="text-zinc-500 text-sm">
            © 2026 Money Control. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Landing;

