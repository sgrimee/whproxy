package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func TestMain(m *testing.M) {
	server.Host = "localhost"
	server.Port = 9876
	go ListenAndServe()
	os.Exit(m.Run())
}

// opening a websocket should return a webhook url
func TestGetHookUrl(t *testing.T) {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/%s", server.Host, server.Port, websocketPath)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	// decode one message from the stream
	var r OpenHookResponse
	websocket.JSON.Receive(ws, &r)
	t.Logf("Received webhook URL: %s\n", r.Url)
	if !strings.HasPrefix(r.Url, "http") {
		t.Fatal("Invalid webhook url: %s", r.Url)
	}
	// make sure hook is removed and only after closing the websocket
	url = fmt.Sprintf("http://%s:%d/%s", server.Host, server.Port, healthzPath)
	var sr *HealthzResponse
	if sr, err = getHealthz(url); err != nil {
		t.Fatal(err)
	}
	if sr.NbOpenHooks != 1 {
		t.Fatal("Incorrect number of open hooks before closing: %d", sr.NbOpenHooks)
	}
	ws.Close()
	if sr, err = getHealthz(url); err != nil {
		t.Fatal(err)
	}
	if sr.NbOpenHooks != 0 {
		t.Fatal("Incorrect number of open hooks after closing: %d", sr.NbOpenHooks)
	}
}

// ensure valid json sent to webhook comes into websocket
func TestProxy(t *testing.T) {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/%s", server.Host, server.Port, websocketPath)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(ws)
	// decode one message from the stream
	var r OpenHookResponse
	if err := dec.Decode(&r); err != nil {
		t.Fatal(err)
	}
	t.Logf("Received webhook URL: %s\n", r.Url)
	if !strings.HasPrefix(r.Url, "http") {
		t.Fatal("Invalid webhook url: %s", r.Url)
	}
	// send json on webhook
	type testMsg struct {
		A string
	}
	sMsg := testMsg{A: "test"}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(sMsg); err != nil {
		panic(err)
	}
	t.Log("Sending json on webhook: %s", b.String())
	res, err := http.Post(r.Url, "application/json; charset=utf-8", b)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatal("Error when POSTing test json: %s", res.Status)
	}
	// expect same json incoming on websocket
	var rMsg testMsg
	if err := dec.Decode(&rMsg); err != nil {
		t.Fatal("Error decoding json from webhook: ", err)
	}
	t.Logf("Received message on webhook: %+v\n", rMsg)
	if rMsg != sMsg {
		t.Fatal("Received message is not what was sent: %+v (%T)", rMsg, rMsg)
	}
}
