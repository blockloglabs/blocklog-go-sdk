# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive unit test suite for all packages (config, transport, client, api)
- CONTRIBUTING.md with development guidelines
- Integration test examples

### Changed
- Fixed deadlock in Client.Shutdown() method
- Improved EventBuffer flushing logic

## [1.0.0] - 2024-01-15

### Added
- Initial release of Blocklog Go SDK
- Decision recording with full context (inputs, outputs, confidence, trace/session IDs)
- Event ingestion with batching, retries, and idempotency
- Incident management (create, assign, annotate, resolve, close)
- Human-in-the-loop approvals (request, approve, reject, escalate)
- Forensic replay (root cause analysis, causal graphs, counterfactuals, trace comparison)
- Compliance reports (SOC2, GDPR, custom frameworks)
- Cryptographic verification (Ed25519 signatures, Merkle proofs, receipt verification)
- Team management (members, roles, notifications, webhooks)
- Trace and session queries with filtering
- Flexible configuration via environment variables or programmatic config
- Custom HTTP client support
- Request/response interceptors
- Middleware hooks for event enrichment
- Event signing with Ed25519

### Security
- All API communications over HTTPS
- API key and OAuth token authentication
- Request signing for tamper-proof audit trails

## [0.9.0] - 2023-12-01

### Added
- Beta release for early adopters
- Core decision and event recording APIs
- Basic incident management
- Initial verification endpoints

### Known Issues
- Limited retry configuration options
- No batch event flushing control

## [0.8.0] - 2023-10-15

### Added
- Alpha release for internal testing
- Basic client initialization
- Simple event ingestion
- Minimal decision recording

---

## Versioning Scheme

This project uses [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Release Notes Format

Each release includes:
- **Added** - New features
- **Changed** - Changes in existing functionality
- **Deprecated** - Soon-to-be removed features
- **Removed** - Removed features
- **Fixed** - Bug fixes
- **Security** - Vulnerability fixes

## Links

- [GitHub Releases](https://github.com/blockloglabs/blocklog-go-sdk/releases)
- [Documentation](https://docs.blocklogsecurity.com)
- [Issues](https://github.com/blockloglabs/blocklog-go-sdk/issues)