import { useEffect, useRef } from 'react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { getProfile } from '../store/slices/authSlice';

/**
 * AuthInitializer loads user profile on app start if token exists.
 * This ensures baseCurrency and user data are available immediately.
 */
const AuthInitializer: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useAppDispatch();
  const { token, user, loading } = useAppSelector((state) => state.auth);
  const initialized = useRef(false);

  useEffect(() => {
    // Load profile only once if token exists but user is not loaded
    if (token && !user && !loading && !initialized.current) {
      initialized.current = true;
      dispatch(getProfile());
    }
  }, [dispatch, token, user, loading]);

  return <>{children}</>;
};

export default AuthInitializer;

