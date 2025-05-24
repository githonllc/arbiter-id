package main

import (
	"fmt"
	"log"
	"time"

	"github.com/githonllc/arbiterid"
)

func main() {
	// Initialize a new node.
	// Node ID must be unique for each server instance (0-3).
	// Disable strict monotonicity checks for this demo to allow different ID types
	node, err := arbiterid.NewNode(0, arbiterid.WithStrictMonotonicityCheck(false))
	if err != nil {
		log.Fatalf("Failed to create arbiterid node: %v", err)
	}

	// Define ID types (0-1023)
	const UserIDType arbiterid.IDType = 1
	const PostIDType arbiterid.IDType = 512
	const CommentIDType arbiterid.IDType = 256

	fmt.Println("=== ArbiterID Examples ===")

	// Example 1: Generate User ID
	fmt.Println("\n1. Generating User ID:")
	userID, err := node.Generate(UserIDType)
	if err != nil {
		log.Fatalf("Failed to generate UserID: %v", err)
	}
	fmt.Printf("   Raw int64: %d\n", userID.Int64())
	fmt.Printf("   String:    %s\n", userID.String())
	fmt.Printf("   Base58:    %s\n", userID.Base58())
	fmt.Printf("   Base64:    %s\n", userID.Base64())
	fmt.Printf("   Base32:    %s\n", userID.Base32())
	fmt.Printf("   ISO Time:  %s\n", userID.TimeISO())

	// Example 2: Generate Post ID
	fmt.Println("\n2. Generating Post ID:")
	postID, err := node.Generate(PostIDType)
	if err != nil {
		log.Fatalf("Failed to generate PostID: %v", err)
	}
	fmt.Printf("   Raw int64: %d\n", postID.Int64())
	fmt.Printf("   Base64:    %s\n", postID.Base64())

	// Example 3: Extract components
	fmt.Println("\n3. Extracting ID Components:")
	IDType, tsMillis, nodeID, seq := userID.Components()
	fmt.Printf("   Type:      %d\n", IDType)
	fmt.Printf("   Timestamp: %d ms (since epoch)\n", tsMillis)
	fmt.Printf("   Node ID:   %d\n", nodeID)
	fmt.Printf("   Sequence:  %d\n", seq)

	// Example 4: Generate multiple IDs in sequence
	fmt.Println("\n4. Generating Multiple IDs:")
	for i := 0; i < 5; i++ {
		id, err := node.Generate(CommentIDType)
		if err != nil {
			log.Fatalf("Failed to generate ID %d: %v", i, err)
		}
		fmt.Printf("   ID %d: %s (seq: %d)\n", i+1, id.Base58(), id.Seq())
	}

	// Example 5: Parse IDs from different encodings
	fmt.Println("\n5. Parsing IDs from Strings:")

	// Parse from decimal string
	parsedFromString, err := arbiterid.ParseString(userID.String())
	if err != nil {
		log.Fatalf("Failed to parse from string: %v", err)
	}
	fmt.Printf("   Parsed from string: %s -> %d (match: %t)\n",
		userID.String(), parsedFromString.Int64(), parsedFromString == userID)

	// Parse from Base58
	parsedFromBase58, err := arbiterid.ParseBase58(userID.Base58())
	if err != nil {
		log.Fatalf("Failed to parse from Base58: %v", err)
	}
	fmt.Printf("   Parsed from Base58: %s -> %d (match: %t)\n",
		userID.Base58(), parsedFromBase58.Int64(), parsedFromBase58 == userID)

	// Parse from Base64
	parsedFromBase64, err := arbiterid.ParseBase64(userID.Base64())
	if err != nil {
		log.Fatalf("Failed to parse from Base64: %v", err)
	}
	fmt.Printf("   Parsed from Base64: %s -> %d (match: %t)\n",
		userID.Base64(), parsedFromBase64.Int64(), parsedFromBase64 == userID)

	// Example 6: Generate ID with specific timestamp
	fmt.Println("\n6. Generating ID with Specific Timestamp:")
	specificTime := time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC)
	timestampID, err := node.GenerateWithTimestamp(UserIDType, specificTime)
	if err != nil {
		log.Fatalf("Failed to generate ID with timestamp: %v", err)
	}
	fmt.Printf("   Generated at %s\n", specificTime.Format(time.RFC3339))
	fmt.Printf("   ID: %s\n", timestampID.Base58())
	fmt.Printf("   Extracted time: %s\n", timestampID.TimeISO())

	// Example 7: Demonstrate K-sortable property
	fmt.Println("\n7. Demonstrating K-sortable Property:")
	var ids []arbiterid.ID
	for i := 0; i < 3; i++ {
		id, err := node.Generate(UserIDType)
		if err != nil {
			log.Fatalf("Failed to generate ID: %v", err)
		}
		ids = append(ids, id)
		time.Sleep(time.Millisecond) // Small delay to ensure different timestamps
	}

	fmt.Println("   Generated IDs (should be in ascending order):")
	for i, id := range ids {
		fmt.Printf("   %d: %d (time: %s)\n", i+1, id.Int64(), id.TimeISO())
	}

	fmt.Println("\n=== Examples Complete ===")
}
