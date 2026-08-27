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

// AccountPatchRequest is meant to be the allowlist of fields a client can
// change. VULNERABILITY: this handler ignores the allowlist and instead
// merges the entire raw request body onto the account — mass assignment.
func HandlePatchAccount(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")
	userID := r.Header.Get("X-User-Id")

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	account, ok := store.GetAccount(accountID)
	if !ok || account.OwnerID != userID {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// VULNERABLE: decoding directly into the domain struct instead of a
	// restricted DTO means ANY field present in the request body —
	// including Balance — gets applied.
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store.UpdateAccount(accountID, account)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
