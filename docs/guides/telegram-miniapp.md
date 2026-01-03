# Telegram Mini App Setup

This guide explains how to set up Money Control as a Telegram Mini App.

## Prerequisites

- Telegram account
- Public HTTPS URL (for production) or ngrok (for development)

## Step 1: Create a Bot

1. Open Telegram and search for [@BotFather](https://t.me/BotFather)
2. Send `/newbot` command
3. Follow the prompts:
   - Enter bot name (e.g., "Money Control")
   - Enter bot username (e.g., "MoneyControlBot")
4. Save the **bot token** you receive

## Step 2: Configure Environment

Add to your `.env` file:

```bash
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

## Step 3: Set Up Web App URL

### For Development (using ngrok)

1. Install ngrok:
   ```bash
   brew install ngrok  # macOS
   # or download from https://ngrok.com/download
   ```

2. Start your frontend:
   ```bash
   cd frontend/web
   npm run dev
   ```

3. Start ngrok tunnel:
   ```bash
   ngrok http 3000
   ```

4. Copy the HTTPS URL (e.g., `https://abc123.ngrok.io`)

### Configure in BotFather

1. Open [@BotFather](https://t.me/BotFather)
2. Send `/mybots`
3. Select your bot
4. Go to **Bot Settings** → **Menu Button** → **Configure Menu Button**
5. Enter your URL:
   - Development: `https://abc123.ngrok.io`
   - Production: `https://yourdomain.com`

## Step 4: Test the Mini App

1. Open your bot in Telegram
2. Click the menu button (hamburger icon at bottom left)
3. The Mini App should load

## How It Works

### Authentication Flow

1. User opens Mini App in Telegram
2. Telegram provides `initData` with user info and hash
3. Frontend sends `initData` to `/api/v1/auth/telegram`
4. Backend validates hash using bot token
5. User is created/logged in automatically

### initData Validation

The backend validates Telegram initData using HMAC-SHA256:

```go
// 1. Create secret key
secretKey = HMAC-SHA256("WebAppData", botToken)

// 2. Create data-check-string
dataCheckString = "auth_date=...\nquery_id=...\nuser=..."

// 3. Calculate hash
calculatedHash = HMAC-SHA256(secretKey, dataCheckString)

// 4. Compare with provided hash
valid = calculatedHash == providedHash
```

## Frontend Integration

The Mini App SDK is already integrated:

```typescript
// hooks/useTelegram.ts
import { useTelegram } from '../hooks/useTelegram';

const { webApp } = useTelegram();

// Access Telegram data
if (webApp) {
  console.log(webApp.initData);        // Raw init data
  console.log(webApp.initDataUnsafe);  // Parsed init data
  console.log(webApp.colorScheme);     // 'light' or 'dark'
}
```

## Styling for Telegram

The Mini App automatically adapts to Telegram's theme:

```typescript
if (webApp) {
  // Use Telegram theme colors
  webApp.setHeaderColor('#0a0b14');
  webApp.setBackgroundColor('#0a0b14');

  // Expand to full screen
  webApp.expand();
}
```

## Troubleshooting

### "Invalid hash" Error

- Check that `TELEGRAM_BOT_TOKEN` is correct
- Ensure the initData is being sent correctly

### Mini App Not Loading

- Verify ngrok is running
- Check that the URL is HTTPS
- Try clearing Telegram cache

### "User data not found" Error

- The initData must include user information
- This happens when the Mini App is opened without proper context

## Production Deployment

For production:

1. Deploy frontend to a server with HTTPS
2. Update BotFather with production URL
3. Set environment variables on your server
4. Remove ngrok dependency

