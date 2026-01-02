import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAppDispatch } from '../store/hooks';
import { setCredentials, getProfile } from '../store/slices/authSlice';
import { Sparkles } from 'lucide-react';

const AuthCallback: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const error = searchParams.get('error');

    if (error) {
      navigate('/login?error=' + error);
      return;
    }

    if (accessToken && refreshToken) {
      // Save tokens and fetch profile
      dispatch(setCredentials({ accessToken, refreshToken }));
      dispatch(getProfile()).then(() => {
        navigate('/');
      });
    } else {
      navigate('/login?error=missing_tokens');
    }
  }, [searchParams, dispatch, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0a0b14]">
      <div className="text-center">
        <div className="w-16 h-16 mx-auto mb-6 rounded-xl bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-violet-500/25">
          <Sparkles className="w-8 h-8 text-white" />
        </div>
        <div className="relative mb-4">
          <div className="w-12 h-12 mx-auto border-4 border-violet-500/20 rounded-full" />
          <div className="absolute inset-0 w-12 h-12 mx-auto border-4 border-transparent border-t-violet-500 rounded-full animate-spin" />
        </div>
        <p className="text-zinc-400">Completing authentication...</p>
      </div>
    </div>
  );
};

export default AuthCallback;

