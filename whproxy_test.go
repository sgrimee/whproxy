package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func TestMain(m *testing.M) {
	config.Host = "localhost"
	config.Port = 9876
	config.Validate = false
	ssl := false
	go ListenAndServe(ssl)
	os.Exit(m.Run())
}

// opening a websocket should return a webhook url
func TestGetHookUrl(t *testing.T) {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/%s", config.Host, config.Port, websocketPath)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
	// decode one message from the stream
	var r SocketResponse
	websocket.JSON.Receive(ws, &r)
	t.Logf("Received webhook URL: %s\n", r.Url)
	if !strings.HasPrefix(r.Url, "http") {
		t.Error("Invalid webhook url: %s", r.Url)
	}
	// make sure hook is removed and only after closing the websocket
	url = fmt.Sprintf("http://%s:%d/%s", config.Host, config.Port, healthzPath)
	var sr *HealthzResponse
	if sr, err = getHealthz(url); err != nil {
		t.Error(err)
	}
	if sr.NbOpenHooks != 1 {
		t.Error("Incorrect number of open hooks before closing: %d", sr.NbOpenHooks)
	}
	ws.Close()
	if sr, err = getHealthz(url); err != nil {
		t.Error(err)
	}
	if sr.NbOpenHooks != 0 {
		t.Error("Incorrect number of open hooks after closing: %d", sr.NbOpenHooks)
	}
}

// ensure valid json sent to webhook comes into websocket
func TestProxy(t *testing.T) {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/%s", config.Host, config.Port, websocketPath)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
	dec := json.NewDecoder(ws)
	// decode inital response on websocket
	var r SocketResponse
	if err := dec.Decode(&r); err != nil {
		t.Error(err)
	}
	t.Log("Received webhook URL: ", r.Url)
	if !strings.HasPrefix(r.Url, "http") {
		t.Error("Invalid webhook url: ", r.Url)
	}
	// send json on webhook
	b := []byte(`{"test": "abc123"}`)
	t.Log("Sending json on webhook: ", string(b[:]))
	res, err := http.Post(r.Url, "application/json; charset=utf-8", bytes.NewReader(b))
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Error("Error when POSTing test json: ", res.Status)
	}
	// expect same json incoming on websocket
	var rMsg SocketEvent
	if err := dec.Decode(&rMsg); err != nil {
		t.Error("Error decoding json from websocket: ", err)
	}
	t.Logf("Received event on websocket: %+v\n", rMsg)
	var m map[string]interface{} = rMsg.Data.(map[string]interface{})
	if m["test"] != "abc123" {
		t.Error("Received message differs from expected: %+v (%T)\n", m, m)
	}
}

func TestSignature(t *testing.T) {
	message := []byte(`{"testing":"123"}`)
	secret := []byte("d42813e6-0b3d-4142-be90-af8f5a2c508v")
	knownSig, _ := hex.DecodeString("992fabf5cbdfaaeec1c607bc478f3a8cc88c1797")
	if !validSignature(message, knownSig, secret) {
		t.Error("Computed signature differs from known-good")
	}

}
