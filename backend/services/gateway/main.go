package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
	"github.com/radmickey/money-control/backend/pkg/cache"
	"github.com/radmickey/money-control/backend/pkg/config"
	"github.com/radmickey/money-control/backend/pkg/health"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/resilience"
	"github.com/radmickey/money-control/backend/services/gateway/handlers"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

const version = "1.0.0"

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize health checker
	healthChecker := health.NewHealthChecker(version)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTAccessDuration,
		cfg.JWTRefreshDuration,
	)

	// Initialize OAuth manager
	oauthManager := auth.NewOAuthManager(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	// Initialize gRPC service proxies (with retry policy and keepalive)
	serviceProxy, err := proxy.NewServiceProxy(proxy.Config{
		AuthServiceURL:         cfg.AuthServiceURL,
		AccountsServiceURL:     cfg.AccountsServiceURL,
		TransactionsServiceURL: cfg.TransactionsServiceURL,
		AssetsServiceURL:       cfg.AssetsServiceURL,
		CurrencyServiceURL:     cfg.CurrencyServiceURL,
		InsightsServiceURL:     cfg.InsightsServiceURL,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize some service proxies: %v", err)
		// Continue anyway for development
	}
	if serviceProxy != nil {
		defer serviceProxy.Close()
	}

	// Try to connect to Redis for caching and rate limiting (optional)
	var redisCache *cache.Cache
	var rateLimiter *middleware.RateLimiter
	if cfg.RedisURL != "" {
		redisCache, err = cache.New(cache.Config{
			URL:    cfg.RedisURL,
			Prefix: "gateway",
		})
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis (rate limiting disabled): %v", err)
		} else {
			defer redisCache.Close()
			rateLimiter = middleware.NewRateLimiter(redisCache.Client(), "gateway")

			// Register Redis health check
			healthChecker.Register("redis", health.RedisCheck(func(ctx context.Context) error {
				return redisCache.Client().Ping(ctx).Err()
			}))
		}
	}

	// Setup Gin router
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORS())

	// Only enable rate limiting if Redis is available
	if rateLimiter != nil {
		router.Use(middleware.RateLimitMiddleware(rateLimiter, middleware.RateLimitConfig{
			Requests: 100,
			Window:   time.Minute,
		}))
	}

	// Health endpoints
	registerHealthEndpoints(router, healthChecker, redisCache)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Register route handlers
	handlers.RegisterRoutes(v1, serviceProxy, jwtManager, oauthManager)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API Gateway v%s starting on port %s", version, cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Gateway stopped")
}

// registerHealthEndpoints registers health check endpoints
func registerHealthEndpoints(router *gin.Engine, healthChecker *health.HealthChecker, redisCache *cache.Cache) {
	// Basic health check (liveness probe)
	router.GET("/health", func(c *gin.Context) {
		check := healthChecker.Liveness(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{
			"status":    check.Status,
			"service":   "gateway",
			"version":   version,
			"timestamp": check.Timestamp,
		})
	})

	// Readiness probe - checks all dependencies
	router.GET("/ready", func(c *gin.Context) {
		report := healthChecker.Readiness(c.Request.Context())

		statusCode := http.StatusOK
		if report.Status != health.StatusUp {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, report)
	})

	// Detailed health check
	router.GET("/health/details", func(c *gin.Context) {
		report := healthChecker.Check(c.Request.Context())

		// Add circuit breaker stats
		cbStats := resilience.GlobalManager.AllStats()

		c.JSON(http.StatusOK, gin.H{
			"health":           report,
			"circuit_breakers": cbStats,
		})
	})

	// Circuit breaker status endpoint
	router.GET("/health/circuits", func(c *gin.Context) {
		stats := resilience.GlobalManager.AllStats()
		c.JSON(http.StatusOK, stats)
	})
}
