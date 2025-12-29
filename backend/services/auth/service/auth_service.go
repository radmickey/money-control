package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/radmickey/money-control/backend/pkg/auth"
	"github.com/radmickey/money-control/backend/services/auth/models"
	"github.com/radmickey/money-control/backend/services/auth/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActive      = errors.New("user account is not active")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	oauthStateRepo   *repository.OAuthStateRepository
	jwtManager       *auth.JWTManager
	oauthManager     *auth.OAuthManager
	refreshDuration  time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	oauthStateRepo *repository.OAuthStateRepository,
	jwtManager *auth.JWTManager,
	oauthManager *auth.OAuthManager,
	refreshDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		oauthStateRepo:   oauthStateRepo,
		jwtManager:       jwtManager,
		oauthManager:     oauthManager,
		refreshDuration:  refreshDuration,
	}
}

// RegisterInput holds registration input data
type RegisterInput struct {
	Email        string
	Password     string
	FirstName    string
	LastName     string
	BaseCurrency string
}

// AuthResult holds authentication result
type AuthResult struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	// Validate password
	if err := auth.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Set default currency if not provided
	if input.BaseCurrency == "" {
		input.BaseCurrency = "USD"
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		BaseCurrency: input.BaseCurrency,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateAuthResult(ctx, user)
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	// Find user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if err := auth.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	return s.generateAuthResult(ctx, user)
}

// GoogleAuth authenticates or registers a user via Google OAuth
func (s *AuthService) GoogleAuth(ctx context.Context, code string) (*AuthResult, error) {
	// Get user info from Google
	googleUser, err := s.oauthManager.AuthenticateWithGoogle(ctx, code)
	if err != nil {
		return nil, err
	}

	// Try to find existing user by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleUser.ID)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	if user == nil {
		// Try to find by email
		user, err = s.userRepo.GetByEmail(ctx, googleUser.Email)
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}

		if user != nil {
			// Link Google account to existing user
			user.GoogleID = &googleUser.ID
			if user.AvatarURL == "" {
				user.AvatarURL = googleUser.Picture
			}
			if err := s.userRepo.Update(ctx, user); err != nil {
				return nil, err
			}
		} else {
			// Create new user
			user = &models.User{
				Email:        googleUser.Email,
				GoogleID:     &googleUser.ID,
				FirstName:    googleUser.GivenName,
				LastName:     googleUser.FamilyName,
				AvatarURL:    googleUser.Picture,
				BaseCurrency: "USD",
				IsActive:     true,
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return s.generateAuthResult(ctx, user)
}

// TelegramAuthInput holds Telegram auth data
type TelegramAuthInput struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
}

// TelegramAuth authenticates or registers a user via Telegram
func (s *AuthService) TelegramAuth(ctx context.Context, input TelegramAuthInput) (*AuthResult, error) {
	// Try to find existing user by Telegram ID
	user, err := s.userRepo.GetByTelegramID(ctx, input.ID)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	if user == nil {
		// Create new user - generate unique email for Telegram users
		email := uuid.New().String() + "@telegram.user"
		username := input.Username
		user = &models.User{
			Email:            email,
			TelegramID:       &input.ID,
			TelegramUsername: &username,
			FirstName:        input.FirstName,
			LastName:         input.LastName,
			BaseCurrency:     "USD",
			IsActive:         true,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return s.generateAuthResult(ctx, user)
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *AuthService) GetGoogleAuthURL(ctx context.Context) (string, error) {
	// Generate and store state
	state := uuid.New().String()
	oauthState := &models.OAuthState{
		State:     state,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.oauthStateRepo.Create(ctx, oauthState); err != nil {
		return "", err
	}

	return s.oauthManager.GetGoogleAuthURL(state), nil
}

// ValidateOAuthState validates an OAuth state
func (s *AuthService) ValidateOAuthState(ctx context.Context, state string) error {
	return s.oauthStateRepo.Validate(ctx, state)
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	// Validate refresh token in database
	rt, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Revoke old refresh token
	_ = s.refreshTokenRepo.Revoke(ctx, refreshToken)

	// Generate new tokens
	return s.generateAuthResult(ctx, user)
}

// Logout revokes the refresh token
func (s *AuthService) Logout(ctx context.Context, userID, refreshToken string) error {
	return s.refreshTokenRepo.Revoke(ctx, refreshToken)
}

// LogoutAll revokes all refresh tokens for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
}

// GetProfile returns user profile
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID, firstName, lastName, baseCurrency string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if baseCurrency != "" {
		user.BaseCurrency = baseCurrency
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateToken validates an access token
func (s *AuthService) ValidateToken(token string) (*auth.AccessClaims, error) {
	return s.jwtManager.ValidateAccessToken(token)
}

func (s *AuthService) generateAuthResult(ctx context.Context, user *models.User) (*AuthResult, error) {
	// Generate token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(s.refreshDuration),
	}

	if err := s.refreshTokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

