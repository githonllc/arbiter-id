# ArbiterID Project Structure

This document outlines the complete project structure for the ArbiterID Go library.

## Project Overview

ArbiterID is a high-performance, distributed unique ID generator library for Go, inspired by Twitter's Snowflake algorithm. It generates 63-bit K-sortable unique identifiers with the following structure:

- **Type**: 10 bits (0-1023) for ID categorization
- **Timestamp**: 41 bits (milliseconds since epoch)
- **Node ID**: 2 bits (0-3) for distributed systems
- **Sequence**: 10 bits (0-1023) for same-millisecond generation

## Directory Structure

```
arbiterid/
├── .github/
│   └── workflows/
│       ├── go.yml              # CI/CD workflow
│       └── release.yml         # Release workflow
├── examples/
│   ├── simple/
│   │   ├── go.mod              # Example module
│   │   └── main.go             # Basic usage examples
│   └── service/                # Production HTTP service
│       ├── main.go             # HTTP server implementation
│       ├── go.mod              # Service module
│       ├── go.sum              # Dependencies lock
│       ├── Dockerfile          # Container build configuration
│       ├── docker-compose.yml  # Multi-instance deployment
│       ├── nginx.conf          # Load balancer configuration
│       ├── test-api.sh         # API testing script
│       └── README.md           # Service documentation
├── .gitignore                  # Git ignore rules
├── .golangci.yml              # Linter configuration
├── CHANGELOG.md               # Version history
├── CLAUDE.md                  # AI development guidance
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # MIT License
├── Makefile                   # Build automation
├── PROJECT_STRUCTURE.md       # This file
├── README.md                  # Main documentation
├── arbiterid.go              # Main library code
├── arbiterid_test.go         # Comprehensive tests
└── go.mod                    # Go module definition
```

## Core Files

### Main Library
- **`arbiterid.go`**: Core implementation with ID generation, encoding, and parsing
- **`arbiterid_test.go`**: Comprehensive test suite with >95% coverage
- **`go.mod`**: Go module definition and dependencies

### Documentation
- **`README.md`**: Main documentation with usage examples and API reference
- **`CLAUDE.md`**: AI development guidance and architecture overview
- **`CHANGELOG.md`**: Version history and release notes
- **`CONTRIBUTING.md`**: Guidelines for contributors
- **`PROJECT_STRUCTURE.md`**: This structure overview

### Configuration
- **`.golangci.yml`**: Linter configuration for code quality
- **`.gitignore`**: Git ignore patterns
- **`Makefile`**: Build automation and development commands

### CI/CD
- **`.github/workflows/go.yml`**: Continuous integration workflow
- **`.github/workflows/release.yml`**: Automated release workflow

## Examples Directory

### Simple Example (`examples/simple/`)
- **`main.go`**: Comprehensive library usage examples
- **`go.mod`**: Example module configuration

### HTTP Service (`examples/service/`)

A production-ready HTTP service for distributed ID generation:

#### Core Files
- **`main.go`**: Complete HTTP server implementation with:
  - RESTful API endpoints (`/generate`, `/health`, `/info`)
  - JSON request/response handling
  - Environment variable configuration
  - Error handling and validation
  - Batch ID generation support

#### Deployment Files
- **`Dockerfile`**: Multi-stage Docker build with:
  - Alpine-based runtime image
  - Non-root user security
  - Health check configuration
  - Optimized for production

- **`docker-compose.yml`**: Multi-instance deployment with:
  - 4 service instances (Node IDs 0-3)
  - Nginx load balancer
  - Health checks
  - Network configuration

- **`nginx.conf`**: Load balancer configuration with:
  - Upstream server definitions
  - Health check endpoints
  - Connection pooling
  - Request routing
  - Status monitoring

#### Testing and Documentation
- **`test-api.sh`**: Comprehensive API testing script
- **`README.md`**: Complete service documentation with:
  - API endpoint descriptions
  - Usage examples
  - Deployment guides
  - Client integration examples
  - Performance characteristics
  - Troubleshooting guide

## Key Features

### Core Library Features
- Thread-safe concurrent generation
- Clock drift resilience with configurable warnings
- Sequence rollover protection
- Configurable strict monotonicity checks
- Quiet mode for production environments
- Multiple encoding formats (Base58, Base64, Base32, Binary, Decimal)
- JSON marshaling/unmarshaling support
- Component extraction methods

