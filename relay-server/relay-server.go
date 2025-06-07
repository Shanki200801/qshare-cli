package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/shanki200801/qshare/internal/relay"
)

type room struct {
	sender               net.Conn
	receiver             net.Conn
	createdAt            time.Time
	retryable            bool
	lastActivity         time.Time
	senderDisconnected   bool
	receiverDisconnected bool
}

var (
	rooms = make(map[string]*room)
	mu    sync.Mutex
)

func main() {
	// Minimal HTTP handler for Render health check
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Println("Starting HTTP health check handler on :8080")
		http.ListenAndServe(":8080", nil)
	}()
	ln, err := net.Listen("tcp", ":4000")
	// if error, log and exit
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	log.Println("Server is running on port 4000")
	relay.StartRateLimitCleanup()
	go cleanupRooms()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		log.Printf("New connection from %s", conn.RemoteAddr())
		// Rate limiting
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		code := peekHandshakeCode(conn)
		allowed, reason := relay.CheckAndRecordRateLimit(ip, code)
		if !allowed {
			log.Printf("Connection from %s for code %s rejected: %s", ip, code, reason)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() //Defer runs after function exits, i.e it's similar to finally in other language
	code, role, retryable, err := handshakeWithRetry(conn)
	if err != nil {
		// Try to extract code from handshake for failed tracking
		codeFromHandshake := peekHandshakeCode(conn)
		if codeFromHandshake != "" {
			allowed, triesLeft, blocked, blockMsg := relay.CheckAndRecordFailedHandshake(codeFromHandshake)
			if !allowed {
				if blocked {
					conn.Write([]byte(blockMsg + "\n"))
					log.Printf("Handshake failed from %s: %v (code %s blocked)", conn.RemoteAddr(), err, codeFromHandshake)
				} else {
					msg := fmt.Sprintf("Invalid code or key. You have %d tries remaining before this code is blocked.\n", triesLeft)
					conn.Write([]byte(msg))
					log.Printf("Handshake failed from %s: %v (code %s, %d tries left)", conn.RemoteAddr(), err, codeFromHandshake, triesLeft)
				}
			}
		}
		log.Printf("Handshake failed from %s: %v", conn.RemoteAddr(), err)
		return
	}
	mu.Lock()
	r, ok := rooms[code]
	if !ok {
		r = &room{createdAt: time.Now(), retryable: retryable, lastActivity: time.Now()}
		rooms[code] = r
		log.Printf("Room created for code %s at %v (retryable=%v)", code, r.createdAt, retryable)
	}
	if role == "sender" {
		r.sender = conn
		r.senderDisconnected = false
		log.Printf("Sender joined room %s from %s", code, conn.RemoteAddr())
	} else {
		r.receiver = conn
		r.receiverDisconnected = false
		log.Printf("Receiver joined room %s from %s", code, conn.RemoteAddr())
	}
	r.lastActivity = time.Now()
	// If both sender and receiver are set, start the piping
	if r.sender != nil && r.receiver != nil {
		go pipeWithNotify(r, "sender")
		go pipeWithNotify(r, "receiver")
		if !r.retryable {
			delete(rooms, code)
			log.Printf("Room %s completed and deleted (not retryable)", code)
		}
	}
	mu.Unlock()
	//Block here until the connection is closed
	select {}
}

func handshakeWithRetry(conn net.Conn) (code, role string, retryable bool, err error) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", "", false, err
	}
	//Expect: code:role[:retry]\n (e.g 5-sku-transfer:sender:retry)
	var retryStr string
	n, _ := fmt.Sscanf(line, "%s[^:]:%s:%s\n", &code, &role, &retryStr)
	if n == 3 && retryStr == "retry" {
		log.Printf("Handshake: code=%s, role=%s, retryable=true, from=%s", code, role, conn.RemoteAddr())
		return code, role, true, nil
	}
	n, _ = fmt.Sscanf(line, "%s[^:]:%s\n", &code, &role)
	if n != 2 {
		return "", "", false, fmt.Errorf("invalid handshake: %q", line)
	}
	log.Printf("Handshake: code=%s, role=%s, retryable=false, from=%s", code, role, conn.RemoteAddr())
	return code, role, false, nil
}

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

// cleanupRooms periodically removes abandoned or expired rooms
func cleanupRooms() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		now := time.Now()
		for code, r := range rooms {
			if r.sender == nil || r.receiver == nil {
				if now.Sub(r.createdAt) > 10*time.Minute {
					log.Printf("Cleaning up abandoned room %s (created at %v)", code, r.createdAt)
					delete(rooms, code)
				}
			} else if r.retryable {
				if r.senderDisconnected || r.receiverDisconnected {
					if now.Sub(r.lastActivity) > 2*time.Minute {
						log.Printf("Cleaning up retryable room %s after disconnect window", code)
						delete(rooms, code)
					}
				}
			}
		}
		mu.Unlock()
	}
}

// peekHandshakeCode peeks at the handshake to extract the code for rate limiting
func peekHandshakeCode(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader := bufio.NewReader(conn)
	line, err := reader.Peek(128)
	conn.SetReadDeadline(time.Time{})
	if err != nil {
		return ""
	}
	var code, role, retry string
	n, _ := fmt.Sscanf(string(line), "%s[^:]:%s:%s\n", &code, &role, &retry)
	if n >= 1 {
		return code
	}
	n, _ = fmt.Sscanf(string(line), "%s[^:]:%s\n", &code, &role)
	if n >= 1 {
		return code
	}
	return ""
}
