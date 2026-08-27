package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const StreamName = "security-events"

type Event struct {
	Agent     string                 `json:"agent"`
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// NewClient connects to a local Redis instance.
func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// Publish serializes an Event and adds it to the security-events stream.
func Publish(ctx context.Context, rdb *redis.Client, evt Event) error {
	if evt.Timestamp == 0 {
		evt.Timestamp = time.Now().Unix()
	}

	payloadJSON, err := json.Marshal(evt.Payload)
	if err != nil {
		return err
	}

	return rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: map[string]interface{}{
			"agent":     evt.Agent,
			"type":      evt.Type,
			"timestamp": evt.Timestamp,
			"payload":   string(payloadJSON),
		},
	}).Err()
}

// Subscribe blocks and reads new events from the stream as they arrive,
// calling handler for each one. lastID controls where to start reading —
// use "$" to only get events published from this moment forward.
func Subscribe(ctx context.Context, rdb *redis.Client, lastID string, handler func(Event)) error {
	for {
		streams, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{StreamName, lastID},
			Block:   0, // block forever until a new event arrives
			Count:   10,
		}).Result()
		if err != nil {
			return err
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				var evt Event
				evt.Agent, _ = msg.Values["agent"].(string)
				evt.Type, _ = msg.Values["type"].(string)

				var payload map[string]interface{}
				if payloadStr, ok := msg.Values["payload"].(string); ok {
					json.Unmarshal([]byte(payloadStr), &payload)
				}
				evt.Payload = payload

				handler(evt)
				lastID = msg.ID // move forward so we don't reprocess this one
			}
		}
	}
}