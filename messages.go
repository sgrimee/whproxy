package main

// json that is sent back to the websocket client open connection
type OpenHookResponse struct {
	Url string `json:"url"`
}

// acceptable json to foward
// Fully open for now
type HookMsg interface{}

// json sent back on /status requests
type StatusResponse struct {
	NbOpenHooks int `json:"nb_open_hooks"`
}
