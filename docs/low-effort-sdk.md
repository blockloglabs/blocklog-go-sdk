# Low-effort instrumentation strategy

## Vision

Blocklog is execution, control, and verification infrastructure for autonomous
AI agents. Manual instrumentation remains the reliable baseline, while future
integrations should make common execution boundaries observable with a small,
explicit opt-in.

## Current boundary

The core SDK exposes a client, event ingestion, and the canonical
`models.AgentEvent` envelope. It deliberately does not rewrite source code,
patch global HTTP state, capture payloads implicitly, or claim that receipts are
available before the backend issues them.

## Proposed integration shape

```text
framework adapter -> lifecycle hook -> AgentEvent collector -> Client
```

Adapters should be separate packages and must depend on a small hook interface,
not on framework-specific logic in the core client. Hooks should receive context,
capture start/end/error and parent relationships, and permit a caller-provided
redaction policy before data leaves the process.

Initial candidates are `net/http` client/server middleware, OpenTelemetry-aware
services, popular Go LLM clients, MCP clients/servers, SQL/database wrappers,
and job/queue workers. Each integration must be opt-in, composable, no-op when
unconfigured, and independently versioned.

## Future `blocklog init`

A future CLI may inspect declared dependencies, propose supported adapters, and
write only explicitly approved configuration. It must never silently modify
application code or enable raw prompt/body capture. Projects need opt-out at the
integration, event-type, and field-redaction levels.
