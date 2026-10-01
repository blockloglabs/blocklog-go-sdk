# Contributing to Blocklog Go SDK

Thank you for your interest in contributing to the Blocklog Go SDK! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- A Blocklog account (for integration testing)

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/blocklog-go.git
   cd blocklog-go
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/blockloglabs/blocklog-go-sdk.git
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Run tests to verify setup:
   ```bash
   go test ./...
   ```

## Development Workflow

### Branching Strategy

- `main` - Stable release branch
- `develop` - Integration branch for ongoing development
- Feature branches: `feature/description-of-feature`
- Bug fix branches: `fix/description-of-bug`
- Release branches: `release/vX.Y.Z`

### Making Changes

1. Create a new branch from `develop`:
   ```bash
   git checkout develop
   git pull upstream develop
   git checkout -b feature/your-feature-name
   ```

2. Make your changes with clear, focused commits

3. Write tests for new functionality

4. Run the full test suite:
   ```bash
   go test ./...
   ```

5. Run linter:
   ```bash
   golangci-lint run
   ```

6. Push to your fork and create a Pull Request

## Code Style

### Go Code Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (run `gofmt -w .`)
- Follow the existing code style in the repository
- Add comments for exported functions and types
- Keep functions small and focused

### Error Handling

- Use the error types defined in `transport/` package
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check for specific error types using `errors.Is()`

### Testing

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies (HTTP servers, etc.)
- Aim for >80% code coverage
- Run tests with race detector: `go test -race ./...`

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): brief description

Longer description if needed

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add entry to CHANGELOG.md (under "Unreleased" section)
4. Request review from maintainers
5. Address review comments
6. Maintainers will merge after approval

## Reporting Issues

### Bug Reports

Include:
- Go version (`go version`)
- SDK version
- Steps to reproduce
- Expected vs actual behavior
- Error messages and stack traces
- Minimal code sample

### Feature Requests

Include:
- Use case description
- Proposed API design
- Any relevant examples
- Willingness to implement

## Release Process

Releases are managed by maintainers:

1. Create release branch from `develop`
2. Update version in relevant files
3. Update CHANGELOG.md
4. Tag release: `git tag vX.Y.Z`
5. Push tag: `git push origin vX.Y.Z`
6. GitHub Actions builds and publishes release

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you agree to uphold this code.

## Getting Help

- Check existing issues and PRs
- Read the [documentation](https://docs.blocklogsecurity.com)
- Join our [Discord community](https://discord.gg/blocklog)
- Email: support@blocklogsecurity.com

## License

By contributing, you agree that your contributions will be licensed under the MIT License.