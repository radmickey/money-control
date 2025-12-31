package handlers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/radmickey/money-control/backend/pkg/auth"
	"github.com/radmickey/money-control/backend/pkg/middleware"
	"github.com/radmickey/money-control/backend/pkg/utils"
	authpb "github.com/radmickey/money-control/backend/proto/auth"
	"github.com/radmickey/money-control/backend/services/gateway/proxy"
)

// AuthHandler handles auth-related requests
type AuthHandler struct {
	proxy        *proxy.ServiceProxy
	oauthManager *auth.OAuthManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(sp *proxy.ServiceProxy, oauth *auth.OAuthManager) *AuthHandler {
	return &AuthHandler{
		proxy:        sp,
		oauthManager: oauth,
	}
}

// generateState generates a random state for OAuth
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email        string `json:"email" binding:"required,email"`
		Password     string `json:"password" binding:"required,min=8"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		BaseCurrency string `json:"base_currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Auth.Register(c.Request.Context(), &authpb.RegisterRequest{
		Email:        req.Email,
		Password:     req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BaseCurrency: req.BaseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, gin.H{
		"user":          resp.User,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Auth.Login(c.Request.Context(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		utils.Unauthorized(c, "Invalid email or password")
		return
	}

	utils.Success(c, gin.H{
		"user":          resp.User,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Auth.RefreshToken(c.Request.Context(), &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		utils.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	utils.Success(c, gin.H{
		"user":          resp.User,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"token_type":    "Bearer",
	})
}

// GoogleAuthURL returns Google OAuth URL and redirects the user
func (h *AuthHandler) GoogleAuthURL(c *gin.Context) {
	state := generateState()

	// Store state in cookie for verification (short-lived)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	// Generate OAuth URL
	url := h.oauthManager.GetGoogleAuthURL(state)

	// Check if client wants JSON response or redirect
	if c.Query("redirect") == "false" {
		utils.Success(c, gin.H{"url": url})
		return
	}

	// Redirect to Google OAuth
	c.Redirect(302, url)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.BadRequest(c, "Authorization code is required")
		return
	}

	// Verify state (optional but recommended)
	// state := c.Query("state")
	// storedState, _ := c.Cookie("oauth_state")
	// if state != storedState {
	// 	utils.BadRequest(c, "Invalid state parameter")
	// 	return
	// }

	resp, err := h.proxy.Auth.GoogleAuth(c.Request.Context(), &authpb.GoogleAuthRequest{
		Code: code,
	})
	if err != nil {
		// Redirect to frontend with error
		c.Redirect(302, "http://localhost:3000/login?error=google_auth_failed")
		return
	}

	// Redirect to frontend with tokens (in a real app, you might use a more secure method)
	redirectURL := "http://localhost:3000/auth/callback?access_token=" + resp.AccessToken +
		"&refresh_token=" + resp.RefreshToken
	c.Redirect(302, redirectURL)
}

// GetProfile returns user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	resp, err := h.proxy.Auth.GetProfile(c.Request.Context(), &authpb.GetProfileRequest{
		UserId: userID,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		BaseCurrency string `json:"base_currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Auth.UpdateProfile(c.Request.Context(), &authpb.UpdateProfileRequest{
		UserId:       userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BaseCurrency: req.BaseCurrency,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	_, err := h.proxy.Auth.Logout(c.Request.Context(), &authpb.LogoutRequest{
		UserId:       userID,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Successfully logged out"})
}

// TelegramAuth handles Telegram WebApp authentication
func (h *AuthHandler) TelegramAuth(c *gin.Context) {
	var req struct {
		InitData string `json:"init_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.proxy.Auth.TelegramAuth(c.Request.Context(), &authpb.TelegramAuthRequest{
		InitData: req.InitData,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user":          resp.User,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
		"token_type":    "Bearer",
	})
}
