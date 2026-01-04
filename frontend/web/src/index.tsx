import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { store } from './store';
import App from './App';
import TelegramProvider from './components/TelegramProvider';
import AuthInitializer from './components/AuthInitializer';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <Provider store={store}>
      <AuthInitializer>
        <TelegramProvider>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </TelegramProvider>
      </AuthInitializer>
    </Provider>
  </React.StrictMode>
);

