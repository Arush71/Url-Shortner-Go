package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Arush71/url-shortener/internal/helpers"
)

type LimiterS struct {
	counter  int
	FirstHit time.Time
	mu       sync.Mutex
}

type IpHolder map[string]*LimiterS
type IPManager struct {
	Holder IpHolder
	muCtr  sync.Mutex
}

func (ipM *IPManager) CheckRateLimit(ip string, limit int) bool {
	ipM.muCtr.Lock()
	value, ok := ipM.Holder[ip]
	if !ok {
		ipM.Holder[ip] = &LimiterS{
			counter:  1,
			FirstHit: time.Now(),
		}
		ipM.muCtr.Unlock()
		return true
	}
	ipM.muCtr.Unlock()
	value.mu.Lock()
	defer value.mu.Unlock()
	if time.Since(value.FirstHit) > time.Minute {
		value.counter = 1
		value.FirstHit = time.Now()
		return true
	}
	if value.counter >= limit {
		return false
	}
	value.counter++
	return true
}
func SetupIpManager() *IPManager {
	return &IPManager{
		Holder: make(IpHolder),
	}
}
func GetClientIp(RemoteA string) string {

	// Note: Switch to proper proxy setup and header stuff.
	host, _, err := net.SplitHostPort(RemoteA)
	if err != nil {
		slog.Warn("failed to parse RemoteAddr, falling back to full address",
			"remote_addr", RemoteA,
			"error", err,
		)
		return RemoteA
	}
	return host
}

func (M *IPManager) CleanUpIp() {
	ticker := time.NewTicker(2 * time.Minute)
	for {
		<-ticker.C
		M.muCtr.Lock()
		for k, v := range M.Holder {
			v.mu.Lock()
			expired := time.Since(v.FirstHit) > 2*time.Minute
			v.mu.Unlock()
			if expired {
				delete(M.Holder, k)
			}
		}
		M.muCtr.Unlock()
	}
}

func (M *IPManager) RateLimitMiddleware(limit int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIp(r.RemoteAddr)
			if !M.CheckRateLimit(ip, limit) {
				helpers.WriteError(w, http.StatusTooManyRequests, helpers.ErrorResponse{
					Error:   "rate_limited",
					Message: "Too many requests. Please try again later.",
				})
				return
			}
			next(w, r)
		}
	}
}
