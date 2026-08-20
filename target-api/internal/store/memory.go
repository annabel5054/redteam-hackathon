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
		"01": {
			ID:               "01",
			Name:             "John Doe",
			Email:            "john.doe@example.com",
			SSN:              "1234567890",
			InternalRiskFlag: false,
		},
		"02": {
			ID:               "02",
			Name:             "Jane Doe",
			Email:            "jane.doe@example.com",
			SSN:              "1234567890",
			InternalRiskFlag: true,
		},
		"03": {
			ID:               "03",
			Name:             "Jim Doe",
			Email:            "jim.doe@example.com",
			SSN:              "1234567891",
			InternalRiskFlag: false,
		},
	}
)

var (
	muAccounts sync.Mutex
	accounts   = map[string]Account{
		"01": {
			ID:      "01",
			OwnerID: "1234567890",
			Balance: 1000,
		},
		"02": {
			ID:      "02",
			OwnerID: "9876543210",
			Balance: 500,
		},
		"03": {
			ID:      "03",
			OwnerID: "1234567891",
			Balance: 2000,
		},
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
	if id == "" {
		return User{}, false
	}
	if id != accounts[id].OwnerID {
		return User{}, false
	}
	user, ok := users[id]
	return user, ok
}
