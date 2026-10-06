# Go SDK architecture

The Go SDK is a thin, context-aware client for Blocklog's API. It owns request
construction, authentication, bounded retries for safe requests, event batching,
and serialization. The backend owns validation, persistence, policy evaluation,
and issuance of verification receipts.

```text
Go application -> Blocklog Go SDK -> HTTP/API layer -> Blocklog backend
                                                     -> execution/event storage
                                                     -> verification/receipts
```

`Client` is safe for normal concurrent use. Reuse one client per application.
Calls accept `context.Context`; cancellation and deadlines are propagated to the
HTTP request and retry wait. GET requests are retried for transient transport,
rate-limit, and server errors. Mutating API calls are not automatically retried;
event ingestion supplies an idempotency key when a retry-safe backend contract is
available.

The portable `models.AgentEvent` envelope establishes the cross-SDK boundary:
execution ID, event ID, type, timestamp, sequence, parent relationship, status,
payload, metadata, provenance, and error. `LLMCall`, `ToolCall`, and
`ExecutionReceipt` are data types only today; they do not imply unsupported API
endpoints.
