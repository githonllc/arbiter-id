# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-05-24

### Added

#### Core Library
- 63-bit Snowflake-inspired ID generation with custom bit allocation
- Support for 4 node IDs (0-3) in distributed systems
- Type field with 10 bits (0-1023) for flexible ID categorization
- Thread-safe concurrent ID generation with mutex protection
- Clock drift resilience and sequence rollover protection
- Configurable strict monotonicity checks with `WithStrictMonotonicityCheck()`
- Quiet mode support with `WithQuietMode()` for production environments
- Multiple encoding formats: decimal, Base2, Base32, Base58, Base64
- JSON marshaling/unmarshaling support preserving string format
- Component extraction methods for debugging and analysis
- Comprehensive error handling for clock anomalies and configuration issues

#### HTTP Service
- Production-ready RESTful HTTP API service (`examples/service/`)
- POST `/generate` endpoint with batch ID generation (1-100 IDs)
- GET `/health` endpoint for service monitoring and health checks
- GET `/info` endpoint for service metadata and configuration
- Support for custom ID types via JSON or query parameters
- Multiple response formats: Base58, Base64, Hex, Int64
- Environment-based configuration (NODE_ID, PORT)
- Docker containerization with multi-stage builds
- Docker Compose deployment with 4-instance setup
- Nginx load balancer configuration with health checks
- API testing script (`test-api.sh`) for comprehensive validation

#### Development & Infrastructure
- Comprehensive test suite with >95% coverage (46 test cases)
- Benchmark tests for performance validation and regression detection
- Race condition testing for concurrent scenarios
- Memory usage testing with quiet mode validation
- GitHub Actions CI/CD workflows for automated testing and releases
- golangci-lint configuration with modern Go linting rules
- Makefile for build automation and development workflow
- Complete documentation suite (README, PROJECT_STRUCTURE, CLAUDE guides)
- Usage examples in `examples/simple/` and `examples/service/`

### Technical Specifications

#### ID Structure (63 bits total)
- **Type**: 10 bits (0-1023) - Application-defined ID categories
- **Timestamp**: 41 bits - Milliseconds since epoch (2025-01-01T00:00:00Z)
- **Node ID**: 2 bits (0-3) - Identifies generating server instance
- **Sequence**: 10 bits (0-1023) - Per-millisecond counter per node
- **Total**: 63 bits (guaranteed positive int64)

#### Performance Characteristics
- High-performance ID generation: >1M IDs/second per node
- Single-threaded generation: ~976 ns/op
- String encoding: ~27 ns/op
- Base58 encoding: ~11 ns/op
- Base64 encoding: ~31 ns/op
- HTTP service response time: <1ms
- Memory-efficient with minimal allocations

#### Encoding Support
- **Decimal**: Human-readable numeric string (default)
- **Base58**: Compact, URL-safe encoding (Bitcoin alphabet)
- **Base64**: Standard URL-safe encoding without padding
- **Base32**: Case-insensitive custom alphabet
- **Binary**: Full 63-bit binary representation
- **Hexadecimal**: Lowercase hex representation

#### Production Features
- Docker containerization with security best practices
- Multi-instance deployment with unique node IDs
- Load balancing configuration for high availability
- Health monitoring and service discovery support
- Comprehensive error handling and logging
- Clock synchronization monitoring
- Input validation and sanitization

### Dependencies
- No external runtime dependencies for core library
- Go 1.21+ required for development
- Docker for containerized deployment
- Nginx for load balancing (optional)

### Breaking Changes
- None (initial release)

### Migration Notes
- This is the initial stable release
- API is considered stable and will follow semantic versioning
- Future breaking changes will increment major version

---

## Release Notes Format

Each release includes:
- **Added**: New features and capabilities
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Now removed features
- **Fixed**: Bug fixes and issue resolutions
- **Security**: Vulnerability fixes and security improvements

## Migration Guides

### From Snowflake Libraries
ArbiterID is designed as a drop-in replacement for most Snowflake implementations with enhanced features:

#### Key Differences
- **Enhanced Type Field**: 10-bit type field (0-1023) instead of typical 8-bit machine/datacenter IDs
- **Modern Epoch**: Custom epoch starts from 2025-01-01 for maximum future range (~69 years)
- **Production Ready**: Built-in quiet mode, comprehensive error handling, and monitoring
- **Multiple Encodings**: Base58, Base64, Base32 support out of the box
- **HTTP Service**: Optional standalone service for microservice architectures
- **Clock Protection**: Advanced clock drift and rollback protection mechanisms

#### Migration Steps
1. Replace Snowflake library import with `github.com/githonllc/arbiterid`
2. Update node initialization: `arbiterid.NewNode(nodeID)`
3. Replace ID generation calls: `node.Generate(idType)`
4. Update ID parsing for new encoding formats if needed
5. Configure quiet mode for production: `WithQuietMode(true)`

### API Stability Promise
Starting from v1.0.0, the API follows strict semantic versioning:
- **Major versions (v2.x.x)**: Breaking API changes requiring code updates
- **Minor versions (v1.x.x)**: Backward-compatible feature additions
- **Patch versions (v1.0.x)**: Bug fixes and performance improvements only

### Deployment Patterns

#### Single Application Integration
```go
node, err := arbiterid.NewNode(nodeID, arbiterid.WithQuietMode(true))
id, err := node.Generate(myIDType)
```

#### Microservice Deployment
```bash
# Deploy HTTP service
docker run -p 8080:8080 -e NODE_ID=0 arbiter-id-service

# Multi-instance with load balancing
docker-compose up
```

#### High Availability Setup
- Deploy 4 instances with NODE_ID 0-3
- Use load balancer (nginx, HAProxy, etc.)
- Monitor via `/health` endpoints
- Scale horizontally within node limits