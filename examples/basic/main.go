// basic records one event using the supported public client API.
package main

import (
	"context"
	"log"
	"os"

	blocklog "github.com/blockloglabs/blocklog-go-sdk"
)

func main() {
	client, err := blocklog.NewClient(blocklog.Config{APIKey: os.Getenv("BLOCKLOG_API_KEY")})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = client.Shutdown(context.Background()) }()

	_, err = client.Event(context.Background(), "agent.execution.completed", map[string]interface{}{
		"agent":  "support-assistant",
		"status": "completed",
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
}
