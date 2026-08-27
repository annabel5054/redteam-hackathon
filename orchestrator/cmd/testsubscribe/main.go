package main

import (
	"context"
	"fmt"

	"github.com/willzfrank/redteam-hackathon/orchestrator/internal/eventbus"
)

func main() {
	rdb := eventbus.NewClient()
	ctx := context.Background()

	fmt.Println("Listening for events...")
	err := eventbus.Subscribe(ctx, rdb, "$", func(evt eventbus.Event) {
		fmt.Printf("Received: agent=%s type=%s payload=%v\n", evt.Agent, evt.Type, evt.Payload)
	})
	if err != nil {
		fmt.Println("error:", err)
	}
}
