# Contributing to ArbiterID

Thank you for your interest in contributing to ArbiterID! This document provides guidelines for contributors.

## Code of Conduct

We expect all contributors to be respectful and considerate. Please maintain a professional and inclusive environment.

## How to Contribute

### Reporting Issues

- Check existing issues before creating new ones
- Use clear, descriptive titles
- Include:
  - Go version
  - Operating system
  - Minimal code example that reproduces the issue
  - Expected vs actual behavior

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test -v ./...`
6. Run linting: `golangci-lint run`
7. Commit with clear messages following [Conventional Commits](https://www.conventionalcommits.org/)
8. Push to your fork: `git push origin feature/your-feature-name`
9. Create a Pull Request

## Development Setup

```bash
# Clone the repository
git clone https://github.com/githonllc/arbiterid.git
cd arbiterid

# Install dependencies
go mod download

# Run tests
go test -v -race ./...

# Run linting
golangci-lint run

# Run benchmarks
go test -bench=. -benchmem
```

## Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Include both unit tests and integration tests
- Test edge cases and error conditions
- Use table-driven tests where appropriate

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with race detection
go test -v -race ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Use `goimports` to organize imports
- Follow the project's golangci-lint configuration
- Write clear, self-documenting code
- Add comments for exported functions and types

### Linting

We use golangci-lint with a custom configuration. Before submitting:

```bash
golangci-lint run
```

## Documentation

- Update README.md if adding new features
- Add godoc comments for exported functions/types
- Include examples in documentation
- Update CHANGELOG.md for user-facing changes

## Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat: add Base32 encoding support
fix: handle clock skew in sequence generation
docs: update README with new examples
test: add concurrency tests for ID generation
```

## Performance Considerations

- Profile code for performance-critical changes
- Include benchmark results for performance improvements
- Consider memory allocation patterns
- Test with realistic workloads

## Security

- Be mindful of potential security implications
- Avoid exposing sensitive information
- Consider timing attacks for cryptographic operations
- Report security issues privately to maintainers

## Release Process

Releases are managed by maintainers:

1. Version numbers follow [Semantic Versioning](https://semver.org/)
2. Changes are documented in CHANGELOG.md
3. Releases are tagged and include release notes
4. Go modules are published automatically

## Questions?

- Check existing issues and discussions
- Create an issue for questions about implementation
- Contact maintainers for security-related concerns

Thank you for contributing to ArbiterID!