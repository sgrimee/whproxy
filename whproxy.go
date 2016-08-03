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
	CertFile,
	KeyFile,
	Host string
	Port     int
	Validate bool
}

var (
	config Config
	ssl    = false
	scheme = "http"
)

func main() {
	flag.StringVar(&config.CertFile, "cert", "", "certificate file (with 'key', activates ssl)")
	flag.StringVar(&config.KeyFile, "key", "", "key file (with 'cert', activates ssl)")
	flag.StringVar(&config.Host, "host", "localhost", "hostname for webhook url")
	flag.IntVar(&config.Port, "port", 12345, "port for the webhook server")
	flag.BoolVar(&config.Validate, "validate", false, "validate signature of incoming webhooks (WIP)")
	showVer := flag.Bool("version", false, "show server version and exit")
	flag.Parse()
	if *showVer {
		fmt.Println("Version: ", version)
		return
	}
	if (config.CertFile != "") && (config.KeyFile != "") {
		ssl = true
		scheme = "https"
	}
	ListenAndServe(ssl)
}

func ListenAndServe(ssl bool) {
	http.Handle("/"+websocketPath, websocket.Handler(wsServe))
	http.HandleFunc("/"+webhooksPath+"/", hookServe)
	http.HandleFunc("/"+healthzPath, healthzServe)

	if ssl {
		log.Printf("SSL Server starting on %s:%d\n", config.Host, config.Port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Port),
			config.CertFile, config.KeyFile, nil))
	} else {
		log.Printf("Server starting on %s:%d\n", config.Host, config.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
	}
}
