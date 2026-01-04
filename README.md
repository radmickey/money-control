# Money Control

Personal finance management application with microservices architecture.

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)
![React](https://img.shields.io/badge/React-18.2-61DAFB.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

<p align="center">
  <img src="docs/screenshots/landing.png" alt="Landing Page" width="100%">
</p>

## Documentation

| Document | Description |
|----------|-------------|
| [Getting Started](docs/guides/getting-started.md) | Installation and setup |
| [Google OAuth](docs/guides/google-oauth.md) | Google authentication configuration |
| [Telegram Mini App](docs/guides/telegram-miniapp.md) | Telegram integration |
| [Resilience Patterns](docs/guides/resilience.md) | Circuit breakers, retries, health checks |
| [API Reference](docs/api/README.md) | REST API documentation |

## Features

- Multi-asset tracking: stocks, crypto, ETFs, real estate, bank accounts
- Multi-currency support with automatic conversion
- Real-time prices via Alpha Vantage and CoinGecko
- Cross-platform: Web, iOS, Android, Telegram
- JWT + OAuth + Telegram authentication
- High availability: circuit breakers, retries, health probes

## Screenshots

| Login | Register |
|-------|----------|
| ![Login](docs/screenshots/login.png) | ![Register](docs/screenshots/register.png) |

| Dashboard | Accounts |
|-----------|----------|
| ![Dashboard](docs/screenshots/dashboard.png) | ![Accounts](docs/screenshots/accounts.png) |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          Clients                                │
│     Web (React)  │  Mobile (React Native)  │  Telegram Bot      │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway (Gin)                            │
│  Rate Limiting │ JWT │ Circuit Breaker │ Health Checks          │
│  Port: 9080                                                     │
│  Endpoints: /health, /ready, /health/circuits                   │
└────────────────────────────┬────────────────────────────────────┘
                             │ gRPC (retry + timeout)
       ┌─────────────────────┼─────────────────────┐
       ▼                     ▼                     ▼
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│    Auth     │       │  Accounts   │       │Transactions │
│  :50051     │       │   :50052    │       │   :50053    │
└─────────────┘       └─────────────┘       └─────────────┘
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   Assets    │       │  Currency   │       │  Insights   │
│   :50054    │       │   :50055    │       │   :50056    │
└─────────────┘       └─────────────┘       └─────────────┘
       │                     │
       ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                      External APIs                              │
│        Alpha Vantage │ CoinGecko │ Frankfurter │ CBR            │
└─────────────────────────────────────────────────────────────────┘
```

### Resilience

| Feature | Configuration |
|---------|---------------|
| Circuit Breaker | 5 failures → open, 30s timeout |
| Retry Policy | 3 attempts, exponential backoff |
| Timeouts | 10s gRPC, 15s HTTP |
| Keepalive | 10s ping interval |

## Project Structure

```
money-control/
├── backend/
│   ├── pkg/
│   │   ├── auth/           # JWT, OAuth
│   │   ├── cache/          # Redis
│   │   ├── converters/     # Type conversions
│   │   ├── database/       # PostgreSQL
│   │   ├── health/         # Health checks
│   │   ├── middleware/     # HTTP middleware
│   │   └── resilience/     # Circuit breaker
│   ├── proto/              # Protocol Buffers
│   └── services/
│       ├── auth/
│       ├── accounts/
│       ├── transactions/
│       ├── assets/
│       ├── currency/
│       ├── insights/
│       └── gateway/
├── frontend/
│   ├── web/                # React + Vite
│   └── mobile/             # React Native
├── docs/
└── docker-compose.yml
```

## Tech Stack

| Layer | Technologies |
|-------|--------------|
| Backend | Go, Gin, gRPC, GORM, PostgreSQL, Redis |
| Frontend | React, Vite, TypeScript, Tailwind, Redux |
| Mobile | React Native, Expo |
| Infrastructure | Docker, Docker Compose |

## Quick Start

```bash
# Clone
git clone https://github.com/radmickey/money-control.git
cd money-control

# Configure
cp .env.example .env

# Run
docker compose up -d

# Access
open http://localhost:3000
```

## Health Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /health` | Liveness probe |
| `GET /ready` | Readiness probe |
| `GET /health/circuits` | Circuit breaker status |

## License

MIT License. See [LICENSE](LICENSE).
