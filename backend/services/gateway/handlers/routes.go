package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.RouterGroup, sp *proxy.ServiceProxy, jwtManager *auth.JWTManager, oauthManager *auth.OAuthManager) {
	authHandler := NewAuthHandler(sp, oauthManager)
	accountsHandler := NewAccountsHandler(sp)
	transactionsHandler := NewTransactionsHandler(sp)
	assetsHandler := NewAssetsHandler(sp)
	currencyHandler := NewCurrencyHandler(sp)
	insightsHandler := NewInsightsHandler(sp)

	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Auth routes (public)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.GET("/google", authHandler.GoogleAuthURL)
		authRoutes.GET("/google/callback", authHandler.GoogleCallback)
		authRoutes.POST("/telegram", authHandler.TelegramAuth)
	}

	// Auth routes (protected)
	authProtected := r.Group("/auth")
	authProtected.Use(authMiddleware)
	{
		authProtected.GET("/profile", authHandler.GetProfile)
		authProtected.PUT("/profile", authHandler.UpdateProfile)
		authProtected.POST("/logout", authHandler.Logout)
	}

	// Accounts routes (protected)
	accountsRoutes := r.Group("/accounts")
	accountsRoutes.Use(authMiddleware)
	{
		accountsRoutes.POST("", accountsHandler.CreateAccount)
		accountsRoutes.GET("", accountsHandler.ListAccounts)
		accountsRoutes.GET("/:id", accountsHandler.GetAccount)
		accountsRoutes.PUT("/:id", accountsHandler.UpdateAccount)
		accountsRoutes.DELETE("/:id", accountsHandler.DeleteAccount)
		accountsRoutes.POST("/:id/sub-accounts", accountsHandler.CreateSubAccount)
		accountsRoutes.GET("/:id/sub-accounts", accountsHandler.ListSubAccounts)
	}

	// Sub-accounts routes (protected)
	subAccountsRoutes := r.Group("/sub-accounts")
	subAccountsRoutes.Use(authMiddleware)
	{
		subAccountsRoutes.GET("/:id", accountsHandler.GetSubAccount)
		subAccountsRoutes.PUT("/:id", accountsHandler.UpdateSubAccount)
		subAccountsRoutes.DELETE("/:id", accountsHandler.DeleteSubAccount)
		subAccountsRoutes.PATCH("/:id/balance", accountsHandler.UpdateSubAccountBalance)
	}

	// Transactions routes (protected)
	transactionsRoutes := r.Group("/transactions")
	transactionsRoutes.Use(authMiddleware)
	{
		transactionsRoutes.POST("", transactionsHandler.CreateTransaction)
		transactionsRoutes.GET("", transactionsHandler.ListTransactions)
		transactionsRoutes.GET("/summary", transactionsHandler.GetSummary)
		transactionsRoutes.GET("/:id", transactionsHandler.GetTransaction)
		transactionsRoutes.PUT("/:id", transactionsHandler.UpdateTransaction)
		transactionsRoutes.DELETE("/:id", transactionsHandler.DeleteTransaction)
	}

	// Assets routes (protected)
	assetsRoutes := r.Group("/assets")
	assetsRoutes.Use(authMiddleware)
	{
		assetsRoutes.POST("", assetsHandler.CreateAsset)
		assetsRoutes.GET("", assetsHandler.ListAssets)
		assetsRoutes.GET("/:id", assetsHandler.GetAsset)
		assetsRoutes.PUT("/:id", assetsHandler.UpdateAsset)
		assetsRoutes.DELETE("/:id", assetsHandler.DeleteAsset)
	}

	// Prices routes (public)
	pricesRoutes := r.Group("/prices")
	{
		pricesRoutes.GET("/:symbol", assetsHandler.GetPrice)
	}

	// Search route (protected)
	r.GET("/search", authMiddleware, assetsHandler.SearchAssets)

	// Currency routes (public)
	currencyRoutes := r.Group("/currencies")
	{
		currencyRoutes.GET("", currencyHandler.ListCurrencies)
		currencyRoutes.GET("/rates/:base", currencyHandler.GetRates)
		currencyRoutes.GET("/rate/:from/:to", currencyHandler.GetRate)
		currencyRoutes.POST("/convert", currencyHandler.Convert)
	}

	// Insights routes (protected)
	insightsRoutes := r.Group("/insights")
	insightsRoutes.Use(authMiddleware)
	{
		insightsRoutes.GET("/net-worth", insightsHandler.GetNetWorth)
		insightsRoutes.GET("/net-worth/history", insightsHandler.GetNetWorthHistory)
		insightsRoutes.GET("/trends", insightsHandler.GetTrends)
		insightsRoutes.GET("/allocation", insightsHandler.GetAllocation)
		insightsRoutes.GET("/dashboard", insightsHandler.GetDashboard)
		insightsRoutes.GET("/cash-flow", insightsHandler.GetCashFlow)
	}

	// Net worth convenience route (protected)
	r.GET("/net-worth", authMiddleware, accountsHandler.GetNetWorth)
}