### HTTP Service Features
- High-performance HTTP API (>1M IDs/second)
- RESTful endpoint design
- Batch ID generation (1-100 IDs)
- Multiple ID formats in responses
- Health monitoring endpoints
- Service information and status
- Environment-based configuration
- Docker containerization
- Multi-instance deployment support
- Load balancing configuration

### Production Features
- Comprehensive error handling
- Detailed logging with quiet mode
- Performance monitoring
- Health checks
- Container security best practices
- High availability deployment patterns

## Development Workflow

### Available Make Commands
```bash
make help          # Show available commands
make build         # Build the project
make test          # Run tests
make test-race     # Run tests with race detection
make test-coverage # Generate coverage report
make benchmark     # Run performance benchmarks
make lint          # Run linter
make fmt           # Format code
make vet           # Run go vet
make examples      # Run examples
make clean         # Clean build artifacts
make ci            # Run CI pipeline
make dev           # Run development checks
make pre-release   # Prepare for release
```

### Testing Strategy
- Unit tests with comprehensive coverage (>95%)
- Race condition testing for concurrent scenarios
- Benchmark tests for performance validation
- Memory usage testing
- Concurrent generation testing
- Clock drift and rollover testing
- API integration testing (HTTP service)

### Code Quality
- golangci-lint for static analysis
- go vet for potential issues
- gofmt for consistent formatting
- Comprehensive test coverage requirements
- Performance regression testing

## HTTP Service Architecture

### API Endpoints

#### POST `/generate`
- Generate single or multiple IDs
- Support for custom ID types (0-1023)
- Batch generation (1-100 IDs)
- JSON and query parameter input
- Multiple output formats

#### GET `/health`
- Service health status
- Node information
- Last generated ID
- Timestamp information

#### GET `/info`
- Service metadata
- Configuration details
- API documentation
- Bit layout information

### Configuration
- **NODE_ID**: 0-3 (must be unique per instance)
- **PORT**: HTTP server port (default: 8080)

### Response Format
```json
{
  "success": true,
  "data": {
    "id": "base58_encoded_id",
    "id_int64": 1234567890123456789,
    "id_base64": "base64_encoded_id",
    "id_hex": "hex_encoded_id",
    "type": 1,
    "time": "2025-01-13T12:34:56.789Z",
    "node": 0,
    "sequence": 1
  }
}
```

## Deployment Patterns

### Single Instance
- Direct library integration
- Embedded in application
- Shared node instance

### Microservice Architecture
- Standalone HTTP service
- API-based integration
- Service discovery support
- Load balancing

### High Availability
- Multi-instance deployment
- Load balancer distribution
- Health check monitoring
- Failover support

## Performance Characteristics

Based on benchmark results:
- **Generation**: ~976 ns/op (single-threaded)
- **Concurrent generation**: >1M IDs/second
- **String encoding**: ~27 ns/op
- **Base58 encoding**: ~11 ns/op
- **Base64 encoding**: ~31 ns/op
- **Parsing**: 11-29 ns/op depending on format
- **HTTP service**: <1ms response time
- **Memory usage**: Minimal allocation during generation

## Security Considerations

### Library
- No external dependencies
- Secure random ID generation
- No sensitive information in IDs

### HTTP Service
- Non-root container execution
- Input validation and sanitization
- Rate limiting recommendations
- HTTPS deployment recommended
- Network security best practices

## Release Process

1. Update `CHANGELOG.md` with new version
2. Run full test suite: `make ci`
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions automatically creates release
6. Update Docker images if needed

## Usage Examples

### Library Integration
```go
node, err := arbiterid.NewNode(0, arbiterid.WithQuietMode(true))
id, err := node.Generate(1)
```

### HTTP Service Usage
```bash
# Start service
cd examples/service
go build -o arbiter-id-service
NODE_ID=0 ./arbiter-id-service

# Generate ID
curl -X POST http://localhost:8080/generate -d '{"id_type": 1}'
```

### Docker Deployment
```bash
# Single instance
docker build -t arbiter-id-service examples/service/
docker run -p 8080:8080 -e NODE_ID=0 arbiter-id-service

# Multi-instance
cd examples/service
docker-compose up
```

## Monitoring and Observability

### Metrics
- ID generation rate
- Error rates
- Response times
- Memory usage
- CPU utilization

### Health Checks
- Service status endpoint
- Database connectivity (if applicable)
- Clock synchronization status
- Node health status

### Logging
- Structured logging support
- Configurable log levels
- Error tracking
- Performance metrics

## License

MIT License - see `LICENSE` file for details.

## Contributing

See `CONTRIBUTING.md` for guidelines on how to contribute to this project.