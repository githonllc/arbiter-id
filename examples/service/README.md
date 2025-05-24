# ArbiterID Generation Service

This is a standalone HTTP service that provides distributed unique ID generation functionality for other applications. The service is based on a Snowflake-like algorithm and generates 63-bit K-sortable unique identifiers.

## Features

- 🚀 High-performance HTTP API
- 🔢 Support for custom ID types (0-1023)
- 📦 Batch ID generation
- 🏥 Health check endpoint
- 📊 Service information endpoint
- 🌐 RESTful API design
- 🔒 Thread-safe
- 📝 Detailed JSON responses

## Quick Start

### 1. Build and Run

```bash
# Enter service directory
cd examples/service

# Build service
go build -o arbiter-id-service

# Run service (default port 8080, node ID 0)
./arbiter-id-service
```

### 2. Environment Variable Configuration

```bash
# Set node ID (0-3, must be unique for each instance)
export NODE_ID=0

# Set port
export PORT=8080

# Run service
./arbiter-id-service
```

### 3. Docker Run

```bash
# Build Docker image
docker build -t arbiter-id-service .

# Run container
docker run -p 8080:8080 -e NODE_ID=0 arbiter-id-service
```

## API Endpoints

### POST /generate

Generate one or more unique IDs.

#### Request Parameters

**JSON Body Parameters:**
```json
{
  "id_type": 1,    // Optional: ID type (0-1023), defaults to 0
  "count": 5       // Optional: generation count (1-100), defaults to 1
}
```

**Query Parameters (Alternative):**
- `type`: ID type (0-1023)
- `count`: Generation count (1-100)

#### Response Examples

**Single ID:**
```json
{
  "success": true,
  "data": {
    "id": "9m4e2mr0ui3e8a215n4g",
    "id_int64": 1234567890123456789,
    "id_base64": "ESaRJ_eMbwk",
    "id_hex": "112a0439fc61b009",
    "type": 1,
    "time": "2025-01-13T12:34:56.789Z",
    "node": 0,
    "sequence": 1
  }
}
```

**Multiple IDs:**
```json
{
  "success": true,
  "data": [
    {
      "id": "9m4e2mr0ui3e8a215n4g",
      "id_int64": 1234567890123456789,
      "id_base64": "ESaRJ_eMbwk",
      "id_hex": "112a0439fc61b009",
      "type": 1,
      "time": "2025-01-13T12:34:56.789Z",
      "node": 0,
      "sequence": 1
    },
    {
      "id": "9m4e2mr0ui3e8a215n4h",
      "id_int64": 1234567890123456790,
      "id_base64": "ESaRJ_eMbwm",
      "id_hex": "112a0439fc61b00a",
      "type": 1,
      "time": "2025-01-13T12:34:56.789Z",
      "node": 0,
      "sequence": 2
    }
  ]
}
```

### GET /health

Health check endpoint to verify the service is running normally.

#### Response Example

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-01-13T12:34:56Z",
    "node_id": 0,
    "last_id": "1234567890123456789"
  }
}
```

### GET /info

Get service information and configuration.

#### Response Example

```json
{
  "success": true,
  "data": {
    "service": "ArbiterID Generation Service",
    "version": "1.0.0",
    "description": "Distributed unique ID generation service using Snowflake-inspired algorithm",
    "node_id": 0,
    "epoch": "2025-01-01T00:00:00.000Z",
    "bit_layout": {
      "type": "10 bits (0-1023)",
      "timestamp": "41 bits (milliseconds since epoch)",
      "node": "2 bits (0-3)",
      "sequence": "10 bits (0-1023)"
    },
    "endpoints": {
      "POST /generate": "Generate new ID(s)",
      "GET /health": "Health check",
      "GET /info": "Service information"
    }
  }
}
```

## Usage Examples

### 1. Generate Single Default ID

```bash
curl -X POST http://localhost:8080/generate
```

### 2. Generate Specific Type ID

```bash
# Using JSON
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"id_type": 42}'

# Using query parameters
curl -X POST http://localhost:8080/generate?type=42
```

### 3. Batch Generate IDs

```bash
# Generate 5 user IDs (type 1)
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{"id_type": 1, "count": 5}'
```

### 4. Health Check

```bash
curl http://localhost:8080/health
```

### 5. Get Service Information

```bash
curl http://localhost:8080/info
```

## ID Format Description

Each generated ID contains the following information:

- **ID**: Base58 encoded string (recommended for APIs)
- **ID Int64**: Raw 64-bit integer
- **ID Base64**: Base64 encoded string
- **ID Hex**: Hexadecimal representation
- **Type**: ID type (0-1023)
- **Time**: ISO formatted generation time
- **Node**: Node ID (0-3)
- **Sequence**: Sequence number (0-1023)

## Deployment Recommendations

### Single Instance Deployment

```bash
# Start service
NODE_ID=0 PORT=8080 ./arbiter-id-service
```

### Multi-Instance Deployment

For high availability, deploy multiple instances with different NODE_IDs:

```bash
# Instance 1
NODE_ID=0 PORT=8081 ./arbiter-id-service &

