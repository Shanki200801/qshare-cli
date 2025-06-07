package relay

import (
	"sync"
)

// Relay is a struct that contains a mutex and a map of rooms
type Relay struct {
	mu sync.Mutex
	rooms map[string]chan []byte
}

// NewRelay creates a new relay instance and returns a pointer to it
// acts like constructors in other languages
func NewRelay() *Relay {
	return &Relay{
		rooms: make(map[string]chan []byte),
	}
}

// Sender calls this to create a room and get a channel
func (r *Relay) CreateRoom(code string) chan []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan []byte)
	r.rooms[code] = ch
	return ch
}

// Receiver calls this to join a room by code
func (r *Relay) JoinRoom(code string) (chan []byte, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.rooms[code]
	return ch, ok
}