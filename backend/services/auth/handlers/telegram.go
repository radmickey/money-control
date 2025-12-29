package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

// TelegramUser represents user data from Telegram
type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsPremium    bool   `json:"is_premium"`
}

// TelegramInitData represents parsed Telegram Web App init data
type TelegramInitData struct {
	QueryID      string        `json:"query_id"`
	User         *TelegramUser `json:"user"`
	AuthDate     int64         `json:"auth_date"`
	Hash         string        `json:"hash"`
	ChatType     string        `json:"chat_type"`
	ChatInstance string        `json:"chat_instance"`
}

// ParseTelegramInitData parses and validates Telegram Web App init data
func ParseTelegramInitData(initData string) (*TelegramInitData, error) {
	// Parse the URL-encoded data
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init data: %w", err)
	}

	// Extract hash
	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash not found in init data")
	}

	// Validate the hash
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken != "" {
		if !validateTelegramHash(values, hash, botToken) {
			return nil, fmt.Errorf("invalid hash")
		}
	}

	// Parse the data
	data := &TelegramInitData{
		QueryID:      values.Get("query_id"),
		Hash:         hash,
		ChatType:     values.Get("chat_type"),
		ChatInstance: values.Get("chat_instance"),
	}

	// Parse auth_date
	if authDateStr := values.Get("auth_date"); authDateStr != "" {
		authDate, err := strconv.ParseInt(authDateStr, 10, 64)
		if err == nil {
			data.AuthDate = authDate
		}
	}

	// Parse user
	if userStr := values.Get("user"); userStr != "" {
		var user TelegramUser
		if err := json.Unmarshal([]byte(userStr), &user); err != nil {
			return nil, fmt.Errorf("failed to parse user data: %w", err)
		}
		data.User = &user
	}

	if data.User == nil {
		return nil, fmt.Errorf("user data not found in init data")
	}

	return data, nil
}

// validateTelegramHash validates the hash of Telegram init data
func validateTelegramHash(values url.Values, hash, botToken string) bool {
	// Create data-check-string
	var keys []string
	for key := range values {
		if key != "hash" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var parts []string
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	dataCheckString := strings.Join(parts, "\n")

	// Create secret key using HMAC-SHA256 of bot token with "WebAppData"
	secretKeyHMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHMAC.Write([]byte(botToken))
	secretKey := secretKeyHMAC.Sum(nil)

	// Create hash of data-check-string
	dataHMAC := hmac.New(sha256.New, secretKey)
	dataHMAC.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(dataHMAC.Sum(nil))

	return calculatedHash == hash
}

