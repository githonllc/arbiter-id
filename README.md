# Arbiter ID

[![Go Reference](https://pkg.go.dev/badge/github.com/githonllc/arbiterid.svg)](https://pkg.go.dev/github.com/githonllc/arbiterid)
[![Go Report Card](https://goreportcard.com/badge/github.com/githonllc/arbiterid)](https://goreportcard.com/report/github.com/githonllc/arbiterid)
[![Build Status](https://github.com/githonllc/arbiterid/actions/workflows/go.yml/badge.svg)](https://github.com/githonllc/arbiterid/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`arbiterid` is a Go library for generating 63-bit, K-sortable, unique identifiers inspired by Twitter's Snowflake. These IDs are designed to be roughly time-ordered and can embed custom information like ID type and the generating node ID. It's built for distributed systems requiring unique IDs with up to 4 generating nodes, high concurrency, and resilience against clock drifts.

## Key Features

*   **Snowflake-like Structure:** 63-bit positive `int64` IDs.
*   **Customizable Bit Allocation:**
    *   **Type:** 10 bits (for categorizing IDs, 0-1023).
    *   **Timestamp:** 41 bits (milliseconds since custom epoch, ~69 years lifespan).
    *   **Node ID:** 2 bits (for differentiating generating servers, 0-3).
    *   **Sequence:** 10 bits (for IDs generated in the same millisecond on the same node, 0-1023).
*   **High Concurrency Support:** Thread-safe generation within a single node instance.
*   **Clock Drift Resilience:** Handles minor clock drifts and protects against clock stalls during sequence rollovers.
*   **Quiet Mode:** Optional suppression of logging output for high-volume production environments.
*   **Multiple Encodings:** Supports decimal string, Base2, Base32 (custom alphabet), Base58, and efficient Base64 (URL-safe) representations.
*   **JSON Marshalling:** Marshals IDs as strings in JSON to preserve precision.
*   **Component Extraction:** Easily extract type, timestamp, node, and sequence from an ID.
*   **HTTP Service:** Production-ready standalone HTTP API service for distributed deployments.

## Installation

```bash
go get github.com/githonllc/arbiterid
```

## Quick Start

### Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/githonllc/arbiterid"
)

func main() {
    // Initialize a new node with quiet mode for production
    node, err := arbiterid.NewNode(0, arbiterid.WithQuietMode(true))
    if err != nil {
        log.Fatalf("Failed to create arbiterid node: %v", err)
    }

    // Define ID types (0-1023)
    const UserIDType arbiterid.IDType = 1
    const PostIDType arbiterid.IDType = 512

    // Generate an ID
    userID, err := node.Generate(UserIDType)
    if err != nil {
        log.Fatalf("Failed to generate UserID: %v", err)
    }
    fmt.Printf("Generated User ID (int64): %d\n", userID)
    fmt.Printf("Generated User ID (string): %s\n", userID.String())
    fmt.Printf("Generated User ID (Base58): %s\n", userID.Base58())

    postID, err := node.Generate(PostIDType)
    if err != nil {
        log.Fatalf("Failed to generate PostID: %v", err)
    }
    fmt.Printf("Generated Post ID (Base64): %s\n", postID.Base64())

    // Extract components from an ID
    IDType, tsMillis, nodeID, seq := userID.Components()
    fmt.Printf("UserID Components: Type=%d, TimestampMillis=%d, Node=%d, Seq=%d\n",
        IDType, tsMillis, nodeID, seq)
    fmt.Printf("UserID Timestamp (ISO): %s\n", userID.TimeISO())
}
```

### HTTP Service

For microservice architectures, use the standalone HTTP service:

```bash
# Build and run the service
cd examples/service
go build -o arbiter-id-service
NODE_ID=0 PORT=8080 ./arbiter-id-service
```

#### API Examples

```bash
# Generate a single ID
curl -X POST http://localhost:8080/generate

# Generate specific type ID
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"id_type": 1}'

# Batch generate IDs
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"id_type": 1, "count": 5}'

# Health check
curl http://localhost:8080/health

# Service information
curl http://localhost:8080/info
```

#### Docker Deployment

```bash
# Single instance
docker build -t arbiter-id-service examples/service/
docker run -p 8080:8080 -e NODE_ID=0 arbiter-id-service

# Multi-instance with load balancer
cd examples/service
docker-compose up
```

## ID Structure (63 bits)

The ID is a 63-bit integer, ensuring it's always positive when stored as an `int64`. The most significant bit is unused (0).

| Component   | Bits | Description                                         | Max Value  |
| :---------- | :--- | :-------------------------------------------------- | :--------- |
| Type        | 10   | Custom application-defined ID type                  | 1023       |
| Timestamp   | 41   | Milliseconds since custom epoch (Jan 1, 2025 UTC)  | ~69 years  |
| Node ID     | 2    | Identifier for the generating node/server           | 3          |
| Sequence    | 10   | Per-millisecond, per-node sequence number           | 1023       |

**Total: 10 + 41 + 2 + 10 = 63 bits**

## Configuration

### Node Options

```go
// Create node with options
node, err := arbiterid.NewNode(nodeID, options...)
```

Available options:

*   `WithStrictMonotonicityCheck(enable bool)`: (Default: `true`) Enables/disables checking that every new ID is strictly greater than the last one.
*   `WithQuietMode(enable bool)`: (Default: `false`) Suppresses most log output for production environments.

### HTTP Service Configuration

Environment variables:
*   `NODE_ID`: Node identifier (0-3, must be unique per instance)
*   `PORT`: HTTP server port (default: 8080)

## Error Handling

The library handles various error conditions:

*   `ErrInvalidNodeID`, `ErrInvalIDType`: Configuration errors.
*   `ErrClockNotAdvancing`: System clock issues during sequence rollover.
*   `ErrMonotonicityViolation`: New ID not greater than previous (when strict checks enabled).
*   Timestamp overflow: Current time exceeds 41-bit limit (~69 years from epoch).

## Encoding and Decoding

Multiple representation formats:

*   `ID.String() string`: Decimal string.
*   `ID.Int64() int64`: Raw `int64` value.
*   `ID.Base2() string`: Binary string.
*   `ID.Base32() string`: Custom Base32 encoded string.
*   `ID.Base58() string`: Base58 encoded string (Bitcoin alphabet).
*   `ID.Base64() string`: URL-safe Base64 encoded string (no padding).

Corresponding parsing functions:

*   `ParseString(s string) (ID, error)`
*   `ParseBase2(s string) (ID, error)`
*   `ParseBase32(s string) (ID, error)`
*   `ParseBase58(s string) (ID, error)`
*   `ParseBase64(s string) (ID, error)`

## Performance

Benchmark results on modern hardware:
- **Generation**: ~976 ns/op (single-threaded)
- **Concurrent generation**: >1M IDs/second
- **String encoding**: ~27 ns/op
- **Base58 encoding**: ~11 ns/op
- **Base64 encoding**: ~31 ns/op

## Production Deployment

### Single Application Integration

```go
// Initialize once per application instance
node, err := arbiterid.NewNode(nodeID, arbiterid.WithQuietMode(true))
if err != nil {
    log.Fatal(err)
}

// Use throughout application
id, err := node.Generate(myIDType)
```

### Microservice Architecture

Deploy the HTTP service and integrate via API:

```go
// Client example
func generateUserID() (string, error) {
    resp, err := http.Post("http://id-service:8080/generate",
        "application/json",
        strings.NewReader(`{"id_type": 1}`))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Success bool `json:"success"`
        Data struct {
            ID string `json:"id"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    if !result.Success {
        return "", errors.New("ID generation failed")
    }

    return result.Data.ID, nil
}
```

### High Availability Setup

```bash
# Deploy multiple instances with unique node IDs
NODE_ID=0 PORT=8081 ./arbiter-id-service &
NODE_ID=1 PORT=8082 ./arbiter-id-service &
NODE_ID=2 PORT=8083 ./arbiter-id-service &
NODE_ID=3 PORT=8084 ./arbiter-id-service &

# Use load balancer (nginx, HAProxy, etc.)
```

## Project Structure

```
arbiterid/
├── arbiterid.go              # Core library implementation
├── arbiterid_test.go         # Comprehensive test suite
├── examples/
│   ├── simple/               # Basic usage examples
│   └── service/              # Production HTTP service
│       ├── main.go           # HTTP server implementation
│       ├── Dockerfile        # Container build
│       ├── docker-compose.yml # Multi-instance deployment
│       ├── nginx.conf        # Load balancer configuration
│       ├── test-api.sh       # API testing script
│       └── README.md         # Service documentation
├── README.md                 # This file
├── PROJECT_STRUCTURE.md      # Detailed project structure
├── CLAUDE.md                 # AI development guidance
└── go.mod                    # Go module definition
```

## Concurrency

All components are safe for concurrent use:
- `Node` instances are thread-safe
- HTTP service handles concurrent requests
- No shared state between different node IDs

## Limitations & Considerations

*   **Node ID Uniqueness:** Each instance must have a unique `nodeID` (0-3).
*   **Timestamp Rollover:** 41-bit timestamp exhausts around year 2094.
*   **System Clock Dependency:** Requires reasonably accurate system clock.
*   **Maximum Throughput:** 1024 unique IDs per millisecond per node.
*   **Network Partitions:** HTTP service instances should be properly load balanced.

## Examples

The `examples/` directory contains:

- **`simple/`**: Basic library usage examples
- **`service/`**: Complete HTTP service with:
  - RESTful API endpoints
  - Docker containerization
  - Multi-instance deployment
  - Load balancing configuration
  - Health checks and monitoring
  - Comprehensive API documentation

## Contributing

Contributions welcome! Please:

1.  Fork the repository
2.  Create your feature branch
3.  Add tests for new functionality
4.  Ensure all tests pass: `make test`
5.  Run linter: `make lint`
6.  Submit a pull request

## License

This project is licensed under the [MIT License](LICENSE).