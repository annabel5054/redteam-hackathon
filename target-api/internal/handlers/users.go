package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/willzfrank/redteam-hackathon/target-api/internal/store"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	user, ok := store.GetUser(userID)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
