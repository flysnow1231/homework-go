package middleware

import (
	"net/http"
	"sync"
	"time"

	"blog/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	lastSeen map[string]time.Time
	r        rate.Limit
	burst    int
	ttl      time.Duration
}

func NewIPRateLimiter(r rate.Limit, burst int, ttl time.Duration) *IPLimiter {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &IPLimiter{
		limiters: make(map[string]*rate.Limiter),
		lastSeen: make(map[string]time.Time),
		r:        r,
		burst:    burst,
		ttl:      ttl,
	}
}

func (m *IPLimiter) get(ip string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	lim, ok := m.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(m.r, m.burst)
		m.limiters[ip] = lim
	}
	m.lastSeen[ip] = time.Now()

	for k, t := range m.lastSeen {
		if time.Since(t) > m.ttl {
			delete(m.lastSeen, k)
			delete(m.limiters, k)
		}
	}
	return lim
}

func (m *IPLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !m.get(ip).Allow() {
			resp.Fail(c, http.StatusTooManyRequests, "rate_limited")
			c.Abort()
			return
		}
		c.Next()
	}
}
