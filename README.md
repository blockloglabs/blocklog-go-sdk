# Blocklog Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/blockloglabs/blocklog-go-sdk)](https://goreportcard.com/report/github.com/blockloglabs/blocklog-go-sdk)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Official Go SDK for [Blocklog](https://blocklogsecurity.com) — Infrastructure for AI Decision-Making.

Record every decision your AI agents make with cryptographic audit trails, human-in-the-loop approvals, forensic replay, and compliance reporting.

## Features

- **Decision Recording** — Capture AI decisions with inputs, outputs, confidence scores, and full context
- **Event Ingestion** — High-throughput event pipeline with batching, retries, and idempotency
- **Incident Management** — Full incident lifecycle: create, assign, annotate, resolve, close
- **Human-in-the-Loop (HITL)** — Request approvals, escalations, and audit trails
- **Forensic Replay** — Root cause analysis, causal graphs, counterfactuals, trace comparison
- **Compliance Reports** — SOC2, GDPR, custom frameworks with sharing and export
- **Cryptographic Verification** — Ed25519 signatures, Merkle proofs, receipt verification
- **Team Management** — Members, roles, notifications, webhooks
- **Trace & Session Queries** — Filter by trace, session, workflow, event type

## Installation

```bash
go get github.com/blockloglabs/blocklog-go-sdk
```

Requires Go 1.22+

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blockloglabs/blocklog-go-sdk"
	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/google/uuid"
)

