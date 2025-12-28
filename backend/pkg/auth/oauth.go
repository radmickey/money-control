package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo holds user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// OAuthManager handles OAuth operations
type OAuthManager struct {
	googleConfig *oauth2.Config
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(googleClientID, googleClientSecret, redirectURL string) *OAuthManager {
	return &OAuthManager{
		googleConfig: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (m *OAuthManager) GetGoogleAuthURL(state string) string {
	return m.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeGoogleCode exchanges an authorization code for tokens
func (m *OAuthManager) ExchangeGoogleCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return m.googleConfig.Exchange(ctx, code)
}

// GetGoogleUserInfo fetches user info from Google using the access token
func (m *OAuthManager) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := m.googleConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// AuthenticateWithGoogle performs the full OAuth flow: exchange code and get user info
func (m *OAuthManager) AuthenticateWithGoogle(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := m.ExchangeGoogleCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return m.GetGoogleUserInfo(ctx, token)
}

