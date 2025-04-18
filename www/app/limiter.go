package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Get(string, any) error
	Set(string, any) error
}

type TokenBucket struct {
	Tokens     float64
	MaxTokens  float64
	RefillRate float64 // tokens per second
	LastRefill time.Time
	mu         sync.Mutex
}

func (tb *TokenBucket) MarshalBinary() ([]byte, error) {
	return json.Marshal(tb)
}

func (tb *TokenBucket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tb)
}

// RateLimiter controls the rate of requests
type RateLimiter struct {
	store        Storage
	tokenspersec float64
	maxburst     float64
	identfn      func(*http.Request) string
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(store Storage, tokenspersec, maxburst float64, identfn func(*http.Request) string) *RateLimiter {
	return &RateLimiter{
		store:        store,
		tokenspersec: tokenspersec,
		maxburst:     maxburst,
		identfn:      identfn,
	}
}

func (rl *RateLimiter) bucket(identifier string) (*TokenBucket, error) {
	var bucket = new(TokenBucket)
	err := rl.store.Get(identifier, bucket)
	if err != nil {
		if err == redis.Nil {
			// Create a new bucket if it doesn't exist
			newbckt := &TokenBucket{
				Tokens:     rl.maxburst,
				MaxTokens:  rl.maxburst,
				RefillRate: rl.tokenspersec,
				LastRefill: time.Now(),
			}
			err = rl.store.Set(identifier, newbckt)
			if err != nil {
				return nil, err
			}
			return newbckt, nil
		}
		return nil, err
	}

	return bucket, nil
}

// Allow checks if a request is allowed based on rate limits
func (rl *RateLimiter) Allow(ident string) (bool, error) {
	bucket, err := rl.bucket(ident)
	if err != nil {
		return false, err
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	bucket.LastRefill = now

	// Add tokens based on time elapsed
	bucket.Tokens = min(bucket.Tokens+(elapsed*bucket.RefillRate), bucket.MaxTokens)

	// Check if we have at least one token
	if bucket.Tokens < 1 {
		return false, nil
	}

	// Consume a token
	bucket.Tokens--

	// Update the bucket in the store
	if err := rl.store.Set(ident, bucket); err != nil {
		return false, err
	}

	return true, nil
}

// Middleware returns HTTP middleware that applies rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get identifier for this request (e.g., IP address)
		identifier := rl.identfn(r)

		allowed, err := rl.Allow(identifier)
		fmt.Println(allowed, err)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			w.Header().Add("Retry-After", "1")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
