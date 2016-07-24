package main

import (
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var (
	hooks map[string]*websocket.Conn
)

func main() {
	hooks = make(map[string]*websocket.Conn)

	http.Handle("/echo", websocket.Handler(echo))
	http.Handle("/websocket", websocket.Handler(wsServe))
	http.HandleFunc("/webhooks/", hookServe)
	http.HandleFunc("/status", statusServe)

	log.Fatal(http.ListenAndServe(":12345", nil))
}

// Echo the data received on the WebSocket.
func echo(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func wsServe(ws *websocket.Conn) {
	id := NewId()
	log.Println("Incoming websocket:, sending id: %s", id)
	hooks[id] = ws // keep track of open sessions
	r := OpenHookResponse{Url: "http://localhost:12345/webhooks/" + id}
	json.NewEncoder(ws).Encode(r)
	// keep reading on websocket to keep it open
	msg := byte.Buffer
	for _, err := ws.Read(msg); err == nil; time.Sleep(1 * time.Second) {
	}
	log.Println("Releasing websocket: %s", id)
	delete(hooks, id)
	ws.Close()
}

func hookServe(hw http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming webhook: %q\n", html.EscapeString(r.URL.Path))
	id := r.URL.Path[len("/webhooks/"):]
	var ws *websocket.Conn
	var ok bool
	if ws, ok = hooks[id]; !ok {
		log.Printf("Id not found: %s\n", id)
		hw.WriteHeader(http.StatusNotFound)
		return
	}
	if ws == nil {
		log.Printf("Id %s points to invalid websocket, we stop tracking it.\n", id)
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
	// send message on websocket
	json.NewEncoder(ws).Encode(message)
	// send ok on webhook
	hw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	hw.WriteHeader(http.StatusOK)
}

func statusServe(w http.ResponseWriter, r *http.Request) {
	log.Printf("Open hooks: %+v", hooks)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	s := &StatusResponse{
		NbOpenHooks: len(hooks),
	}
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}
