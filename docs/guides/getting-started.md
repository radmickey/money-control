# Getting Started

This guide will help you set up and run Money Control locally.

## Prerequisites

- **Docker** and **Docker Compose** (recommended)
- Or for local development:
  - Go 1.24+
  - Node.js 20+
  - PostgreSQL 16+
  - Redis 7+

## Quick Start with Docker

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/money-control.git
cd money-control
```

### 2. Create Environment File

```bash
cat > .env << 'EOF'
# Database
DB_PASSWORD=postgres

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Google OAuth (optional)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:9080/api/v1/auth/google/callback

# Telegram (optional)
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
EOF
```

### 3. Start Services

```bash
# Export environment variables
export $(cat .env | grep -v '^#' | xargs)

# Start all services
docker compose up -d
```

### 4. Access the Application

| Service | URL |
|---------|-----|
| Web App | http://localhost:3000 |
| API Gateway | http://localhost:9080 |
| Health Check | http://localhost:9080/health |

### 5. Create Your First Account

1. Open http://localhost:3000
2. Click **"Get Started"** to register
3. Fill in your details and click **"Create Account"**
4. You'll be redirected to the dashboard

## Local Development

### Backend

```bash
# Install dependencies
cd backend
go mod download

# Start databases
docker compose up -d auth-db accounts-db transactions-db assets-db currency-db insights-db redis

# Run auth service
go run ./services/auth

# Run gateway (in another terminal)
go run ./services/gateway
```

### Frontend

```bash
cd frontend/web
npm install
npm run dev
```

## Next Steps

- [Configure Google OAuth](./google-oauth.md)
- [Set up Telegram Mini App](./telegram-miniapp.md)
- [API Documentation](../api/README.md)

