package middleware

import (
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"awesomeProjectCr/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

const RequestIDKey = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)

		c.Next()
	}
}

// -- Recover ------------------------------------------------------------------

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("panic", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("stack", string(debug.Stack())).
					Msg("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, handler.ErrorResponse{
					Message: "please contact admin",
					Code:    http.StatusInternalServerError,
				})
			}
		}()
		c.Next()
	}
}

// -- Rate Limiter -------------------------------------------------------------

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
	ttl      time.Duration
}

func NewRateLimiterStore(rps int, burst int, ttl time.Duration) *RateLimiterStore {
	store := &RateLimiterStore{
		visitors: make(map[string]*visitor),
		rps:      rate.Every(time.Minute / time.Duration(rps)),
		burst:    burst,
		ttl:      ttl,
	}

	go store.cleanupLoop()

	return store
}

func (s *RateLimiterStore) getVisitor(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.visitors[key]
	if !exists {
		limiter := rate.NewLimiter(s.rps, s.burst)
		s.visitors[key] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (s *RateLimiterStore) cleanupLoop() {
	ticker := time.NewTicker(s.ttl)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for key, v := range s.visitors {
			if time.Since(v.lastSeen) > s.ttl {
				delete(s.visitors, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *RateLimiterStore) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Authorization")
		if key == "" {
			key = c.ClientIP()
		}

		limiter := s.getVisitor(key)
		if !limiter.Allow() {
			log.Warn().
				Str("key", key).
				Str("path", c.Request.URL.Path).
				Msg("rate limit exceeded")

			c.AbortWithStatusJSON(http.StatusTooManyRequests, handler.ErrorResponse{
				Message: "too many requests, please slow down",
				Code:    http.StatusTooManyRequests,
			})
			return
		}

		c.Next()
	}
}
