package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/willzfrank/redteam-hackathon/target-api/internal/store"
)

// UserResponse is the public-safe shape sent to clients — deliberately
// excludes SSN and InternalRiskFlag.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

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

	resp := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
