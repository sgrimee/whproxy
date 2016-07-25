package main

import (
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const (
	webhooksPath = "webhooks"
)

// hookServe hanles an incoming webhook by reading json from it and
// passing it on to the corresponding open websocket
// then the webhook connection is closed
func hookServe(hw http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming webhook: %q\n", html.EscapeString(r.URL.Path))
	id := r.URL.Path[len("/"+webhooksPath+"/"):]
	var ws *websocket.Conn
	var ok bool
	if ws, ok = hooks[id]; !ok {
		log.Printf("Error: webhook %s not found.\n", id)
		hw.WriteHeader(http.StatusNotFound)
		return
	}
	if ws == nil {
		log.Printf("Error: webhook %s points to invalid websocket.\n", id)
		delete(hooks, id)
		hw.WriteHeader(http.StatusNotFound)
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var message HookMsg
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("Invalid json on webhook %s: %s", id, err)
		hw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		hw.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(hw).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	log.Printf("Received json on webhook %s: %+v", id, message)
	// proxy message to websocket
	hw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(ws).Encode(message); err != nil {
		log.Printf("Could not send proxied json from %s to websocket: %s", id, err)
		delete(hooks, id)
		hw.WriteHeader(http.StatusNotFound)
	}
	hw.WriteHeader(http.StatusOK)
}
