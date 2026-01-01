import React, { useEffect, createContext, useContext, useState } from 'react';
import { useTelegram } from '../hooks/useTelegram';

interface TelegramContextType {
  isTelegram: boolean;
  user: {
    id: number;
    first_name: string;
    last_name?: string;
    username?: string;
  } | undefined;
  initData: string | undefined;
  isReady: boolean;
}

const TelegramContext = createContext<TelegramContextType>({
  isTelegram: false,
  user: undefined,
  initData: undefined,
  isReady: false,
});

export const useTelegramContext = () => useContext(TelegramContext);

interface Props {
  children: React.ReactNode;
}

export const TelegramProvider: React.FC<Props> = ({ children }) => {
  const { tg, isTelegram, user, initData } = useTelegram();
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (tg) {
      // Initialize Telegram Web App
      tg.ready();
      tg.expand();

      // Set theme colors to match our app
      tg.setHeaderColor('#0a0b14');
      tg.setBackgroundColor('#0a0b14');

      // Enable closing confirmation for important actions
      tg.enableClosingConfirmation();

      setIsReady(true);

      console.log('Telegram WebApp initialized:', {
        version: tg.version,
        platform: tg.platform,
        colorScheme: tg.colorScheme,
        user: user,
      });
    } else {
      setIsReady(true); // Not in Telegram, proceed normally
    }
  }, [tg, user]);

  return (
    <TelegramContext.Provider value={{ isTelegram, user, initData, isReady }}>
      {children}
    </TelegramContext.Provider>
  );
};

export default TelegramProvider;

