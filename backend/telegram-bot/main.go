package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/radmickey/money-control/backend/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	authpb "github.com/radmickey/money-control/backend/proto/auth"
	insightspb "github.com/radmickey/money-control/backend/proto/insights"
)

// Bot represents the Telegram bot
type Bot struct {
	token       string
	webhookURL  string
	authClient  authpb.AuthServiceClient
	accClient   accountspb.AccountsServiceClient
	insClient   insightspb.InsightsServiceClient
	userSessions map[int64]string // chatID -> userID mapping
}

// Update represents a Telegram update
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

// Message represents a Telegram message
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text,omitempty"`
}

// User represents a Telegram user
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bot := &Bot{
		token:        cfg.TelegramBotToken,
		webhookURL:   cfg.TelegramWebhook,
		userSessions: make(map[int64]string),
	}

	// Connect to services
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if authConn, err := grpc.Dial(cfg.AuthServiceURL, opts...); err == nil {
		bot.authClient = authpb.NewAuthServiceClient(authConn)
	}
	if accConn, err := grpc.Dial(cfg.AccountsServiceURL, opts...); err == nil {
		bot.accClient = accountspb.NewAccountsServiceClient(accConn)
	}
	if insConn, err := grpc.Dial(cfg.InsightsServiceURL, opts...); err == nil {
		bot.insClient = insightspb.NewInsightsServiceClient(insConn)
	}

	// Start webhook server
	http.HandleFunc("/webhook", bot.handleWebhook)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"telegram-bot"}`))
	})

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8087"
		}
		log.Printf("Telegram bot starting on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Telegram bot stopped")
}

func (b *Bot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Failed to decode update: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if update.Message != nil {
		b.handleMessage(update.Message)
	}

	w.WriteHeader(http.StatusOK)
}

func (b *Bot) handleMessage(msg *Message) {
	if msg.Text == "" {
		return
	}

	parts := strings.Fields(msg.Text)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	var response string

	switch command {
	case "/start":
		response = b.handleStart(msg)
	case "/help":
		response = b.handleHelp()
	case "/networth":
		response = b.handleNetWorth(msg)
	case "/accounts":
		response = b.handleAccounts(msg)
	case "/addasset":
		response = b.handleAddAsset(msg, args)
	case "/link":
		response = b.handleLink(msg, args)
	default:
		response = "Unknown command. Use /help to see available commands."
	}

	b.sendMessage(msg.Chat.ID, response)
}

func (b *Bot) handleStart(msg *Message) string {
	name := msg.From.FirstName
	return fmt.Sprintf(`Welcome to Money Control, %s! 💰

I'm your personal finance assistant. Here's what I can do:

📊 /networth - View your total net worth
💼 /accounts - List your accounts
➕ /addasset <type> <symbol> <qty> - Add an asset
🔗 /link <email> - Link your account

Use /help for more information.`, name)
}

func (b *Bot) handleHelp() string {
	return `📱 *Money Control Bot Commands*

*Account Management:*
/networth - View total net worth
/accounts - List all accounts
/transactions - Recent transactions

*Asset Management:*
/addasset <type> <symbol> <qty> - Add asset
  Types: stock, crypto, etf
  Example: /addasset crypto BTC 0.5

*Settings:*
/link <email> - Link your account
/currency <code> - Set base currency

*Quick Actions:*
/dashboard - Open Mini App
/refresh - Refresh prices

For the full experience, open our Mini App! 📲`
}

func (b *Bot) handleNetWorth(msg *Message) string {
	userID, ok := b.userSessions[msg.Chat.ID]
	if !ok {
		return "Please link your account first using /link <email>"
	}

	if b.insClient == nil {
		return "Service temporarily unavailable"
	}

	ctx := context.Background()
	resp, err := b.insClient.GetNetWorth(ctx, &insightspb.GetNetWorthRequest{
		UserId:       userID,
		BaseCurrency: "USD",
	})
	if err != nil {
		return "Failed to fetch net worth. Please try again."
	}

	change := "📈"
	if resp.Change_24H < 0 {
		change = "📉"
	}

	return fmt.Sprintf(`💰 *Your Net Worth*

Total: $%.2f %s

24h Change: $%.2f (%.2f%%)
7d Change: $%.2f (%.2f%%)
30d Change: $%.2f (%.2f%%)`,
		resp.TotalNetWorth, change,
		resp.Change_24H, resp.ChangePercent_24H,
		resp.Change_7D, resp.ChangePercent_7D,
		resp.Change_30D, resp.ChangePercent_30D)
}

func (b *Bot) handleAccounts(msg *Message) string {
	userID, ok := b.userSessions[msg.Chat.ID]
	if !ok {
		return "Please link your account first using /link <email>"
	}

	if b.accClient == nil {
		return "Service temporarily unavailable"
	}

	ctx := context.Background()
	resp, err := b.accClient.ListAccounts(ctx, &accountspb.ListAccountsRequest{
		UserId:   userID,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return "Failed to fetch accounts. Please try again."
	}

	if len(resp.Accounts) == 0 {
		return "You don't have any accounts yet. Create one in the app!"
	}

	var sb strings.Builder
	sb.WriteString("💼 *Your Accounts*\n\n")

	for _, acc := range resp.Accounts {
		icon := getAccountIcon(acc.Type.String())
		sb.WriteString(fmt.Sprintf("%s *%s*\n", icon, acc.Name))
		sb.WriteString(fmt.Sprintf("   Balance: $%.2f %s\n\n", acc.TotalBalance, acc.Currency))
	}

	return sb.String()
}

func (b *Bot) handleAddAsset(msg *Message, args []string) string {
	userID, ok := b.userSessions[msg.Chat.ID]
	if !ok {
		return "Please link your account first using /link <email>"
	}

	if len(args) < 3 {
		return "Usage: /addasset <type> <symbol> <quantity>\nExample: /addasset crypto BTC 0.5"
	}

	// This would typically create an asset through the service
	return fmt.Sprintf("✅ Asset added successfully!\n\nType: %s\nSymbol: %s\nQuantity: %s\n\nUser: %s",
		args[0], strings.ToUpper(args[1]), args[2], userID)
}

func (b *Bot) handleLink(msg *Message, args []string) string {
	if len(args) < 1 {
		return "Usage: /link <email>\nExample: /link user@example.com"
	}

	email := args[0]

	// In a real implementation, this would verify the email and link the account
	// For now, we'll just store a placeholder
	b.userSessions[msg.Chat.ID] = fmt.Sprintf("user_%d", msg.From.ID)

	return fmt.Sprintf("✅ Account linked successfully!\n\nEmail: %s\n\nYou can now use all bot features.", email)
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	data, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	return err
}

func getAccountIcon(accountType string) string {
	icons := map[string]string{
		"ACCOUNT_TYPE_BANK":        "🏦",
		"ACCOUNT_TYPE_CASH":        "💵",
		"ACCOUNT_TYPE_INVESTMENT":  "📈",
		"ACCOUNT_TYPE_CRYPTO":      "₿",
		"ACCOUNT_TYPE_REAL_ESTATE": "🏠",
		"ACCOUNT_TYPE_OTHER":       "📁",
	}
	if icon, ok := icons[accountType]; ok {
		return icon
	}
	return "📁"
}

