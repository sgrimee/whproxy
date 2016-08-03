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
func wsServe(c *websocket.Conn) {
	defer c.Close()
	id := NewUid()
	secret := []byte(NewUid())
	setSession(id, Session{c, secret})
	defer deleteSession(id)
	// send private webhook endpoint to client
	data := SocketResponse{
		Url:    fmt.Sprintf("%s://%s:%d/%s/%s", scheme, config.Host, config.Port, webhooksPath, id),
		Secret: string(secret),
	}
	log.Printf("Incoming websocket from %s, sending: %+v\n", c.Request().RemoteAddr, data)
	websocket.JSON.Send(c, data)
	// read forever on websocket to keep it open
	var msg []byte
	for _, err := c.Read(msg); err == nil; time.Sleep(1 * time.Second) {
	}
}
