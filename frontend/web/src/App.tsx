import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAppSelector } from './store/hooks';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Accounts from './pages/Accounts';
import Transactions from './pages/Transactions';
import Assets from './pages/Assets';
import Settings from './pages/Settings';
import Login from './pages/Login';
import Register from './pages/Register';
import Landing from './pages/Landing';
import AuthCallback from './pages/AuthCallback';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, loading } = useAppSelector((state) => state.auth);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#0a0b14]">
        <div className="relative">
          <div className="w-16 h-16 border-4 border-violet-500/20 rounded-full" />
          <div className="absolute inset-0 w-16 h-16 border-4 border-transparent border-t-violet-500 rounded-full animate-spin" />
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/welcome" replace />;
  }

  return <>{children}</>;
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAppSelector((state) => state.auth);

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/welcome" element={<PublicRoute><Landing /></PublicRoute>} />
      <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
      <Route path="/register" element={<PublicRoute><Register /></PublicRoute>} />
      <Route path="/auth/callback" element={<AuthCallback />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="accounts" element={<Accounts />} />
        <Route path="transactions" element={<Transactions />} />
        <Route path="assets" element={<Assets />} />
        <Route path="settings" element={<Settings />} />
      </Route>
      <Route path="*" element={<Navigate to="/welcome" replace />} />
    </Routes>
  );
};

export default App;
