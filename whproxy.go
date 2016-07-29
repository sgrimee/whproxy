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
	Port,
	SSLPort int
	Validate bool
}

var (
	config Config
)

func main() {
	flag.StringVar(&config.CertFile, "cert", "", "certificate file")
	flag.StringVar(&config.KeyFile, "key", "", "key file")
	flag.StringVar(&config.Host, "host", "localhost", "hostname for webhook url")
	flag.IntVar(&config.Port, "port", 12345, "port for the webhook server")
	flag.IntVar(&config.SSLPort, "sslport", 12346, "SSL port for the webhook server")
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
	var errs = make(chan error)

	http.Handle("/"+websocketPath, websocket.Handler(wsServe))
	http.HandleFunc("/"+webhooksPath+"/", hookServe)
	http.HandleFunc("/"+healthzPath, healthzServe)

	go func() {
		log.Printf("Server starting on %s:%d\n", config.Host, config.Port)
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	}()
	if (config.CertFile != "") && (config.KeyFile != "") {
		go func() {
			log.Printf("SSL Server starting on %s:%d\n", config.Host, config.SSLPort)
			errs <- http.ListenAndServeTLS(fmt.Sprintf(":%d", config.SSLPort),
				config.CertFile, config.KeyFile, nil)
		}()
	}
	log.Fatal(<-errs) // block until one of the servers exits
}
