# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ArbiterID is a Go library for generating 63-bit, K-sortable, unique identifiers inspired by Twitter's Snowflake algorithm. It generates distributed unique IDs with customizable bit allocation for type (10 bits), timestamp (41 bits), node ID (2 bits), and sequence (10 bits). The project includes both a core library and a production-ready HTTP service for microservice architectures.

## Common Development Commands

### Testing
```bash
make test               # Run standard tests
make test-race          # Run tests with race detection
make test-coverage      # Generate HTML coverage report
make benchmark          # Run performance benchmarks
```

### Code Quality
```bash
make lint              # Run golangci-lint
make fmt               # Format code with gofmt and goimports
make vet               # Run go vet
make dev               # Run fmt + vet + test (development workflow)
```

### Build and Examples
```bash
make build             # Build the project
make examples          # Run the simple example
make ci                # Full CI pipeline (deps + vet + lint + test-race)
```

### HTTP Service Development
```bash
cd examples/service
go build -o arbiter-id-service    # Build service
NODE_ID=0 ./arbiter-id-service    # Run service
./test-api.sh                     # Test API endpoints
docker build -t arbiter-id-service .  # Build Docker image
docker-compose up                 # Multi-instance deployment
```

### Dependency Management
```bash
make deps              # Download and verify dependencies
make tidy              # Tidy go.mod
```

## Architecture

### Core Library Components

- **ID Structure**: 63-bit positive int64 with custom bit allocation
- **Node**: Thread-safe ID generator with configurable options
- **Encoding**: Multiple encoding formats (Base58, Base64, Base32, Binary, Decimal)
- **JSON Support**: Custom marshaling/unmarshaling preserving string format

### HTTP Service Components

- **REST API**: Production-ready HTTP endpoints for ID generation
- **Health Monitoring**: Health checks and service information endpoints
- **Docker Support**: Complete containerization with multi-instance deployment
- **Load Balancing**: Nginx configuration for high availability

### Key Files

#### Core Library
- `arbiterid.go`: Main implementation with ID generation, encoding, and parsing
- `arbiterid_test.go`: Comprehensive test suite with >95% coverage

#### HTTP Service (`examples/service/`)
- `main.go`: Complete HTTP server with RESTful API
- `Dockerfile`: Multi-stage container build
- `docker-compose.yml`: Multi-instance deployment configuration
- `nginx.conf`: Load balancer setup
- `test-api.sh`: API testing script
- `README.md`: Service documentation

#### Examples and Documentation
- `examples/simple/main.go`: Basic library usage examples
- `README.md`: Main project documentation
- `PROJECT_STRUCTURE.md`: Detailed project structure
- `CLAUDE.md`: This development guidance file

### Design Principles

- **Thread Safety**: All Node operations are mutex-protected
- **Clock Resilience**: Handles clock drift and rollover scenarios
- **Strict Monotonicity**: Optional enforcement that IDs are strictly increasing
- **Performance**: Optimized for high-throughput generation (>1M IDs/sec)
- **Production Ready**: Quiet mode, comprehensive error handling, monitoring
- **Microservice Friendly**: Standalone HTTP service with Docker support

### ID Bit Layout (63 bits total)
```
[Type: 10 bits][Timestamp: 41 bits][Node: 2 bits][Sequence: 10 bits]
```

- Type: 0-1023 (application-defined categorization)
- Timestamp: Milliseconds since epoch (Jan 1, 2025)
- Node: 0-3 (must be unique per instance)
- Sequence: 0-1023 (per-millisecond counter)

### Configuration Options

#### Node Configuration
```go
node, err := arbiterid.NewNode(nodeID, options...)
```

- `WithStrictMonotonicityCheck(bool)`: Enable/disable monotonicity validation
- `WithQuietMode(bool)`: Suppress logging for production environments

#### HTTP Service Configuration
Environment variables:
- `NODE_ID`: 0-3 (must be unique per instance)
- `PORT`: HTTP server port (default: 8080)

### API Endpoints

#### POST `/generate`
Generate single or multiple IDs with optional type specification.

**Request:**
```json
{
  "id_type": 1,    // Optional: 0-1023, defaults to 0
  "count": 5       // Optional: 1-100, defaults to 1
}
```

**Response:**
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

#### GET `/health`
Health check endpoint for monitoring.

#### GET `/info`
Service information and configuration details.

### Error Handling

The library defines specific error types for different failure modes:
- `ErrInvalidNodeID`: Node ID not in range 0-3
- `ErrInvalIDType`: Type not in range 0-1023
- `ErrClockNotAdvancing`: System clock issues during rollover
- `ErrMonotonicityViolation`: New ID not greater than previous (when strict checks enabled)

