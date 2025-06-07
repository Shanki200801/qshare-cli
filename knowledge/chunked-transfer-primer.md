# Chunked File Transfer & Framing Primer (Go, Python, Java)

## Why Chunked Transfer?
- Reading/writing files in chunks avoids loading the whole file into memory.
- Essential for large files and streaming.

## Framing: How to Know Where Each Chunk Starts/Ends?
- Prepend each chunk with its length (as a fixed-size header, e.g., 4 bytes, big-endian).
- Receiver reads the length, then reads that many bytes for the chunk.

---

## Go Example: Sending Chunks
```go
buf := make([]byte, 64*1024) // 64KB buffer
for {
    n, err := file.Read(buf)
    if n > 0 {
        chunk := buf[:n]
        // Encrypt chunk if needed
        // Write length header
        binary.Write(conn, binary.BigEndian, uint32(len(chunk)))
        // Write chunk
        conn.Write(chunk)
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        // handle error
        break
    }
}
// Optionally, send a zero-length chunk to signal EOF
```

## Go Example: Receiving Chunks
```go
for {
    var chunkLen uint32
    err := binary.Read(conn, binary.BigEndian, &chunkLen)
    if err != nil {
        break // EOF or error
    }
    if chunkLen == 0 {
        break // End of file
    }
    chunk := make([]byte, chunkLen)
    _, err = io.ReadFull(conn, chunk)
    if err != nil {
        break // handle error
    }
    // Decrypt chunk if needed
    // Write chunk to file
}
```

---

## Python Equivalent
```python
import struct
# Sender:
chunk = f.read(65536)
conn.sendall(struct.pack('>I', len(chunk)))
conn.sendall(chunk)
# Receiver:
size = struct.unpack('>I', conn.recv(4))[0]
chunk = conn.recv(size)
```

## Java Equivalent
```java
// Sender:
out.writeInt(chunk.length);
out.write(chunk);
// Receiver:
int size = in.readInt();
byte[] chunk = new byte[size];
in.readFully(chunk);
```

---

## Summary
- Chunked transfer = read/write in pieces, not all at once.
- Use a length header (framing) so the receiver knows how much to read.
- Works the same way in Go, Python, and Java.
- For encrypted transfer, encrypt each chunk before sending, decrypt after receiving. 