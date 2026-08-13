package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
	maxSize  int
	stopChan chan struct{}
	stopOnce sync.Once
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
		maxSize:  10000,
		stopChan: make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() { close(rl.stopChan) })
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for {
		select {
		case <-rl.stopChan:
			return
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		if len(rl.visitors) >= rl.maxSize {
			rl.evictOldest()
		}
		rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if time.Since(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	if v.count >= rl.limit {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

func (rl *RateLimiter) evictOldest() {
	var oldestIP string
	var oldestTime time.Time
	for ip, v := range rl.visitors {
		if oldestIP == "" || v.lastSeen.Before(oldestTime) {
			oldestIP = ip
			oldestTime = v.lastSeen
		}
	}
	if oldestIP != "" {
		delete(rl.visitors, oldestIP)
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow(c.ClientIP()) {
			c.Header("Retry-After", strconv.Itoa(int(rl.window.Seconds())))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
