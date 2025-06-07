package relay

import (
	"sync"
	"time"
)

const (
	ipLimit         = 5
	codeLimit       = 5
	window          = time.Minute
	failedWindow    = 5 * time.Minute
	failedThreshold = 3
	blockDuration   = 10 * time.Minute
)

var (
	ipAttempts       = make(map[string][]time.Time)
	codeAttempts     = make(map[string][]time.Time)
	mu               sync.Mutex
	failedHandshakes = make(map[string][]time.Time)
	blockedCodes     = make(map[string]time.Time)
)

// CheckAndRecordRateLimit checks and records an attempt for the given IP and code.
// Returns (true, "") if allowed, (false, reason) if rate limited.
func CheckAndRecordRateLimit(ip, code string) (bool, string) {
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	// Clean up old entries
	clean := func(attempts []time.Time) []time.Time {
		cutoff := now.Add(-window)
		for len(attempts) > 0 && attempts[0].Before(cutoff) {
			attempts = attempts[1:]
		}
		return attempts
	}
	ipAttempts[ip] = clean(ipAttempts[ip])
	codeAttempts[code] = clean(codeAttempts[code])
	if len(ipAttempts[ip]) >= ipLimit {
		return false, "rate limit exceeded for IP"
	}
	if len(codeAttempts[code]) >= codeLimit {
		return false, "rate limit exceeded for code"
	}
	// Record this attempt
	ipAttempts[ip] = append(ipAttempts[ip], now)
	codeAttempts[code] = append(codeAttempts[code], now)
	return true, ""
}

// CheckAndRecordFailedHandshake tracks failed handshakes per code. Returns (allowed, triesLeft, blocked, blockMsg)
func CheckAndRecordFailedHandshake(code string) (bool, int, bool, string) {
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	// Check if code is blocked
	if until, ok := blockedCodes[code]; ok {
		if now.Before(until) {
			return false, 0, true, "Code temporarily blocked due to too many failed attempts. Try again later."
		}
		delete(blockedCodes, code)
	}
	// Clean up old failed attempts
	cutoff := now.Add(-failedWindow)
	attempts := failedHandshakes[code]
	for len(attempts) > 0 && attempts[0].Before(cutoff) {
		attempts = attempts[1:]
	}
	failedHandshakes[code] = attempts
	// Record this failed attempt
	attempts = append(attempts, now)
	failedHandshakes[code] = attempts
	triesLeft := failedThreshold - len(attempts)
	if triesLeft <= 0 {
		blockedCodes[code] = now.Add(blockDuration)
		return false, 0, true, "Code temporarily blocked due to too many failed attempts. Try again later."
	}
	return true, triesLeft, false, ""
}

// Periodically clean up old entries to avoid memory leaks
func StartRateLimitCleanup() {
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			now := time.Now()
			cutoff := now.Add(-window)
			for ip, attempts := range ipAttempts {
				ipAttempts[ip] = filterRecent(attempts, cutoff)
				if len(ipAttempts[ip]) == 0 {
					delete(ipAttempts, ip)
				}
			}
			for code, attempts := range codeAttempts {
				codeAttempts[code] = filterRecent(attempts, cutoff)
				if len(codeAttempts[code]) == 0 {
					delete(codeAttempts, code)
				}
			}
			// Clean up failed handshakes and blocked codes
			failedCutoff := now.Add(-failedWindow)
			for code, attempts := range failedHandshakes {
				failedHandshakes[code] = filterRecent(attempts, failedCutoff)
				if len(failedHandshakes[code]) == 0 {
					delete(failedHandshakes, code)
				}
			}
			for code, until := range blockedCodes {
				if now.After(until) {
					delete(blockedCodes, code)
				}
			}
			mu.Unlock()
		}
	}()
}

func filterRecent(attempts []time.Time, cutoff time.Time) []time.Time {
	for len(attempts) > 0 && attempts[0].Before(cutoff) {
		attempts = attempts[1:]
	}
	return attempts
}