# Instance 2
NODE_ID=1 PORT=8082 ./arbiter-id-service &

# Instance 3
NODE_ID=2 PORT=8083 ./arbiter-id-service &

# Instance 4
NODE_ID=3 PORT=8084 ./arbiter-id-service &
```

### Load Balancing

Use nginx or other load balancers to distribute requests:

```nginx
upstream arbiter_id_backend {
    server localhost:8081;
    server localhost:8082;
    server localhost:8083;
    server localhost:8084;
}

server {
    listen 80;

    location / {
        proxy_pass http://arbiter_id_backend;
        proxy_set_header Host $host;
    }
}
```

## Performance Characteristics

- **High Throughput**: Each instance can generate 1M+ IDs/second
- **Low Latency**: Single ID generation typically <1ms
- **K-sortable**: IDs sorted by generation time
- **No Duplicates**: Guaranteed globally unique under given configuration
- **Clock Safe**: Handles clock drift and rollback

## Error Handling

The service returns standard HTTP status codes and JSON error responses:

```json
{
  "success": false,
  "error": "ID type must be between 0 and 1023"
}
```

Common errors:
- `400 Bad Request`: Invalid parameters
- `405 Method Not Allowed`: Unsupported HTTP method
- `500 Internal Server Error`: Service internal error

## Monitoring and Logging

The service displays basic information on startup:

```
Starting ArbiterID service on port 8080
Node ID: 0
Available endpoints:
  POST /generate - Generate new ID(s)
  GET  /health   - Health check
  GET  /info     - Service information
  GET  /         - Service information
```

Recommended monitoring metrics:
- `/health` endpoint response time and status
- ID generation QPS
- Error rate
- Memory and CPU usage

## Client Integration

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type GenerateRequest struct {
    IDType *int `json:"id_type,omitempty"`
    Count  *int `json:"count,omitempty"`
}

type IDData struct {
    ID      string `json:"id"`
    IDInt64 int64  `json:"id_int64"`
    Type    int    `json:"type"`
}

type GenerateResponse struct {
    Success bool    `json:"success"`
    Data    IDData  `json:"data"`
    Error   string  `json:"error,omitempty"`
}

func generateID(baseURL string, idType int) (*IDData, error) {
    req := GenerateRequest{IDType: &idType}
    body, _ := json.Marshal(req)

    resp, err := http.Post(baseURL+"/generate", "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result GenerateResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if !result.Success {
        return nil, fmt.Errorf("API error: %s", result.Error)
    }

    return &result.Data, nil
}
```

### Python Client Example

```python
import requests

def generate_id(base_url, id_type=0, count=1):
    """Generate ID"""
    payload = {"id_type": id_type, "count": count}
    response = requests.post(f"{base_url}/generate", json=payload)

    if response.status_code == 200:
        data = response.json()
        if data["success"]:
            return data["data"]

    raise Exception(f"Failed to generate ID: {response.text}")

# Usage example
id_data = generate_id("http://localhost:8080", id_type=1)
print(f"Generated ID: {id_data['id']}")
```

### JavaScript/Node.js Client Example

```javascript
async function generateID(baseURL, idType = 0, count = 1) {
    const response = await fetch(`${baseURL}/generate`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id_type: idType, count: count })
    });

    const data = await response.json();

    if (!data.success) {
        throw new Error(`API error: ${data.error}`);
    }

    return data.data;
}

// Usage example
generateID('http://localhost:8080', 1)
    .then(id => console.log('Generated ID:', id.id))
    .catch(err => console.error('Error:', err));
```

## Troubleshooting

### Common Issues

1. **Service fails to start**
   - Check if port is already in use
   - Verify NODE_ID is in range 0-3

2. **ID generation fails**
   - Check service health status: `curl http://localhost:8080/health`
   - Verify request parameters are correct

3. **Performance issues**
   - Monitor system resource usage
   - Consider adding more service instances

### Log Analysis

The service uses quiet mode to reduce log output, but critical errors are still logged. It's recommended to configure external monitoring tools to collect metrics.

## License

This example code follows the same MIT license as the main project.