# Google OAuth Setup

This guide explains how to set up Google OAuth for Money Control.

## Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click **"Select a project"** → **"New Project"**
3. Enter a project name (e.g., "Money Control")
4. Click **"Create"**

## Step 2: Enable OAuth Consent Screen

1. Navigate to **APIs & Services** → **OAuth consent screen**
2. Select **"External"** user type
3. Fill in the required fields:
   - **App name**: Money Control
   - **User support email**: your email
   - **Developer contact**: your email
4. Click **"Save and Continue"**
5. Skip scopes for now, click **"Save and Continue"**
6. Add test users if needed, click **"Save and Continue"**

## Step 3: Create OAuth Credentials

1. Navigate to **APIs & Services** → **Credentials**
2. Click **"+ Create Credentials"** → **"OAuth client ID"**
3. Select **"Web application"**
4. Name it (e.g., "Money Control Web")
5. Add **Authorized redirect URIs**:
   - For local development: `http://localhost:9080/api/v1/auth/google/callback`
   - For production: `https://yourdomain.com/api/v1/auth/google/callback`
6. Click **"Create"**

## Step 4: Copy Credentials

After creation, you'll see:
- **Client ID**: `xxxxxx.apps.googleusercontent.com`
- **Client Secret**: `GOCSPX-xxxxx`

## Step 5: Configure Environment

Add to your `.env` file:

```bash
GOOGLE_CLIENT_ID=your_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:9080/api/v1/auth/google/callback
```

## Step 6: Restart Services

```bash
export $(cat .env | grep -v '^#' | xargs)
docker compose up -d --force-recreate gateway auth-service
```

## Testing

1. Go to http://localhost:3000/login
2. Click **"Continue with Google"**
3. You should be redirected to Google's login page
4. After authentication, you'll be redirected back to the app

## Troubleshooting

### "redirect_uri_mismatch" Error

Make sure the redirect URI in Google Cloud Console exactly matches:
```
http://localhost:9080/api/v1/auth/google/callback
```

### "Access blocked" Error

Your app might be in testing mode. Add your Google account as a test user in the OAuth consent screen.

### Empty client_id in URL

Environment variables not loaded. Run:
```bash
export $(cat .env | grep -v '^#' | xargs)
```

