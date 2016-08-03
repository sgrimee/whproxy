package main

import (
	"encoding/hex"
	"encoding/json"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	webhooksPath = "webhooks"
)

// hookServe hanles an incoming webhook by reading json from it and
// passing it on to the corresponding open websocket
// then the webhook connection is closed
func hookServe(hw http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming webhook from %s: %q\n", r.RemoteAddr, html.EscapeString(r.URL.Path))
	hSig := r.Header.Get("X-Spark-Signature")
	sig, err := hex.DecodeString(hSig)
	if err != nil {
		log.Printf("Unable to decode signature: %q", hSig)
		hw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if config.Validate && sig == nil {
		log.Println("Discarding webhook because non hex or missing X-Spark-Signature and Validation mode is on")
		hw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	log.Printf("  with signature: %q\n", sig)
	id := r.URL.Path[len("/"+webhooksPath+"/"):]
	var s Session
	var ok bool
	if s, ok = session(id); !ok {
		log.Printf("Error: webhook %s not found.\n", id)
		hw.WriteHeader(http.StatusNotFound)
		return
	}
	if s.Connection == nil {
		log.Printf("Error: webhook %s points to invalid websocket.\n", id)
		deleteSession(id)
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
	log.Printf("Body: \n%s\n", body)
	if config.Validate && !validSignature(body, sig, s.Secret) {
		log.Println("Discarding webhook because X-Spark-Signature is not valid.")
		hw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// we parse the message as a means of validating its json
	var hm HookMsg
	if err := json.Unmarshal(body, &hm); err != nil {
		log.Printf("Invalid json on webhook %s: %s", id, err)
		hw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		hw.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(hw).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	log.Printf("Received json on webhook %s: %+v", id, hm)
	// proxy hook message to websocket in 'data' field of a json
	var ev SocketEvent
	ev.Data = hm

	hw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(s.Connection).Encode(ev); err != nil {
		log.Printf("Could not send proxied json from %s to websocket: %s", id, err)
		deleteSession(id)
		hw.WriteHeader(http.StatusNotFound)
	}
	hw.WriteHeader(http.StatusOK)
}
