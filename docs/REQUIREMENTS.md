# ArbiterID Project Requirements Document

This document records all functional and non-functional requirements for the ArbiterID project.

## Functional Requirements

### FR1. Core ID Generation
- **FR1.1** [✅ Complete] Generate 63-bit positive integer unique identifiers
- **FR1.2** [✅ Complete] Support custom bit allocation: type(10 bits) + timestamp(41 bits) + node ID(2 bits) + sequence(10 bits)
- **FR1.3** [✅ Complete] Support ID type range 0-1023 (10 bits)
- **FR1.4** [✅ Complete] Support node ID range 0-3 (2 bits)
- **FR1.5** [✅ Complete] Support sequence range 0-1023 within same millisecond (10 bits)
- **FR1.6** [✅ Complete] Timestamp based on custom epoch (2025-01-01T00:00:00.000Z)

### FR2. Encoding and Parsing
- **FR2.1** [✅ Complete] Support Base58 encoding/decoding
- **FR2.2** [✅ Complete] Support Base64 encoding/decoding
- **FR2.3** [✅ Complete] Support Base32 encoding/decoding
- **FR2.4** [✅ Complete] Support binary string encoding/decoding
- **FR2.5** [✅ Complete] Support decimal string encoding/decoding
- **FR2.6** [✅ Complete] Provide component extraction methods (type, timestamp, node, sequence)

### FR3. JSON Support
- **FR3.1** [✅ Complete] Implement JSON serialization/deserialization
- **FR3.2** [✅ Complete] Support string format JSON representation
- **FR3.3** [✅ Complete] Support numeric format JSON representation (backward compatibility)

### FR4. HTTP Service
- **FR4.1** [✅ Complete] Provide RESTful API
- **FR4.2** [✅ Complete] POST /generate endpoint for ID generation
- **FR4.3** [✅ Complete] Support batch generation (1-100 IDs)
- **FR4.4** [✅ Complete] Support specifying ID type
- **FR4.5** [✅ Complete] GET /health health check endpoint
- **FR4.6** [✅ Complete] GET /info service information endpoint
- **FR4.7** [✅ Complete] Multi-format response (Base58, Base64, hex, int64)

### FR5. Configuration and Options
- **FR5.1** [✅ Complete] Support strict monotonicity check toggle
- **FR5.2** [✅ Complete] Support quiet mode (production environment)
- **FR5.3** [✅ Complete] Configure HTTP service via environment variables (NODE_ID, PORT)

## Non-Functional Requirements

### NFR1. Performance Requirements
- **NFR1.1** [✅ Complete] Single-threaded generation performance: ~976 ns/op
- **NFR1.2** [✅ Complete] Concurrent generation capability: >1M IDs/second
- **NFR1.3** [✅ Complete] HTTP service response time: <1ms
- **NFR1.4** [✅ Complete] Base58 encoding performance: ~11 ns/op
- **NFR1.5** [✅ Complete] Base64 encoding performance: ~31 ns/op

### NFR2. Concurrent Safety
- **NFR2.1** [✅ Complete] Thread-safe ID generation
- **NFR2.2** [✅ Complete] Use mutex locks to protect Node operations
- **NFR2.3** [✅ Complete] Verify through race condition testing

### NFR3. Clock Resilience
- **NFR3.1** [✅ Complete] Handle clock drift and rollback
- **NFR3.2** [✅ Complete] Time waiting mechanism when sequence is exhausted
- **NFR3.3** [✅ Complete] Warning logs for clock anomalies
- **NFR3.4** [✅ Complete] Configure clock rollback wait parameters

### NFR4. Error Handling
- **NFR4.1** [✅ Complete] Define specific error types
- **NFR4.2** [✅ Complete] Node ID range validation
- **NFR4.3** [✅ Complete] ID type range validation
- **NFR4.4** [✅ Complete] Monotonicity violation detection
- **NFR4.5** [✅ Complete] Clock anomaly detection

### NFR5. Maintainability
- **NFR5.1** [✅ Complete] Code test coverage >95%
- **NFR5.2** [✅ Complete] Performance benchmark testing
- **NFR5.3** [✅ Complete] Code quality checks (golangci-lint)
- **NFR5.4** [✅ Complete] Automated CI/CD pipeline

### NFR6. Deployment and Operations
- **NFR6.1** [✅ Complete] Docker containerization support
- **NFR6.2** [✅ Complete] Multi-instance deployment configuration
- **NFR6.3** [✅ Complete] Load balancing configuration (Nginx)
- **NFR6.4** [✅ Complete] Health check mechanism
- **NFR6.5** [✅ Complete] Service discovery support

### NFR7. Compatibility
- **NFR7.1** [✅ Complete] Go 1.19+ compatibility
- **NFR7.2** [✅ Complete] Cross-platform support
- **NFR7.3** [✅ Complete] Standard HTTP protocol compatibility

## Quality Requirements

### QR1. Reliability
- **QR1.1** [✅ Complete] ID uniqueness guarantee
- **QR1.2** [✅ Complete] Uniqueness in distributed environments
- **QR1.3** [✅ Complete] Time ordering characteristics (K-sortable)

### QR2. Availability
- **QR2.1** [✅ Complete] High availability deployment patterns
- **QR2.2** [✅ Complete] Failover support
- **QR2.3** [✅ Complete] Graceful error handling

### QR3. Scalability
- **QR3.1** [✅ Complete] Horizontal scaling support (maximum 4 nodes)
- **QR3.2** [✅ Complete] Microservice architecture friendly
- **QR3.3** [✅ Complete] API version compatibility considerations

## Security Requirements

### SR1. Data Security
- **SR1.1** [✅ Complete] No sensitive information disclosure
- **SR1.2** [✅ Complete] Secure Docker container configuration (non-root user)

### SR2. Service Security
- **SR2.1** [✅ Complete] Input validation and range checking
- **SR2.2** [✅ Complete] Error messages do not leak internal state

## Constraints

### C1. Technical Constraints
- **C1.1** Implemented using Go language
- **C1.2** 63-bit positive integer limitation
- **C1.3** Maximum support for 4 nodes (2-bit node ID)
- **C1.4** Timestamp precision at millisecond level
- **C1.5** Time range approximately 69 years (41-bit timestamp)

### C2. Business Constraints
- **C2.1** IDs must maintain incremental characteristics
- **C2.2** IDs generated by the same node must be unique
- **C2.3** IDs generated by different nodes must be unique

## Future Requirements

### Future.1 Enhanced Features
- **F1.1** [TODO] Support more encoding formats
- **F1.2** [TODO] Support custom time precision
- **F1.3** [TODO] Support more node counts

### Future.2 Operations Enhancement
- **F2.1** [TODO] Monitoring metrics export (Prometheus)
- **F2.2** [TODO] Distributed tracing support
- **F2.3** [TODO] Configuration hot reload

### Future.3 Performance Optimization
- **F3.1** [TODO] Memory pool optimization
- **F3.2** [TODO] Lock-free concurrency optimization
- **F3.3** [TODO] Batch generation optimization

---

**Document Version**: 1.0
**Last Updated**: 2025-01-13
**Status Legend**: ✅ Complete, 🟡 In Progress, ❌ Not Started, [TODO] Planned