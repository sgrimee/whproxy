package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	Host string
	Port int
}

var (
	server Server
	hooks  map[string]*websocket.Conn
)

func init() {
	hooks = make(map[string]*websocket.Conn)
}

func main() {
	flag.StringVar(&server.Host, "host", "localhost", "hostname for the webhook server")
	flag.IntVar(&server.Port, "port", 12345, "port for the webhook server")
	flag.Parse()
	ListenAndServe()
}

func ListenAndServe() {
	http.Handle("/"+websocketPath, websocket.Handler(wsServe))
	http.HandleFunc("/"+webhooksPath+"/", hookServe)
	http.HandleFunc("/"+healthzPath, healthzServe)

	log.Printf("Server starting on %s:%d\n", server.Host, server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil))
}
