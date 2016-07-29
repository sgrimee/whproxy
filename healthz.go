package main

import (
	"encoding/json"
	"net/http"
)

const (
	healthzPath = "healthz"
)

// json sent back on /status requests
type HealthzResponse struct {
	NbOpenHooks int `json:"nb_open_hooks"`
}

// healthzServe returns infos on running server
// currently the number of currently open hooks
func healthzServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	s := &HealthzResponse{
		NbOpenHooks: connCount(),
	}
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}

// getHealthz queries the /status endpoint (used by tests)
func getHealthz(url string) (*HealthzResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	var sr HealthzResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