func main() {
	// Initialize from environment variables (BLOCKLOG_API_KEY or BLOCKLOG_ACCESS_TOKEN)
	client, err := blocklog.InitFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown(context.Background())

	ctx := context.Background()

	// Record an AI decision
	traceID := uuid.New()
	decision, err := client.Decisions.Create(ctx, &models.DecisionCreateRequest{
		DecisionType: "TRADE_EXECUTE",
		Agent:        strPtr("trading-bot"),
		Asset:        strPtr("AAPL"),
		Confidence:   float64Ptr(0.95),
		Inputs: map[string]interface{}{
			"price": 150.25, "quantity": 100, "signal": "BUY",
		},
		Outputs: map[string]interface{}{
			"order_id": "ord_123", "filled_price": 150.30,
		},
		TraceID: &traceID,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decision recorded: %s\n", decision.ID)

	// Verify cryptographically
	verify, _ := client.Verify.Decision(ctx, decision.ID)
	fmt.Printf("Verified: %v\n", verify.Status)
}

func strPtr(s string) *string { return &s }
func float64Ptr(f float64) *float64 { return &f }
```

## Configuration

| Environment Variable | Description | Default |
|---------------------|-------------|---------|
| `BLOCKLOG_API_KEY` | API key for authentication | Required |
| `BLOCKLOG_ACCESS_TOKEN` | OAuth access token | Alternative to API key |
| `BLOCKLOG_ENDPOINT` | API endpoint | `https://blocklogsecurity.com/api/v1` |
| `BLOCKLOG_TIMEOUT` | Request timeout (seconds) | `10` |
| `BLOCKLOG_RETRY_COUNT` | Max retry attempts | `3` |
| `BLOCKLOG_BATCH_SIZE` | Event batch size | `100` |
| `BLOCKLOG_FLUSH_INTERVAL` | Auto-flush interval (seconds) | `2` |
| `BLOCKLOG_DEBUG` | Enable debug logging | `false` |
| `BLOCKLOG_SIGNING_KEY` | Ed25519 private key for signing | Optional |
| `BLOCKLOG_SIGNING_ALG` | Signing algorithm | `ed25519` |

Or configure programmatically:

```go
cfg := config.NewConfig()
cfg.APIKey = "blk_..."
cfg.Endpoint = "https://api.mycompany.com/v1"
cfg.Timeout = 30 * time.Second
cfg.RetryCount = 5
cfg.Debug = true

client, err := blocklog.Init(cfg)
```

## Core Concepts

### Client Initialization

```go
// Global singleton (recommended for most apps)
client, _ := blocklog.InitFromEnv()
defer client.Shutdown(ctx)

// Or explicit client for advanced use cases
client, _ := blocklog.NewClient(&blocklog.ClientOptions{
    Config: cfg,
})
```

### Event Recording

```go
// Immediate event (bypasses buffer)
client.Event(ctx, "USER_ACTION", map[string]interface{}{
    "user_id": "123", "action": "login",
}, map[string]interface{}{"source": "web"})

// Batched event (high throughput)
client.Enqueue(ctx, "METRIC", map[string]interface{}{
    "cpu": 45.2, "host": "server-1",
}, nil)

// Flush batched events
client.Flush(ctx)
```

### Decision Recording

```go
decision, _ := client.Decisions.Create(ctx, &models.DecisionCreateRequest{
    DecisionType: "CLASSIFICATION",
    Agent:        strPtr("fraud-detector"),
    Model:        strPtr("xgboost-v3"),
    Confidence:   float64Ptr(0.98),
    Inputs: map[string]interface{}{"transaction": txData},
    Outputs: map[string]interface{}{"is_fraud": true, "score": 0.98},
    TraceID:    &traceID,
    SessionID:  &sessionID,
})
```

### Incident Lifecycle

```go
// Create incident
incident, _ := client.Incidents.Create(ctx, &models.IncidentCreateRequest{
    Title:       "Anomalous pattern detected",
    Severity:    models.IncidentSeverityHigh,
    Description: strPtr("Details..."),
    TraceID:     decision.TraceID,
})

// Fluent handle for lifecycle
handle := api.NewIncidentHandle(incident, client.Incidents)
handle.Assign(ctx, "analyst@co.com", strPtr("Investigating"))
handle.Annotate(ctx, "Found root cause", strPtr("analyst"))
handle.Resolve(ctx, "False positive", nil, nil)
handle.Close(ctx, "Reviewed", "approved")
```

### Human-in-the-Loop Approvals

```go
// Request approval
approval, _ := client.Approvals.Request(ctx, &models.ApprovalRequest{
    DecisionID: strPtr(decision.ID),
    Reason:     "Exceeds risk threshold",
    Reviewer:   strPtr("risk-team@co.com"),
})

// Approve (by reviewer)
client.Approvals.Approve(ctx, approval.ID, map[string]interface{}{
    "status": "approved", "reason": "Risk acceptable",
})

// Reject
client.Approvals.Reject(ctx, approval.ID, "reviewer@co.com", "Too risky", nil)

// Escalate
client.Approvals.Escalate(ctx, &models.ApprovalEscalateRequest{
    CurrentReviewer:  "reviewer@co.com",
    EscalationTarget: "senior-risk@co.com",
    EscalationReason: "Requires senior review",
})
```

### Forensic Replay

```go
// Create replay session
replay, _ := client.Replay.Create(ctx, &models.ReplayCreateRequest{
    TraceID: traceID,
})

// Wrap for fluent API
session := api.NewReplaySession(replay, client.Replay)

// Root cause analysis
rootCause, _ := session.RootCause(ctx)
// rootCause.Detected, rootCause.RootCauseType, rootCause.Confidence, ...

// Causal graph
graph, _ := session.CausalGraph(ctx)
// graph.Nodes, graph.Edges

// Counterfactual (what-if)
counterfactual, _ := session.Counterfactual(ctx, tokenID, map[string]interface{}{
    "input_field": "modified_value",
})

// Compare two traces
comparison, _ := session.Compare(ctx, otherTraceID)
```

### Compliance Reports

```go
report, _ := client.Compliance.Generate(ctx, &models.ComplianceGenerateRequest{
    Title:      "SOC2 Report - Q1 2024",
    ReportKind: "design_partner_readiness",
    ScopeType:  "company",
    Framework:  strPtr("SOC2"),
    DateFrom:   timePtr(time.Now().AddDate(0, -3, 0)),
    DateTo:     timePtr(time.Now()),
})

// Share with auditors
share, _ := client.Compliance.Share(ctx, report.ID, &models.ComplianceShareRequest{
    Recipients:          []string{"auditor@firm.com"},
    ExpiresInDays:       intPtr(30),
    CreateAuditorAPIKey: true,
})
```

### Cryptographic Verification

```go
// Verify a log entry
logVerify, _ := client.Verify.Log(ctx, logID)
// logVerify.Status, logVerify.MerkleProof, logVerify.BatchProof

// Verify a batch
batchVerify, _ := client.Verify.Batch(ctx, batchID)
// batchVerify.Status, batchVerify.Signature, batchVerify.SignedAt

// Verify a decision (all evidence)
decVerify, _ := client.Verify.Decision(ctx, decision.ID)
// decVerify.Status, decVerify.MerkleProof, decVerify.Signature
```

### Team Management

```go
// List teams
teams, _ := client.Teams.List(ctx)

// Create team
team, _ := client.Teams.Create(ctx, &models.TeamCreateRequest{
    Name:               "Security Team",
    SlackWebhookURL:    strPtr("https://hooks.slack.com/..."),
    DefaultSLAMinutes:  intPtr(60),
})

// Manage members
client.Teams.Members.Add(ctx, team.ID, &models.TeamMemberAddRequest{
    Email:              strPtr("new@co.com"),
    Role:               models.TeamRoleAdmin,
    IsOnCall:           true,
    NotificationChannels: []models.TeamNotificationChannel{
        models.TeamNotificationChannelSlack,
        models.TeamNotificationChannelEmail,
    },
})
```

### Tracing & Session Queries

```go
// List traces with filters
traces, _ := client.Traces.List(ctx, &models.TraceListParams{
    EventType: strPtr("DECISION_COMPLETE"),
    FromTS:    timePtr(time.Now().Add(-24 * time.Hour)),
    Limit:     50,
})

// Session timeline
timeline, _ := client.Traces.SessionTimeline(ctx, sessionID.String(), &models.SessionTimelineParams{
    Limit: 100,
})
```

## Advanced Usage

### Custom HTTP Client

```go
tr := transport.NewTransport(
    "https://custom-endpoint.com",
    "api-key",
    "",
    30*time.Second,
    false,
)
tr.client = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client, _ := blocklog.NewClient(&blocklog.ClientOptions{
    Config:    cfg,
    Transport: tr,
})
```

### Request/Response Interceptors

```go
client.Transport().AddRequestInterceptor(func(method, path, url string, headers map[string]string, body []byte) {
    log.Printf("→ %s %s", method, url)
})

client.Transport().AddResponseInterceptor(func(method, path string, status int, ok bool) {
    log.Printf("← %s %s [%d]", method, path, status)
})
```

### Middleware Hooks

```go
client.AddHook(func(payload map[string]interface{}) map[string]interface{} {
    // Add correlation IDs, enrichment, etc.
    payload["correlation_id"] = generateCorrelationID()
    return payload
})
```

### Signing Events

```go
cfg := config.NewConfig()
cfg.SigningKey = "your-ed25519-private-key"
cfg.EnableSigning = true

// Events will be signed automatically
client.Event(ctx, "SENSITIVE_ACTION", payload, nil)
```

## Error Handling

```go
import (
    "github.com/blockloglabs/blocklog-go-sdk/transport"
)

decision, err := client.Decisions.Create(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, transport.ErrAuthentication):
        // Invalid credentials
    case errors.Is(err, transport.ErrAuthorization):
        // Insufficient permissions
    case errors.Is(err, transport.ErrNotFound):
        // Resource not found
    case errors.Is(err, transport.ErrRateLimit):
        // Rate limited - retry with backoff
    case errors.Is(err, transport.ErrValidation):
        // Invalid request payload
    case errors.Is(err, transport.ErrServerError):
        // 5xx error - retry
    default:
        // Network or unknown error
    }
}
```

## Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires BLOCKLOG_API_KEY)
go test -tags=integration ./...
```

## Project Structure

```
blocklog-go/
├── api/              # Domain API clients
│   ├── base.go       # Base HTTP client with retry
│   ├── decisions.go  # Decisions API
│   ├── incidents.go  # Incidents API + fluent handle
│   ├── approvals.go  # HITL approvals API
│   ├── traces.go     # Traces/sessions API
│   ├── replay.go     # Forensic replay API + fluent session
│   ├── compliance.go # Compliance reports API
│   ├── verify.go     # Verification API
│   ├── teams.go      # Teams API
│   ├── auth.go       # Auth API
│   ├── executions.go # Executions API
│   ├── execution_gateway.go
│   └── receipts.go   # Receipt verification
├── config/           # Configuration
├── models/           # Data models
├── transport/        # HTTP transport + retry
├── examples/         # Usage examples
└── client.go         # Main client
```

## Requirements

- Go 1.22+
- Blocklog account (get API key at https://blocklogsecurity.com)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- Documentation: https://docs.blocklogsecurity.com
- Issues: https://github.com/blockloglabs/blocklog-go-sdk/issues
- Email: support@blocklogsecurity.com