### Performance Characteristics

- **Library Generation**: ~976 ns/op (single-threaded)
- **Concurrent Generation**: >1M IDs/second
- **HTTP Service**: <1ms response time
- **String Encoding**: ~27 ns/op
- **Base58 Encoding**: ~11 ns/op
- **Base64 Encoding**: ~31 ns/op

### Testing Strategy

#### Core Library Tests
- Unit tests for all public APIs
- Race condition testing for concurrent generation
- Benchmark tests for performance validation
- Component extraction and encoding round-trip tests
- Edge case testing (clock drift, sequence rollover, etc.)
- Memory usage testing with quiet mode

#### HTTP Service Tests
- API endpoint testing
- Request/response validation
- Error handling verification
- Health check validation
- Docker container testing
- Load balancing validation

### Deployment Patterns

#### Single Application Integration
```go
// Initialize once per application
node, err := arbiterid.NewNode(nodeID, arbiterid.WithQuietMode(true))
if err != nil {
    log.Fatal(err)
}

// Use throughout application
id, err := node.Generate(myIDType)
```

#### Microservice Architecture
```bash
# Deploy HTTP service
cd examples/service
docker build -t arbiter-id-service .
docker run -p 8080:8080 -e NODE_ID=0 arbiter-id-service

# Client integration
curl -X POST http://id-service:8080/generate -d '{"id_type": 1}'
```

#### High Availability Setup
```bash
# Multi-instance deployment
docker-compose up

# Manual scaling
NODE_ID=0 PORT=8081 ./arbiter-id-service &
NODE_ID=1 PORT=8082 ./arbiter-id-service &
NODE_ID=2 PORT=8083 ./arbiter-id-service &
NODE_ID=3 PORT=8084 ./arbiter-id-service &
```

### Development Workflow

#### Local Development
1. Make changes to code
2. Run tests: `make test`
3. Check code quality: `make lint`
4. Test examples: `make examples`
5. Test HTTP service: `cd examples/service && ./test-api.sh`

#### Service Development
1. Modify service code in `examples/service/main.go`
2. Build: `go build -o arbiter-id-service`
3. Test locally: `NODE_ID=0 ./arbiter-id-service`
4. Test API: `./test-api.sh`
5. Build Docker: `docker build -t arbiter-id-service .`
6. Test deployment: `docker-compose up`

#### Integration Testing
1. Start service: `docker-compose up -d`
2. Run integration tests
3. Check health: `curl http://localhost/health`
4. Load test if needed
5. Clean up: `docker-compose down`

### Monitoring and Observability

#### Health Checks
- `/health` endpoint for service status
- Last generated ID tracking
- Node health monitoring
- Clock synchronization status

#### Metrics Collection
- ID generation rate
- Error rates and types
- Response time distribution
- Memory and CPU usage
- Container health status

#### Logging
- Structured JSON logs (recommended)
- Configurable log levels
- Error tracking and alerting
- Performance metrics logging
- Quiet mode for high-volume environments

### Security Considerations

#### Library Security
- No external dependencies
- Secure ID generation
- No sensitive data in IDs
- Input validation

#### Service Security
- Non-root container execution
- Input sanitization and validation
- Rate limiting (recommended)
- HTTPS deployment (recommended)
- Network security best practices

### Code Style and Quality

#### Go Style Guidelines
- Follow standard Go conventions
- Use golangci-lint configuration in `.golangci.yml`
- Line length limit: 120 characters
- Comprehensive error messages with context
- Prefer explicit error handling

#### Documentation Standards
- Godoc comments for all public APIs
- README updates for new features
- API documentation for HTTP service
- Example code for complex features

#### Testing Standards
- >95% test coverage requirement
- Benchmark tests for performance-critical code
- Race condition testing for concurrent code
- Integration tests for HTTP service
- Error path testing

### Common Development Tasks

#### Adding New Encoding Format
1. Add encoding method to ID type
2. Add parsing function
3. Add tests for round-trip conversion
4. Add benchmarks
5. Update documentation

#### Modifying HTTP Service
1. Update `examples/service/main.go`
2. Update API tests in `test-api.sh`
3. Update service documentation
4. Test Docker deployment
5. Update nginx config if needed

#### Performance Optimization
1. Add benchmarks for target functionality
2. Profile with `go tool pprof`
3. Optimize hot paths
4. Verify no regression in existing benchmarks
5. Update performance documentation

#### Error Handling Enhancement
1. Define new error types if needed
2. Add comprehensive error messages
3. Update error handling in HTTP service
4. Add tests for error conditions
5. Update API documentation