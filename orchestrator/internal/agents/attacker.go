package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/willzfrank/redteam-hackathon/orchestrator/internal/eventbus"
)

const targetBaseURL = "http://localhost:8081"

type attempt struct {
	name        string
	method      string
	path        string
	headers     map[string]string
	body        string
	expectBlock bool
}

func attackerChecklist() []attempt {
	return []attempt{
		{
			name:        "BOLA: read account 02 while claiming to own account 01",
			method:      "GET",
			path:        "/accounts/02",
			headers:     map[string]string{"X-User-Id": "1234567890"},
			expectBlock: true,
		},
		{
			name:        "BOLA: read account 01 with no user id at all",
			method:      "GET",
			path:        "/accounts/01",
			headers:     map[string]string{},
			expectBlock: true,
		},
		{
			name:        "Data exposure: check /users response for SSN leakage",
			method:      "GET",
			path:        "/users",
			headers:     map[string]string{"X-User-Id": "01"},
			expectBlock: false,
		},
		{
			name:        "Mass assignment: patch account 01 nickname, sneak in balance",
			method:      "PATCH",
			path:        "/accounts/01",
			headers:     map[string]string{"X-User-Id": "1234567890"},
			body:        `{"nickname":"Renamed","balance":999999999}`,
			expectBlock: true,
		},
	}
}

func RunAttacker(ctx context.Context, rdb *redis.Client) error {
	client := &http.Client{Timeout: 5 * time.Second}

	for _, a := range attackerChecklist() {
		var bodyReader io.Reader
		if a.body != "" {
			bodyReader = strings.NewReader(a.body)
		}

		req, err := http.NewRequest(a.method, targetBaseURL+a.path, bodyReader)
		if err != nil {
			return err
		}
		if a.body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		for k, v := range a.headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()

		vulnerable := detectVulnerability(a, resp.StatusCode, body)

		evtType := "attack_failed"
		if vulnerable {
			evtType = "vulnerability_found"
		}

		err = eventbus.Publish(ctx, rdb, eventbus.Event{
			Agent: "attacker",
			Type:  evtType,
			Payload: map[string]interface{}{
				"attempt":    a.name,
				"path":       a.path,
				"statusCode": resp.StatusCode,
				"vulnerable": vulnerable,
			},
		})
		if err != nil {
			return err
		}

		fmt.Printf("[attacker] %s -> status=%d vulnerable=%v\n", a.name, resp.StatusCode, vulnerable)
	}

	return nil
}

func detectVulnerability(a attempt, status int, body map[string]interface{}) bool {
	switch a.name {
	case "BOLA: read account 02 while claiming to own account 01",
		"BOLA: read account 01 with no user id at all":
		return status == http.StatusOK

	case "Data exposure: check /users response for SSN leakage":
		_, hasSSN := body["ssn"]
		_, hasRiskFlag := body["internalRiskFlag"]
		return hasSSN || hasRiskFlag

	case "Mass assignment: patch account 01 nickname, sneak in balance":
		balance, ok := body["Balance"].(float64)
		return ok && balance == 999999999

	default:
		return false
	}
}