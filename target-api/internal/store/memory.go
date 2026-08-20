package store

import "sync"

type Account struct {
	ID      string
	OwnerID string
	Balance int64
}

type User struct {
	ID               string
	Name             string
	Email            string
	SSN              string
	InternalRiskFlag bool
}

var (
	muUsers sync.Mutex
	users   = map[string]User{
		"01": {ID: "01", Name: "John Doe", Email: "john.doe@example.com", SSN: "123-45-6789", InternalRiskFlag: false},
		"02": {ID: "02", Name: "Jane Doe", Email: "jane.doe@example.com", SSN: "987-65-4321", InternalRiskFlag: true},
		"03": {ID: "03", Name: "Jim Doe", Email: "jim.doe@example.com", SSN: "111-22-3333", InternalRiskFlag: false},
	}
)

var (
	muAccounts sync.Mutex
	accounts   = map[string]Account{
		"01": {ID: "01", OwnerID: "1234567890", Balance: 1000},
		"02": {ID: "02", OwnerID: "9876543210", Balance: 500},
		"03": {ID: "03", OwnerID: "1234567891", Balance: 2000},
	}
)

func GetAccount(id string) (Account, bool) {
	muAccounts.Lock()
	defer muAccounts.Unlock()
	account, ok := accounts[id]
	return account, ok
}

func UpdateAccount(id string, acc Account) {
	if id == "" || acc.ID == "" {
		return
	}
	muAccounts.Lock()
	defer muAccounts.Unlock()
	accounts[id] = acc
}

func GetUser(id string) (User, bool) {
	muUsers.Lock()
	defer muUsers.Unlock()
	user, ok := users[id]
	return user, ok
}

// --- Rate limiting state for /login ---

type loginTracker struct {
	mu sync.Mutex
	// key -> unix timestamps (seconds) of recent attempts
	attempts map[string][]int64
}

var loginAttempts = &loginTracker{attempts: make(map[string][]int64)}

// RecordLoginAttempt records an attempt for key (e.g. IP address) at time now (unix seconds) and returns how many attempts have occurred within the trailing windowSeconds.

func RecordLoginAttempt(key string, now int64, windowSeconds int64) int {
	loginAttempts.mu.Lock()
	defer loginAttempts.mu.Unlock()

	cutoff := now - windowSeconds
	existing := loginAttempts.attempts[key]

	kept := existing[:0]
	for _, ts := range existing {
		if ts > cutoff {
			kept = append(kept, ts)
		}
	}
	kept = append(kept, now)
	loginAttempts.attempts[key] = kept

	return len(kept)
}
