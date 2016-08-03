package main

// concurrent management of open sessions

import (
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Session struct {
	Connection *websocket.Conn
	Secret     []byte
}

var (
	mu       sync.Mutex
	sessions = make(map[string]Session)
)

func session(id string) (Session, bool) {
	mu.Lock()
	defer mu.Unlock()
	s, ok := sessions[id]
	return s, ok
}

func setSession(id string, s Session) {
	mu.Lock()
	defer mu.Unlock()
	sessions[id] = s
}

func deleteSession(id string) {
	log.Printf("Removing websocket: %s\n", id)
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, id)
}

func sessionsCount() int {
	mu.Lock()
	defer mu.Unlock()
	return len(sessions)
}
