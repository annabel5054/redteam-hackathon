package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/willzfrank/redteam-hackathon/target-api/internal/store"
)

const (
	rateLimitWindowSeconds = 60
	rateLimitMaxAttempts   = 5
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	key, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		key = r.RemoteAddr
	}
	count := store.RecordLoginAttempt(key, time.Now().Unix(), rateLimitWindowSeconds)

	if count > rateLimitMaxAttempts {
		w.Header().Set("Retry-After", "60")
		http.Error(w, "Too many login attempts, try again later", http.StatusTooManyRequests)
		return
	}

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
