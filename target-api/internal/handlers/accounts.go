package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willzfrank/redteam-hackathon/target-api/internal/store"
)

func HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	account, ok := store.GetAccount(accountID)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if account.OwnerID != userID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
