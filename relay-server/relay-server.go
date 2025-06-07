package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"bufio"
)

type room struct {
	sender net.Conn
	receiver net.Conn
}

var (
	rooms = make(map[string]*room)
	mu sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp",":4000")
	// if error, log and exit
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	log.Println("Server is running on port 4000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() //Defer runs after function exits, i.e it's similar to finally in other language
	code, role, err := handshake(conn)
	if err != nil {
		log.Printf("Handshake failed: %v", err)
		return
	}
	mu.Lock()
	r, ok := rooms[code]
	if !ok {
		r = &room{}
		rooms[code] = r
	}
	if role == "sender" {
		r.sender = conn
	} else {
		r.receiver = conn
	}
	// If both sender and receiver are set, start the piping
	if r.sender != nil && r.receiver != nil {
		go pipe(r.sender, r.receiver)
		go pipe(r.receiver, r.sender)
		delete(rooms, code)
	}
	mu.Unlock()
	//Block here until the connection is closed
	select {}
}

func handshake(conn net.Conn) (code, role string, err error){
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err

	}

	//Expect: code:role\n (e.g 5-sku-transfer:sender)
	n, _ := fmt.Sscanf(line, "%s[^:]:%s\n", &code, &role)
	if n != 2 {
		return "", "", fmt.Errorf("invalid handshake: %q", line)
	}
	return code, role, nil
}

func pipe(src, dst net.Conn) {
	io.Copy(dst, src)
}