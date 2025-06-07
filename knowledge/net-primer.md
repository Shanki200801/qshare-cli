# Go net Primer: TCP Server & Client

## TCP Server in Go
```go
ln, err := net.Listen("tcp", ":4000") // Listen on port 4000
if err != nil { /* handle error */ }
for {
    conn, err := ln.Accept() // Wait for a client
    if err != nil { continue }
    go handleConnection(conn) // Handle each client in a goroutine
}
```
- `net.Listen("tcp", addr)` — starts a TCP server.
- `Accept()` — waits for a client to connect.
- Each connection is a net.Conn (like a socket).
- Use goroutines to handle multiple clients at once.

**Python equivalent:**
```python
import socket
s = socket.socket()
s.bind(('', 4000))
s.listen()
while True:
    conn, addr = s.accept()
    # handle conn in a thread/process
```

**Java equivalent:**
```java
ServerSocket server = new ServerSocket(4000);
while (true) {
    Socket client = server.accept();
    // handle client in a thread
}
```

---

## TCP Client in Go
```go
conn, err := net.Dial("tcp", "server-address:4000")
if err != nil { /* handle error */ }
// Use conn to read/write
```
- `net.Dial` connects to a TCP server.
- You can use conn.Write([]byte) and conn.Read([]byte) to send/receive data.

**Python equivalent:**
```python
import socket
s = socket.socket()
s.connect(('server-address', 4000))
# s.send(), s.recv()
```

**Java equivalent:**
```java
Socket socket = new Socket("server-address", 4000);
// socket.getInputStream(), socket.getOutputStream()
```

---

## Reading/Writing Data
- Go's net.Conn implements io.Reader and io.Writer.
- You can use bufio, encoding/gob, or just send raw bytes/strings.

**Example:**
```go
conn.Write([]byte("hello\n"))
buf := make([]byte, 1024)
n, err := conn.Read(buf)
```

---

## Deployment Notes
- The relay server is just a Go binary—no dependencies, very lightweight.
- You can deploy it on any free/cheap VPS (Fly.io, Railway, Render, etc.), or even a free tier cloud VM.
- Listen on 0.0.0.0:PORT to accept connections from anywhere.
- For public deployment, consider using TLS (can add later).

---

## Summary
- Go's net package is simple and powerful for TCP.
- The relay server will be a standalone Go program, easy to deploy.
- Sender/receiver CLI will connect to it over TCP. 