package main

import (
	"context"
	"log"

	"github.com/willzfrank/redteam-hackathon/orchestrator/internal/agents"
	"github.com/willzfrank/redteam-hackathon/orchestrator/internal/eventbus"
)

func main() {
	rdb := eventbus.NewClient()
	ctx := context.Background()

	if err := agents.RunAttacker(ctx, rdb); err != nil {
		log.Fatal(err)
	}
}