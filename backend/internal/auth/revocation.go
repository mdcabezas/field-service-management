package auth

import (
	"sync"
	"time"
)

type RevocationService struct {
	mu       sync.RWMutex
	revoked  map[string]time.Time
	stopChan chan struct{}
	stopOnce sync.Once
}

func NewRevocationService() *RevocationService {
	r := &RevocationService{
		revoked:  make(map[string]time.Time),
		stopChan: make(chan struct{}),
	}
	go r.cleanup()
	return r
}

func (r *RevocationService) Revoke(jti string, expiresAt time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.revoked[jti] = expiresAt
}

func (r *RevocationService) IsRevoked(jti string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.revoked[jti]
	return exists
}

func (r *RevocationService) Stop() {
	r.stopOnce.Do(func() { close(r.stopChan) })
}

func (r *RevocationService) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.mu.Lock()
			now := time.Now()
			for jti, exp := range r.revoked {
				if now.After(exp) {
					delete(r.revoked, jti)
				}
			}
			r.mu.Unlock()
		}
	}
}
