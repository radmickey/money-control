package proxy

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	accountspb "github.com/radmickey/money-control/backend/proto/accounts"
	assetspb "github.com/radmickey/money-control/backend/proto/assets"
	authpb "github.com/radmickey/money-control/backend/proto/auth"
	currencypb "github.com/radmickey/money-control/backend/proto/currency"
	insightspb "github.com/radmickey/money-control/backend/proto/insights"
	transactionspb "github.com/radmickey/money-control/backend/proto/transactions"
)

// Default timeouts and retry configuration
const (
	DefaultConnectTimeout = 5 * time.Second
	DefaultRequestTimeout = 10 * time.Second
	MaxRetryAttempts      = 3
)

// Retry policy for gRPC calls - handles transient failures
var retryPolicy = `{
	"methodConfig": [{
		"name": [{"service": ""}],
		"waitForReady": true,
		"timeout": "10s",
		"retryPolicy": {
			"maxAttempts": 3,
			"initialBackoff": "0.1s",
			"maxBackoff": "1s",
			"backoffMultiplier": 2.0,
			"retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED", "RESOURCE_EXHAUSTED"]
		}
	}]
}`

// Keepalive parameters to maintain connection health
var keepaliveParams = keepalive.ClientParameters{
	Time:                10 * time.Second, // Ping server every 10s if no activity
	Timeout:             3 * time.Second,  // Wait 3s for ping ack
	PermitWithoutStream: true,             // Allow pings even without active streams
}

// Config holds service proxy configuration
type Config struct {
	AuthServiceURL         string
	AccountsServiceURL     string
	TransactionsServiceURL string
	AssetsServiceURL       string
	CurrencyServiceURL     string
	InsightsServiceURL     string
}

// ServiceProxy manages gRPC connections to microservices
type ServiceProxy struct {
	connections []*grpc.ClientConn

	Auth         authpb.AuthServiceClient
	Accounts     accountspb.AccountsServiceClient
	Transactions transactionspb.TransactionsServiceClient
	Assets       assetspb.AssetsServiceClient
	Currency     currencypb.CurrencyServiceClient
	Insights     insightspb.InsightsServiceClient
}

// NewServiceProxy creates a new service proxy with resilient connections
func NewServiceProxy(cfg Config) (*ServiceProxy, error) {
	sp := &ServiceProxy{
		connections: make([]*grpc.ClientConn, 0),
	}

	// Common dial options for all connections
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithKeepaliveParams(keepaliveParams),
		// Connection pooling via round-robin when multiple backends available
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}

	// Connect to Auth service
	if cfg.AuthServiceURL != "" {
		conn, err := grpc.Dial(cfg.AuthServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Auth = authpb.NewAuthServiceClient(conn)
	}

	// Connect to Accounts service
	if cfg.AccountsServiceURL != "" {
		conn, err := grpc.Dial(cfg.AccountsServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Accounts = accountspb.NewAccountsServiceClient(conn)
	}

	// Connect to Transactions service
	if cfg.TransactionsServiceURL != "" {
		conn, err := grpc.Dial(cfg.TransactionsServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Transactions = transactionspb.NewTransactionsServiceClient(conn)
	}

	// Connect to Assets service
	if cfg.AssetsServiceURL != "" {
		conn, err := grpc.Dial(cfg.AssetsServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Assets = assetspb.NewAssetsServiceClient(conn)
	}

	// Connect to Currency service
	if cfg.CurrencyServiceURL != "" {
		conn, err := grpc.Dial(cfg.CurrencyServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Currency = currencypb.NewCurrencyServiceClient(conn)
	}

	// Connect to Insights service
	if cfg.InsightsServiceURL != "" {
		conn, err := grpc.Dial(cfg.InsightsServiceURL, opts...)
		if err != nil {
			return nil, err
		}
		sp.connections = append(sp.connections, conn)
		sp.Insights = insightspb.NewInsightsServiceClient(conn)
	}

	return sp, nil
}

// Close closes all gRPC connections
func (sp *ServiceProxy) Close() {
	for _, conn := range sp.connections {
		conn.Close()
	}
}
