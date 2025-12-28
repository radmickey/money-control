package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache provides Redis caching functionality
type Cache struct {
	client *redis.Client
	prefix string
}

// Config holds Redis configuration
type Config struct {
	URL    string
	Prefix string
}

// New creates a new Redis cache instance
func New(cfg Config) (*Cache, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "mc"
	}

	return &Cache{
		client: client,
		prefix: prefix,
	}, nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Client returns the underlying Redis client
func (c *Cache) Client() *redis.Client {
	return c.client
}

// key generates a prefixed cache key
func (c *Cache) key(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Set stores a value in cache with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.Set(ctx, c.key(key), data, expiration).Err()
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// GetString retrieves a string value from cache
func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, c.key(key)).Result()
}

// SetString stores a string value in cache
func (c *Cache) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	return c.client.Set(ctx, c.key(key), value, expiration).Err()
}

// Delete removes a value from cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.key(k)
	}
	return c.client.Del(ctx, prefixedKeys...).Err()
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, c.key(key)).Result()
	return n > 0, err
}

// SetNX sets a value if key doesn't exist (for distributed locks)
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.SetNX(ctx, c.key(key), data, expiration).Result()
}

// Incr increments a counter
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.key(key)).Result()
}

// Expire sets expiration on a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, c.key(key), expiration).Err()
}

// HSet sets hash fields
func (c *Cache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.client.HSet(ctx, c.key(key), values...).Err()
}

// HGet gets a hash field
func (c *Cache) HGet(ctx context.Context, key, field string) (string, error) {
	return c.client.HGet(ctx, c.key(key), field).Result()
}

// HGetAll gets all hash fields
func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, c.key(key)).Result()
}

// Keys cache key constants
const (
	KeyAssetPrice    = "asset:price:%s"
	KeyExchangeRate  = "currency:rate:%s:%s"
	KeyExchangeRates = "currency:rates:%s"
	KeyUserSession   = "session:%s"
	KeyRateLimit     = "ratelimit:%s:%s"
)

// AssetPriceKey generates cache key for asset price
func AssetPriceKey(symbol string) string {
	return fmt.Sprintf(KeyAssetPrice, symbol)
}

// ExchangeRateKey generates cache key for exchange rate
func ExchangeRateKey(from, to string) string {
	return fmt.Sprintf(KeyExchangeRate, from, to)
}

// ExchangeRatesKey generates cache key for all rates with base currency
func ExchangeRatesKey(base string) string {
	return fmt.Sprintf(KeyExchangeRates, base)
}

// UserSessionKey generates cache key for user session
func UserSessionKey(userID string) string {
	return fmt.Sprintf(KeyUserSession, userID)
}

// RateLimitKey generates cache key for rate limiting
func RateLimitKey(identifier, endpoint string) string {
	return fmt.Sprintf(KeyRateLimit, identifier, endpoint)
}

