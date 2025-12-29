package handlers

import (
	"context"

	pb "github.com/radmickey/money-control/backend/proto/auth"
	"github.com/radmickey/money-control/backend/services/auth/models"
	"github.com/radmickey/money-control/backend/services/auth/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the AuthServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(authService *service.AuthService) *GRPCHandler {
	return &GRPCHandler{
		authService: authService,
	}
}

// Register registers a new user
func (h *GRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	result, err := h.authService.Register(ctx, service.RegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BaseCurrency: req.BaseCurrency,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         userToProto(result.User),
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// Login logs in a user
func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	result, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         userToProto(result.User),
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// RefreshToken refreshes the access token
func (h *GRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	result, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         userToProto(result.User),
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// GetProfile gets the user profile
func (h *GRPCHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.User, error) {
	user, err := h.authService.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return userToProto(user), nil
}

// UpdateProfile updates the user profile
func (h *GRPCHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.User, error) {
	user, err := h.authService.UpdateProfile(ctx, req.UserId, req.FirstName, req.LastName, req.BaseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return userToProto(user), nil
}

// GoogleAuth handles Google OAuth authentication
func (h *GRPCHandler) GoogleAuth(ctx context.Context, req *pb.GoogleAuthRequest) (*pb.AuthResponse, error) {
	result, err := h.authService.GoogleAuth(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "google auth failed: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         userToProto(result.User),
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// TelegramAuth handles Telegram WebApp authentication
func (h *GRPCHandler) TelegramAuth(ctx context.Context, req *pb.TelegramAuthRequest) (*pb.AuthResponse, error) {
	// Parse initData
	telegramData, err := ParseTelegramInitData(req.InitData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid telegram init data: %v", err)
	}

	result, err := h.authService.TelegramAuth(ctx, service.TelegramAuthInput{
		ID:        telegramData.User.ID,
		FirstName: telegramData.User.FirstName,
		LastName:  telegramData.User.LastName,
		Username:  telegramData.User.Username,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "telegram auth failed: %v", err)
	}

	return &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         userToProto(result.User),
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// ValidateToken validates an access token
func (h *GRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := h.authService.ValidateToken(req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

// Logout logs out a user
func (h *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authService.Logout(ctx, req.UserId, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout: %v", err)
	}

	return &pb.LogoutResponse{Success: true}, nil
}

// Helper functions
func userToProto(u *models.User) *pb.User {
	googleID := ""
	if u.GoogleID != nil {
		googleID = *u.GoogleID
	}
	var telegramID int64 = 0
	if u.TelegramID != nil {
		telegramID = *u.TelegramID
	}
	telegramUsername := ""
	if u.TelegramUsername != nil {
		telegramUsername = *u.TelegramUsername
	}
	return &pb.User{
		Id:               u.ID,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		BaseCurrency:     u.BaseCurrency,
		GoogleId:         googleID,
		TelegramId:       telegramID,
		TelegramUsername: telegramUsername,
		CreatedAt:        timestamppb.New(u.CreatedAt),
		UpdatedAt:        timestamppb.New(u.UpdatedAt),
	}
}
