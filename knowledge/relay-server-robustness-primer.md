# Go Primer: Relay Server Robustness & Reliability

## Why Robustness?
- Prevents resource leaks and zombie sessions.
- Allows users to recover from disconnects and get clear feedback.

---

## Core Concepts & Libraries
- **Structs:** Used to track room state (connections, timestamps, flags).
- **sync.Mutex:** Ensures safe concurrent access to the room map.
- **Goroutines:** For background cleanup and piping data.
- **Error Handling:** All errors are logged, and users get clear messages.

---

## Example: Room Struct for State Tracking
```go
type room struct {
    sender    net.Conn
    receiver  net.Conn
    createdAt time.Time
    retryable bool
    lastActivity time.Time
    senderDisconnected bool
    receiverDisconnected bool
}
```
**What does this do?**
- Tracks both connections, timestamps, and retry/disconnect state for each transfer code.

---

## Example: Room Expiration & Cleanup
```go
func cleanupRooms() {
    for {
        time.Sleep(time.Minute)
        mu.Lock()
        now := time.Now()
        for code, r := range rooms {
            // Remove abandoned rooms after 10 min
            if r.sender == nil || r.receiver == nil {
                if now.Sub(r.createdAt) > 10*time.Minute {
                    delete(rooms, code)
                }
            } else if r.retryable {
                // Remove retryable rooms 2 min after disconnect
                if r.senderDisconnected || r.receiverDisconnected {
                    if now.Sub(r.lastActivity) > 2*time.Minute {
                        delete(rooms, code)
                    }
                }
            }
        }
        mu.Unlock()
    }
}
```
**What does this do?**
- Periodically scans and removes rooms that are abandoned or expired.
- Uses a goroutine and mutex for concurrency safety.

---

## Example: Retryable Rooms & Disconnect Notification
```go
func pipeWithNotify(r *room, who string) {
    var src, dst net.Conn
    if who == "sender" {
        src = r.sender
        dst = r.receiver
    } else {
        src = r.receiver
        dst = r.sender
    }
    if src == nil || dst == nil {
        return
    }
    _, err := io.Copy(dst, src)
    if err != nil {
        log.Printf("Pipe error (%s): %v", who, err)
    }
    // Notify the other side
    if dst != nil {
        dst.Write([]byte("DISCONNECT\n"))
    }
    mu.Lock()
    if who == "sender" {
        r.senderDisconnected = true
    } else {
        r.receiverDisconnected = true
    }
    r.lastActivity = time.Now()
    mu.Unlock()
}
```
**What does this do?**
- Pipes data between sender and receiver.
- If one side disconnects, notifies the other and updates room state for retry logic.

---

## Example: User Feedback on Handshake Failure
```go
if codeFromHandshake != "" {
    allowed, triesLeft, blocked, blockMsg := relay.CheckAndRecordFailedHandshake(codeFromHandshake)
    if !allowed {
        if blocked {
            conn.Write([]byte(blockMsg + "\n"))
        } else {
            msg := fmt.Sprintf("Invalid code or key. You have %d tries remaining before this code is blocked.\n", triesLeft)
            conn.Write([]byte(msg))
        }
    }
}
```
**What does this do?**
- Tells the user how many tries remain before a code is blocked, or if the code is temporarily blocked.

---

## Go Keywords & Libraries
- `struct`, `sync.Mutex`, `goroutine`, `defer`, `map`, `time.Time`, `io.Copy`, `log`, `fmt.Sprintf` 