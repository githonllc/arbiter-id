#!/bin/bash

# ArbiterID Service API Test Script

BASE_URL="http://localhost:8080"

echo "=== ArbiterID Service API Test ==="
echo

# Check if service is running
echo "1. Check service health status..."
curl -s "${BASE_URL}/health" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/health"
echo
echo

# Get service information
echo "2. Get service information..."
curl -s "${BASE_URL}/info" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/info"
echo
echo

# Generate default ID
echo "3. Generate default ID (type 0)..."
curl -s -X POST "${BASE_URL}/generate" | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/generate"
echo
echo

# Generate specific type ID
echo "4. Generate user ID (type 1)..."
curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 1}' | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 1}'
echo
echo

# Generate ID using query parameters
echo "5. Generate order ID using query parameters (type 100)..."
curl -s -X POST "${BASE_URL}/generate?type=100" | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/generate?type=100"
echo
echo

# Batch generate IDs
echo "6. Batch generate 3 comment IDs (type 200)..."
curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 200, "count": 3}' | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 200, "count": 3}'
echo
echo

# Test error cases
echo "7. Test invalid ID type..."
curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 9999}' | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/generate" \
  -H "Content-Type: application/json" \
  -d '{"id_type": 9999}'
echo
echo

# Test invalid HTTP method
echo "8. Test invalid HTTP method..."
curl -s -X GET "${BASE_URL}/generate" | jq '.' 2>/dev/null || curl -s -X GET "${BASE_URL}/generate"
echo
echo

echo "=== Test Complete ==="