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

// json that is sent back to the websocket client open connection
type OpenHookResponse struct {
	Url string `json:"url"`
}

// wsServe gemerates a private webhook endpoint for each incoming websocket
// The websocket is kept open so incoming webhook data can be proxied to it
func wsServe(ws *websocket.Conn) {
	id := NewUid()
	hooks[id] = ws // keep track of open sessions
	// send private webhook endpoint to client
	data := OpenHookResponse{Url: fmt.Sprintf("http://%s:%d/%s/%s",
		server.Host, server.Port, webhooksPath, id)}
	log.Println("Incoming websocket, sending: %+v", data)
	websocket.JSON.Send(ws, data)
	// read forever on websocket to keep it open
	var msg []byte
	for _, err := ws.Read(msg); err == nil; time.Sleep(1 * time.Second) {
	}
	log.Println("Releasing websocket: %s", id)
	delete(hooks, id)
	ws.Close()
}
