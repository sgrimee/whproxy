package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var (
	version string
)

type Config struct {
	Host     string
	Port     int
	Validate bool
}

var (
	config Config
)

func main() {
	flag.StringVar(&config.Host, "host", "localhost", "hostname for the webhook server")
	flag.IntVar(&config.Port, "port", 12345, "port for the webhook server")
	flag.BoolVar(&config.Validate, "validate", false, "validate signature of incoming webhooks (WIP)")
	showVer := flag.Bool("version", false, "show server version and exit")
	flag.Parse()
	if *showVer {
		fmt.Println("Version: ", version)
		return
	}
	ListenAndServe()
}

func ListenAndServe() {
	http.Handle("/"+websocketPath, websocket.Handler(wsServe))
	http.HandleFunc("/"+webhooksPath+"/", hookServe)
	http.HandleFunc("/"+healthzPath, healthzServe)

	log.Printf("Server starting on %s:%d\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
