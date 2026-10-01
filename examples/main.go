// Example usage of the Blocklog Go SDK
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blockloglabs/blocklog-go-sdk"
	"github.com/blockloglabs/blocklog-go-sdk/api"
	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/google/uuid"
)

func main() {
	// Initialize client from environment variables
	// Set BLOCKLOG_API_KEY or BLOCKLOG_ACCESS_TOKEN in your environment
	client, err := blocklog.InitFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}
	defer client.Shutdown(context.Background())

	// Or initialize with explicit config
	// cfg := config.NewConfig()
	// cfg.APIKey = "blk_..."
	// cfg.Endpoint = "https://blocklogsecurity.com/api/v1"
	// client, err := blocklog.Init(cfg)

	ctx := context.Background()

	// Example 1: Record a simple event
	fmt.Println("=== Recording Events ===")
	resp, err := client.Event(ctx, "USER_LOGIN", map[string]interface{}{
		"user_id": "user_123",
		"method":  "oauth",
	}, map[string]interface{}{
		"source": "my-app",
	})
	if err != nil {
		log.Printf("Event error: %v", err)
	} else {
		fmt.Printf("Event recorded: %+v\n", resp)
	}

	// Example 2: Create a decision
	fmt.Println("\n=== Creating Decision ===")
	traceID := uuid.New()
	sessionID := uuid.New()
	decision, err := client.Decisions.Create(ctx, &models.DecisionCreateRequest{
		DecisionType: "TRADE_EXECUTE",
		Agent:        strPtr("trading-bot"),
		AgentID:      strPtr("bot_001"),
		Model:        strPtr("gpt-4"),
		Asset:        strPtr("AAPL"),
		Confidence:   float64Ptr(0.95),
		Inputs: map[string]interface{}{
			"price":    150.25,
			"quantity": 100,
			"signal":   "BUY",
		},
		Outputs: map[string]interface{}{
			"order_id":     "ord_abc123",
			"filled_qty":   100,
			"filled_price": 150.30,
		},
		TraceID:   &traceID,
		SessionID: &sessionID,
	})
	if err != nil {
		log.Printf("Decision error: %v", err)
	} else {
		fmt.Printf("Decision created: %s\n", decision.ID)
	}

	// Example 3: Create an incident
	fmt.Println("\n=== Creating Incident ===")
	incident, err := client.Incidents.Create(ctx, &models.IncidentCreateRequest{
		Title:       "Anomalous trading pattern detected",
		Severity:    models.IncidentSeverityHigh,
		Description: strPtr("Multiple rapid trades detected outside normal parameters"),
		TraceID:     decision.TraceID,
		Metadata: map[string]interface{}{
			"detection_rule": "velocity_check",
			"threshold":      1000,
		},
	})
	if err != nil {
		log.Printf("Incident error: %v", err)
	} else {
		fmt.Printf("Incident created: %s (status: %s)\n", incident.ID, incident.Status)

		// Use fluent handle for incident lifecycle
		handle := api.NewIncidentHandle(incident, client.Incidents)
		handle.Assign(ctx, "analyst@company.com", strPtr("Investigating unusual volume"))
		handle.Resolve(ctx, "False positive - authorized batch trade", nil, nil)
		handle.Close(ctx, "Reviewed and approved", "approved")
	}

	// Example 4: Request approval (HITL)
	fmt.Println("\n=== Requesting Approval ===")
	approval, err := client.Approvals.Request(ctx, &models.ApprovalRequest{
		DecisionID: strPtr(decision.ID),
		Reason:     "Trade exceeds $500k threshold, requires manual review",
		Reviewer:   strPtr("risk-team@company.com"),
		Metadata: map[string]interface{}{
			"threshold": 500000,
			"amount":    750000,
		},
	})
	if err != nil {
		log.Printf("Approval error: %v", err)
	} else {
		fmt.Printf("Approval requested: %s\n", approval.ID)
	}

	// Example 5: Create a replay session for forensics
	fmt.Println("\n=== Creating Replay Session ===")
	replay, err := client.Replay.Create(ctx, &models.ReplayCreateRequest{
		TraceID: traceID,
		Metadata: map[string]interface{}{
			"purpose": "investigation",
		},
	})
	if err != nil {
		log.Printf("Replay error: %v", err)
	} else {
		fmt.Printf("Replay session created: %s\n", replay.ID)

		// Wrap with ReplaySession for fluent API
		replaySession := api.NewReplaySession(replay, client.Replay)

		// Get root cause analysis
		rootCause, err := replaySession.RootCause(ctx)
		if err != nil {
			log.Printf("Root cause error: %v", err)
		} else {
			fmt.Printf("Root cause: %+v\n", rootCause)
		}

		// Get causal graph
		causalGraph, err := replaySession.CausalGraph(ctx)
		if err != nil {
			log.Printf("Causal graph error: %v", err)
		} else {
			fmt.Printf("Causal graph nodes: %d, edges: %d\n", len(causalGraph.Nodes), len(causalGraph.Edges))
		}
	}

	// Example 6: Generate compliance report
	fmt.Println("\n=== Generating Compliance Report ===")
	traceIDStr := traceID.String()
	report, err := client.Compliance.Generate(ctx, &models.ComplianceGenerateRequest{
		Title:      "SOC2 Trading Compliance Report",
		ReportKind: "design_partner_readiness",
		ScopeType:  "trace",
		Framework:  strPtr("SOC2"),
		TraceID:    &traceIDStr,
		DateFrom:   timePtr(time.Now().AddDate(0, -1, 0)),
		DateTo:     timePtr(time.Now()),
	})
	if err != nil {
		log.Printf("Compliance error: %v", err)
	} else {
		fmt.Printf("Report generated: %s\n", report.ID)
	}

	// Example 7: Verify a decision
	fmt.Println("\n=== Verifying Decision ===")
	verify, err := client.Verify.Decision(ctx, decision.ID)
	if err != nil {
		log.Printf("Verify error: %v", err)
	} else {
		fmt.Printf("Verification status: %v\n", verify.Status)
	}

	// Example 8: List teams
	fmt.Println("\n=== Listing Teams ===")
	teams, err := client.Teams.List(ctx)
	if err != nil {
		log.Printf("Teams error: %v", err)
	} else {
		for _, team := range teams {
			fmt.Printf("Team: %s (%s)\n", team.Name, team.Slug)
		}
	}

	// Example 9: Query traces
	fmt.Println("\n=== Querying Traces ===")
	eventType := "DECISION_COMPLETE"
	traces, err := client.Traces.List(ctx, &models.TraceListParams{
		EventType: &eventType,
		Limit:     10,
	})
	if err != nil {
		log.Printf("Traces error: %v", err)
	} else {
		fmt.Printf("Found %d traces\n", len(traces.Items))
	}

	// Example 10: Batch events (enqueue and flush)
	fmt.Println("\n=== Batch Events ===")
	for i := 0; i < 5; i++ {
		client.Enqueue(ctx, "METRIC_RECORD", map[string]interface{}{
			"metric": "cpu_usage",
			"value":  45.2 + float64(i),
			"host":   "server-01",
		}, map[string]interface{}{})
	}

	flushResp, err := client.Flush(ctx)
	if err != nil {
		log.Printf("Flush error: %v", err)
	} else {
		fmt.Printf("Flushed %d events\n", flushResp.Ingested)
	}

	fmt.Println("\n=== All examples completed ===")
}

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
