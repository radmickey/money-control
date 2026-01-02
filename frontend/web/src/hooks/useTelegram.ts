// Telegram Web App SDK hook
// https://core.telegram.org/bots/webapps

declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp;
    };
  }
}

interface TelegramWebApp {
  initData: string;
  initDataUnsafe: {
    user?: {
      id: number;
      first_name: string;
      last_name?: string;
      username?: string;
      language_code?: string;
      is_premium?: boolean;
    };
    auth_date: number;
    hash: string;
    query_id?: string;
  };
  version: string;
  platform: string;
  colorScheme: 'light' | 'dark';
  themeParams: {
    bg_color?: string;
    text_color?: string;
    hint_color?: string;
    link_color?: string;
    button_color?: string;
    button_text_color?: string;
    secondary_bg_color?: string;
  };
  isExpanded: boolean;
  viewportHeight: number;
  viewportStableHeight: number;
  headerColor: string;
  backgroundColor: string;
  isClosingConfirmationEnabled: boolean;
  BackButton: {
    isVisible: boolean;
    onClick: (callback: () => void) => void;
    offClick: (callback: () => void) => void;
    show: () => void;
    hide: () => void;
  };
  MainButton: {
    text: string;
    color: string;
    textColor: string;
    isVisible: boolean;
    isProgressVisible: boolean;
    isActive: boolean;
    setText: (text: string) => void;
    onClick: (callback: () => void) => void;
    offClick: (callback: () => void) => void;
    show: () => void;
    hide: () => void;
    enable: () => void;
    disable: () => void;
    showProgress: (leaveActive?: boolean) => void;
    hideProgress: () => void;
  };
  HapticFeedback: {
    impactOccurred: (style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft') => void;
    notificationOccurred: (type: 'error' | 'success' | 'warning') => void;
    selectionChanged: () => void;
  };
  ready: () => void;
  expand: () => void;
  close: () => void;
  sendData: (data: string) => void;
  openLink: (url: string, options?: { try_instant_view?: boolean }) => void;
  openTelegramLink: (url: string) => void;
  showPopup: (params: {
    title?: string;
    message: string;
    buttons?: Array<{
      id?: string;
      type?: 'default' | 'ok' | 'close' | 'cancel' | 'destructive';
      text?: string;
    }>;
  }, callback?: (buttonId: string) => void) => void;
  showAlert: (message: string, callback?: () => void) => void;
  showConfirm: (message: string, callback?: (confirmed: boolean) => void) => void;
  enableClosingConfirmation: () => void;
  disableClosingConfirmation: () => void;
  setHeaderColor: (color: 'bg_color' | 'secondary_bg_color' | string) => void;
  setBackgroundColor: (color: 'bg_color' | 'secondary_bg_color' | string) => void;
}

export const useTelegram = () => {
  const tg = window.Telegram?.WebApp;
  const isTelegram = !!tg;
  const user = tg?.initDataUnsafe?.user;
  const initData = tg?.initData;

  const onClose = () => {
    tg?.close();
  };

  const onToggleButton = () => {
    if (tg?.MainButton.isVisible) {
      tg.MainButton.hide();
    } else {
      tg?.MainButton.show();
    }
  };

  const hapticFeedback = (type: 'success' | 'error' | 'warning' | 'light' | 'medium' | 'heavy') => {
    if (!tg?.HapticFeedback) return;

    if (['success', 'error', 'warning'].includes(type)) {
      tg.HapticFeedback.notificationOccurred(type as 'success' | 'error' | 'warning');
    } else {
      tg.HapticFeedback.impactOccurred(type as 'light' | 'medium' | 'heavy');
    }
  };

  return {
    tg,
    isTelegram,
    user,
    initData,
    onClose,
    onToggleButton,
    hapticFeedback,
    queryId: tg?.initDataUnsafe?.query_id,
  };
};

export default useTelegram;

