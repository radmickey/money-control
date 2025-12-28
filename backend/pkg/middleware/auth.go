package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// UserIDKey is the context key for user ID
	UserIDKey = "userID"
	// UserEmailKey is the context key for user email
	UserEmailKey = "userEmail"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if token is present, but doesn't require it
func OptionalAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	if !strings.HasPrefix(header, BearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimPrefix(header, BearerPrefix)
	if token == "" {
		return "", errors.New("token is required")
	}

	return token, nil
}

// GetUserID extracts user ID from Gin context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail extracts user email from Gin context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// MustGetUserID extracts user ID from Gin context or panics
func MustGetUserID(c *gin.Context) string {
	userID, ok := GetUserID(c)
	if !ok {
		panic("user ID not found in context")
	}
	return userID
}

// GRPCAuthKey is the context key for gRPC auth
type GRPCAuthKey struct{}

// GRPCUserContext holds user info for gRPC context
type GRPCUserContext struct {
	UserID string
	Email  string
}

// NewGRPCContext creates a context with user info for gRPC calls
func NewGRPCContext(ctx context.Context, userID, email string) context.Context {
	return context.WithValue(ctx, GRPCAuthKey{}, &GRPCUserContext{
		UserID: userID,
		Email:  email,
	})
}

// GetGRPCUser extracts user info from gRPC context
func GetGRPCUser(ctx context.Context) (*GRPCUserContext, bool) {
	user, ok := ctx.Value(GRPCAuthKey{}).(*GRPCUserContext)
	return user, ok
}
