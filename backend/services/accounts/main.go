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
	"github.com/radmickey/money-control/backend/pkg/config"
	"github.com/radmickey/money-control/backend/pkg/database"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	pb "github.com/radmickey/money-control/backend/proto/accounts"
	"github.com/radmickey/money-control/backend/services/accounts/handlers"
	"github.com/radmickey/money-control/backend/services/accounts/models"
	"github.com/radmickey/money-control/backend/services/accounts/repository"
	"github.com/radmickey/money-control/backend/services/accounts/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadForService("ACCOUNTS")
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
	if err := db.Migrate(&models.Account{}, &models.SubAccount{}, &models.BalanceHistory{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db.DB)
	subAccountRepo := repository.NewSubAccountRepository(db.DB)
	balanceHistoryRepo := repository.NewBalanceHistoryRepository(db.DB)

	// Initialize service
	accountService := service.NewAccountService(accountRepo, subAccountRepo, balanceHistoryRepo)

	// Initialize JWT manager for auth middleware
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTAccessDuration,
		cfg.JWTRefreshDuration,
	)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcHandler := handlers.NewGRPCHandler(accountService)
	pb.RegisterAccountsServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Accounts gRPC server starting on port %s", cfg.GRPCPort)
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
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "accounts"})
	})

	// Protected routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(jwtManager))
	handlers.RegisterHTTPRoutes(v1, accountService)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("Accounts HTTP server starting on port %s", cfg.HTTPPort)
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

	log.Println("Accounts service stopped")
}

