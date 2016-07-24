package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func TestMain(m *testing.M) {
	//go ListenAndServe()
	os.Exit(m.Run())
}

// should receive the text that was sent
func TestDialEcho(t *testing.T) {
	origin := "http://local"
	url := "ws://localhost:12345/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	sentText := []byte("hello, world!\n")
	if _, err := ws.Write(sentText); err != nil {
		t.Fatal(err)
	}
	var rcvText = make([]byte, 512)
	var n int
	if n, err = ws.Read(rcvText); err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(sentText, rcvText[:n]) != 0 {
		t.Fatalf("Received: %s", rcvText[:n])
	}
}

// opening a websocket should return a webhook url
func TestGetHookUrl(t *testing.T) {
	origin := "http://localhost/"
	url := "ws://localhost:12345/websocket"
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
	// TODO: add test to return hooks count via /status
	// TODO: close the websocket ?
	// TODO: add test to return hooks count via /status
}

// ensure valid json sent to webhook comes into websocket
func TestProxy(t *testing.T) {
	origin := "http://localhost/"
	url := "ws://localhost:12345/websocket"
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
		a string
	}
	sMsg := testMsg{a: "test"}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(sMsg); err != nil {
		panic(err)
	}
	t.Log("Sending json on webhook: %s", b)
	res, err := http.Post(r.Url, "application/json; charset=utf-8", b)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatal("Error when POSTing test json: %s", res.Status)
	}
	// expect same json incoming on websocket
	var rMsg HookMsg
	if err := dec.Decode(&rMsg); err != nil {
		t.Fatal("Error decoding json from webhook: ", err)
	}
	t.Logf("Received message on webhook: %+v\n", rMsg)
	if rMsg != sMsg {
		t.Fatal("Received message is not what was sent: %+v", rMsg)
	}
}
