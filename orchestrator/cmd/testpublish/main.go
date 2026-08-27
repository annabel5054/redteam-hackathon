package main

import (
	"context"
	"fmt"
	"log"

	"github.com/willzfrank/redteam-hackathon/orchestrator/internal/eventbus"
)

func main() {
	rdb := eventbus.NewClient()
	ctx := context.Background()

	err := eventbus.Publish(ctx, rdb, eventbus.Event{
		Agent: "attacker",
		Type:  "attack_attempt",
		Payload: map[string]interface{}{
			"endpoint": "/accounts/02",
			"method":   "GET",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Event published successfully")
}