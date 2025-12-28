package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Environment string
	Debug       bool

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// JWT
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// External APIs
	AlphaVantageAPIKey  string
	AlphaVantagePremium bool
	ExchangeRatesAPIKey string
	CoinGeckoAPIURL     string

	// Telegram
	TelegramBotToken string
	TelegramWebhook  string

	// Service Ports
	GRPCPort string
	HTTPPort string

	// Service URLs (for inter-service communication)
	AuthServiceURL         string
	AccountsServiceURL     string
	TransactionsServiceURL string
	AssetsServiceURL       string
	CurrencyServiceURL     string
	InsightsServiceURL     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		// Server
		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       getEnvBool("DEBUG", true),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// Redis
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTAccessDuration:  getEnvDuration("JWT_ACCESS_DURATION", 15*time.Minute),
		JWTRefreshDuration: getEnvDuration("JWT_REFRESH_DURATION", 7*24*time.Hour),

		// Google OAuth
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),

		// External APIs
		AlphaVantageAPIKey:  getEnv("ALPHA_VANTAGE_API_KEY", ""),
		AlphaVantagePremium: getEnvBool("ALPHA_VANTAGE_PREMIUM", false),
		ExchangeRatesAPIKey: getEnv("EXCHANGERATES_API_KEY", ""),
		CoinGeckoAPIURL:     getEnv("COINGECKO_API_URL", "https://api.coingecko.com/api/v3"),

		// Telegram
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramWebhook:  getEnv("TELEGRAM_WEBHOOK_URL", ""),

		// Service Ports
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		// Service URLs
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "localhost:50051"),
		AccountsServiceURL:     getEnv("ACCOUNTS_SERVICE_URL", "localhost:50052"),
		TransactionsServiceURL: getEnv("TRANSACTIONS_SERVICE_URL", "localhost:50053"),
		AssetsServiceURL:       getEnv("ASSETS_SERVICE_URL", "localhost:50054"),
		CurrencyServiceURL:     getEnv("CURRENCY_SERVICE_URL", "localhost:50055"),
		InsightsServiceURL:     getEnv("INSIGHTS_SERVICE_URL", "localhost:50056"),
	}

	return cfg, nil
}

// LoadForService loads service-specific configuration
func LoadForService(serviceName string) (*Config, error) {
	_ = godotenv.Load()

	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	// Override database URL with service-specific one
	dbURLKey := serviceName + "_DB_URL"
	if dbURL := getEnv(dbURLKey, ""); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	// Override ports with service-specific ones
	grpcPortKey := serviceName + "_GRPC_PORT"
	if port := getEnv(grpcPortKey, ""); port != "" {
		cfg.GRPCPort = port
	}

	httpPortKey := serviceName + "_HTTP_PORT"
	if port := getEnv(httpPortKey, ""); port != "" {
		cfg.HTTPPort = port
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

