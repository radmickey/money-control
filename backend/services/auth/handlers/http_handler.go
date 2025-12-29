package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	"github.com/radmickey/money-control/backend/services/auth/models"
	"github.com/radmickey/money-control/backend/services/auth/repository"
	"github.com/radmickey/money-control/backend/services/auth/service"
)

// HTTPHandler handles HTTP requests for auth
type HTTPHandler struct {
	authService *service.AuthService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(authService *service.AuthService) *HTTPHandler {
	return &HTTPHandler{
		authService: authService,
	}
}

// RegisterRoutes registers HTTP routes
func (h *HTTPHandler) RegisterRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.GET("/google", h.GoogleAuthURL)
		auth.GET("/google/callback", h.GoogleCallback)
	}

	// Protected routes
	protected := r.Group("/auth")
	protected.Use(authMiddleware)
	{
		protected.GET("/profile", h.GetProfile)
		protected.PUT("/profile", h.UpdateProfile)
		protected.POST("/logout", h.Logout)
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BaseCurrency string `json:"base_currency"`
}

// Register handles user registration
func (h *HTTPHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	input := service.RegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BaseCurrency: req.BaseCurrency,
	}

	result, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		if err == repository.ErrUserExists {
			utils.Conflict(c, "User with this email already exists")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, gin.H{
		"user":          userToResponse(result.User),
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login
func (h *HTTPHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			utils.Unauthorized(c, "Invalid email or password")
			return
		}
		if err == service.ErrUserNotActive {
			utils.Forbidden(c, "Account is not active")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user":          userToResponse(result.User),
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken handles token refresh
func (h *HTTPHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == repository.ErrInvalidToken {
			utils.Unauthorized(c, "Invalid or expired refresh token")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user":          userToResponse(result.User),
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// GoogleAuthURL returns the Google OAuth URL
func (h *HTTPHandler) GoogleAuthURL(c *gin.Context) {
	url, err := h.authService.GetGoogleAuthURL(c.Request.Context())
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"url": url,
	})
}

// GoogleCallback handles Google OAuth callback
func (h *HTTPHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		utils.BadRequest(c, "Authorization code is required")
		return
	}

	// Validate state
	if state != "" {
		if err := h.authService.ValidateOAuthState(c.Request.Context(), state); err != nil {
			utils.BadRequest(c, "Invalid OAuth state")
			return
		}
	}

	result, err := h.authService.GoogleAuth(c.Request.Context(), code)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	// In production, you might redirect to frontend with tokens
	utils.Success(c, gin.H{
		"user":          userToResponse(result.User),
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// GetProfile returns the current user's profile
func (h *HTTPHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, userToResponse(user))
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BaseCurrency string `json:"base_currency"`
}

// UpdateProfile updates the current user's profile
func (h *HTTPHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, req.FirstName, req.LastName, req.BaseCurrency)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, userToResponse(user))
}

// Logout handles user logout
func (h *HTTPHandler) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Successfully logged out"})
}

// UserResponse represents user response
type UserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BaseCurrency string `json:"base_currency"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	CreatedAt    string `json:"created_at"`
}

func userToResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		BaseCurrency: u.BaseCurrency,
		AvatarURL:    u.AvatarURL,
		CreatedAt:    u.CreatedAt.Format(http.TimeFormat),
	}
}

