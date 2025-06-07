# Go Primer: Rate Limiting & Abuse Prevention in a TCP Relay Server

## Why Rate Limiting?
- Prevents abuse (DoS, brute-force, flooding) by limiting how often clients can connect.
- Protects both per-user (IP) and per-resource (code) usage.

---

## Core Concepts & Libraries
- **sync.Mutex**: Ensures safe concurrent access to shared maps.
- **time.Time, time.Duration**: Used for tracking and comparing timestamps.
- **Goroutines**: Background tasks for periodic cleanup.
- **Maps**: Store per-IP and per-code attempt history.

---

## Example: Per-IP and Per-Code Rate Limiting
```go
var (
    ipAttempts   = make(map[string][]time.Time)
    codeAttempts = make(map[string][]time.Time)
    mu           sync.Mutex
)

const (
    ipLimit   = 5
    codeLimit = 5
    window    = time.Minute
)

func CheckAndRecordRateLimit(ip, code string) (bool, string) {
    mu.Lock()
    defer mu.Unlock()
    now := time.Now()
    // Clean up old entries
    ipAttempts[ip] = filterRecent(ipAttempts[ip], now.Add(-window))
    codeAttempts[code] = filterRecent(codeAttempts[code], now.Add(-window))
    if len(ipAttempts[ip]) >= ipLimit {
        return false, "rate limit exceeded for IP"
    }
    if len(codeAttempts[code]) >= codeLimit {
        return false, "rate limit exceeded for code"
    }
    ipAttempts[ip] = append(ipAttempts[ip], now)
    codeAttempts[code] = append(codeAttempts[code], now)
    return true, ""
}

func filterRecent(attempts []time.Time, cutoff time.Time) []time.Time {
    for len(attempts) > 0 && attempts[0].Before(cutoff) {
        attempts = attempts[1:]
    }
    return attempts
}
```
**What does this do?**
- Tracks connection attempts per IP and per code in a rolling 1-minute window.
- Uses a mutex to prevent race conditions (multiple goroutines accessing the maps).
- If the limit is exceeded, returns false and a reason.

---

## Example: Failed Handshake Lockout
```go
var (
    failedHandshakes = make(map[string][]time.Time)
    blockedCodes     = make(map[string]time.Time)
)
const (
    failedWindow   = 5 * time.Minute
    failedThreshold = 3
    blockDuration   = 10 * time.Minute
)

func CheckAndRecordFailedHandshake(code string) (bool, int, bool, string) {
    mu.Lock()
    defer mu.Unlock()
    now := time.Now()
    if until, ok := blockedCodes[code]; ok && now.Before(until) {
        return false, 0, true, "Code temporarily blocked due to too many failed attempts. Try again later."
    }
    // Clean up old failed attempts
    cutoff := now.Add(-failedWindow)
    attempts := filterRecent(failedHandshakes[code], cutoff)
    failedHandshakes[code] = attempts
    attempts = append(attempts, now)
    failedHandshakes[code] = attempts
    triesLeft := failedThreshold - len(attempts)
    if triesLeft <= 0 {
        blockedCodes[code] = now.Add(blockDuration)
        return false, 0, true, "Code temporarily blocked due to too many failed attempts. Try again later."
    }
    return true, triesLeft, false, ""
}
```
**What does this do?**
- Tracks failed handshake attempts per code.
- If a code gets 3 failed attempts in 5 minutes, it is blocked for 10 minutes.
- Returns how many tries are left before block.

---

## Periodic Cleanup with Goroutines
```go
func StartRateLimitCleanup() {
    go func() {
        for {
            time.Sleep(time.Minute)
            mu.Lock()
            // ...clean up old entries from all maps...
            mu.Unlock()
        }
    }()
}
```
**What does this do?**
- Runs a background goroutine to periodically clean up old entries, preventing memory leaks.

---

## How to Extend
- Change thresholds by editing the constants.
- For distributed/multi-instance servers, use a shared store (e.g., Redis) instead of in-memory maps.
- Add more granular limits (e.g., per user agent) by tracking additional keys.

---

## Go Keywords & Libraries
- `sync.Mutex`, `map`, `time.Time`, `time.Duration`, `goroutine`, `defer`, `make`, `append`
- No third-party libraries required for this logic.

---

## Implementation Overview
- **Per-IP Limit:** Max 5 connections per minute per IP address.
- **Per-Code Limit:** Max 5 connections per minute per code (transfer code).
- **Failed Handshake Lockout:** If a code gets 3 failed handshakes in 5 minutes, it is blocked for 10 minutes.
- **Cleanup:** Old entries are periodically removed to avoid memory leaks.

---

## Go Concepts Used
- **Maps:** Used to track attempts per IP and per code.
- **sync.Mutex:** Ensures safe concurrent access to maps.
- **Goroutines:** Used for periodic cleanup in the background.
- **Time Windows:** Attempts are only counted within a rolling window (e.g., last 1 or 5 minutes).

---

## Extending/Customizing
- Change thresholds by editing constants in the rate limiting file.
- Add more granular limits (e.g., per user agent) by tracking additional keys.
- For distributed systems, use Redis or another shared store instead of in-memory maps.

---

## Keywords
- DoS protection, brute-force prevention, concurrency, mutex, goroutine, map, time window, lockout, cleanup. 