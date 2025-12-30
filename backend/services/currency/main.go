package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
	pkgcache "github.com/radmickey/money-control/backend/pkg/cache"
	"github.com/radmickey/money-control/backend/pkg/config"
	"github.com/radmickey/money-control/backend/pkg/database"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	pb "github.com/radmickey/money-control/backend/proto/currency"
	"github.com/radmickey/money-control/backend/services/currency/handlers"
	"github.com/radmickey/money-control/backend/services/currency/models"
	"github.com/radmickey/money-control/backend/services/currency/providers"
	"github.com/radmickey/money-control/backend/services/currency/repository"
	"github.com/radmickey/money-control/backend/services/currency/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadForService("CURRENCY")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(database.Config{
		URL:   cfg.DatabaseURL,
		Debug: cfg.Debug,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(&models.ExchangeRate{}, &models.RateHistory{}, &models.Currency{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to Redis
	redisCache, err := pkgcache.New(pkgcache.Config{
		URL:    cfg.RedisURL,
		Prefix: "currency",
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Initialize exchange rates client
	ratesClient := providers.NewExchangeRatesClient(cfg.ExchangeRatesAPIKey)

	// Initialize repositories
	currencyRepo := repository.NewCurrencyRepository(db.DB)
	rateRepo := repository.NewExchangeRateRepository(db.DB)
	historyRepo := repository.NewRateHistoryRepository(db.DB)

	// Seed currencies
	if err := currencyRepo.SeedCurrencies(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed currencies: %v", err)
	}

	// Initialize service
	currencyService := service.NewCurrencyService(
		currencyRepo, rateRepo, historyRepo,
		ratesClient, redisCache, "USD",
	)

	// Start rate updater (update every hour)
	currencyService.StartRateUpdater(1 * time.Hour)
	defer currencyService.Stop()

	// Initialize JWT manager for auth middleware
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTAccessDuration,
		cfg.JWTRefreshDuration,
	)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcHandler := handlers.NewGRPCHandler(currencyService)
	pb.RegisterCurrencyServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Currency gRPC server starting on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "currency"})
	})

	// Routes (some public, some protected)
	v1 := router.Group("/api/v1")
	handlers.RegisterHTTPRoutes(v1, currencyService, jwtManager)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("Currency HTTP server starting on port %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Currency service stopped")
}

