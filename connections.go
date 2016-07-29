package main

// concurrent management of open connections

import (
	"golang.org/x/net/websocket"
	"sync"
)

var (
	mu          sync.Mutex
	connections = make(map[string]*websocket.Conn)
)

// func init() {
// 	hooks = make(map[string]*websocket.Conn)
// }

func conn(id string) (*websocket.Conn, bool) {
	mu.Lock()
	defer mu.Unlock()
	c, ok := connections[id]
	return c, ok
}

func setConn(id string, c *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	connections[id] = c
}

func deleteConn(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(connections, id)
}

func connCount() int {
	mu.Lock()
	defer mu.Unlock()
	return len(connections)
}
