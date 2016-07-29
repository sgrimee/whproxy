package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

const (
	websocketPath = "websocket"
)

type SocketEvent struct {
	Data      HookMsg `json:"data"`
	Signature string  `json:"signature,omitempty"`
}

// json that is sent back to the websocket client open connection
type SocketResponse struct {
	Url    string `json:"url"`
	Secret string `json:"secret,omitempty"`
}

// wsServe gemerates a private webhook endpoint for each incoming websocket
// The websocket is kept open so incoming webhook data can be proxied to it
func wsServe(ws *websocket.Conn) {
	id := NewUid()
	setConn(id, ws) // keep track of open sessions
	// send private webhook endpoint to client
	data := SocketResponse{Url: fmt.Sprintf("http://%s:%d/%s/%s",
		config.Host, config.Port, webhooksPath, id)}
	log.Printf("Incoming websocket from %s, sending: %+v\n", ws.Request().RemoteAddr, data)
	websocket.JSON.Send(ws, data)
	// read forever on websocket to keep it open
	var msg []byte
	for _, err := ws.Read(msg); err == nil; time.Sleep(1 * time.Second) {
	}
	log.Printf("Releasing websocket: %s\n", id)
	deleteConn(id)
	ws.Close()
}
