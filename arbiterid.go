package arbiterid

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

// Constants for bit allocation in the ID
const (
	// Epoch is the starting timestamp (2025-01-01T00:00:00.000Z)
	Epoch int64 = 1735718400000

	// Bit allocation for different sections of the ID
	TypeBits      uint8 = 10 // 10 bits for type
	TimestampBits uint8 = 41 // 41 bits for timestamp (milliseconds since Epoch)
	NodeBits      uint8 = 2  // 2 bits for node ID (0-3)
	SeqBits       uint8 = 10 // 10 bits for sequence (0-1023)

	// Total bits: 10 (Type) + 41 (Timestamp) + 2 (Node) + 10 (Sequence) = 63 bits.
	// This leaves the most significant bit (sign bit for int64) as 0, ensuring positive IDs.

	// Calculate maximum values for each section
	TypeMax      uint16 = (1 << TypeBits) - 1 // Max value for 10 bits (0-1023)
	TimestampMax int64  = (1 << TimestampBits) - 1
	NodeMax      int64  = (1 << NodeBits) - 1
	SeqMax       int64  = (1 << SeqBits) - 1

	// Calculate bit shifts for each section
	SeqShift  = 0
	NodeShift = SeqBits
	TimeShift = NodeShift + NodeBits
	TypeShift = TimeShift + TimestampBits

	// Masks for extracting parts of the ID
	SeqMask       int64 = SeqMax << SeqShift
	NodeMask      int64 = NodeMax << NodeShift
	TimestampMask int64 = TimestampMax << TimeShift
	TypeMask      int64 = int64(TypeMax) << TypeShift // Cast TypeMax for int64 mask

	// Configuration for clock rollover protection
	maxRolloverWaitAttempts   = 2000
	rolloverWaitCheckInterval = time.Microsecond * 50
	maxEarlyAttempts          = 10 // Maximum attempts to check for fresh time
)

// Encoding maps for Base32 and Base58
const (
	encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"
	encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

// Error definitions
var (
	ErrInvalidNodeID         = errors.New("arbiterid: invalid node ID")
	ErrInvalIDType           = errors.New("arbiterid: ID type must be between 0 and 1023") // Updated for 10 bits
	ErrInvalidBase58         = errors.New("arbiterid: invalid base58 string")
	ErrInvalidBase32         = errors.New("arbiterid: invalid base32 string")
	ErrMonotonicityViolation = errors.New("arbiterid: generated ID is not strictly greater than the last ID")
	ErrClockNotAdvancing     = errors.New("arbiterid: system clock appears to be stuck or moving backward excessively")
	ErrBase64InvalidLength   = errors.New("arbiterid: invalid base64 ID length, expected 8 decoded bytes")
)

// Decoding maps, initialized in init()
var (
	decodeBase32Map [256]byte
	decodeBase58Map [256]byte
)

func init() {
	for i := range decodeBase32Map {
		decodeBase32Map[i] = 0xFF
	}
	for i := range decodeBase58Map {
		decodeBase58Map[i] = 0xFF
	}
	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}
}

// ID represents an arbiterid unique identifier
type ID int64

// IDType is the type used for different categories of IDs, now supporting up to 1023.
type IDType uint16

// Node generates and manages unique IDs
type Node struct {
	mu                       sync.Mutex
	epoch                    time.Time
	lastID                   ID
	node                     int64
	time                     int64
	seq                      int64
	clockWarningCount        int64
	strictMonotonicityChecks bool
	quietMode                bool // Suppresses most log output for testing
}

// NodeOption is a functional option for configuring a Node
type NodeOption func(*Node)

// WithStrictMonotonicityCheck enables or disables checking that new IDs are strictly greater than the last.
// Default is true. If you want to allow different ID types, set to false.
func WithStrictMonotonicityCheck(enable bool) NodeOption {
	return func(n *Node) {
		n.strictMonotonicityChecks = enable
	}
}

// WithQuietMode enables or disables quiet mode to suppress most log output.
// Default is false. Set to true to reduce logging during testing or high-volume environments.
func WithQuietMode(enable bool) NodeOption {
	return func(n *Node) {
		n.quietMode = enable
	}
}

// NewNode creates a new Node for generating IDs with the given options
func NewNode(nodeID int, options ...NodeOption) (*Node, error) {
	if int64(nodeID) < 0 || int64(nodeID) > NodeMax {
		return nil, fmt.Errorf("%w: got %d, max %d", ErrInvalidNodeID, nodeID, NodeMax)
	}

	epochTime := time.Unix(Epoch/1000, (Epoch%1000)*1000000).UTC()

	n := &Node{
		node:                     int64(nodeID),
		epoch:                    epochTime,
		time:                     0,
		seq:                      0,
		lastID:                   0,
		strictMonotonicityChecks: true,
		clockWarningCount:        0,
	}

	for _, option := range options {
		option(n)
	}
	if !n.quietMode {
		log.Printf("ArbiterID Node initialized: ID=%d, StrictMonotonicityChecks=%t, QuietMode=%t", n.node, n.strictMonotonicityChecks, n.quietMode)
	}
	return n, nil
}

// Generate creates a new unique ID with the given type and current timestamp.
// This method includes clock rollover detection for production safety.
func (n *Node) Generate(idType IDType) (ID, error) {
	if uint16(idType) > TypeMax {
		return 0, fmt.Errorf("%w: got %d, max %d", ErrInvalIDType, idType, TypeMax)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UTC().Sub(n.epoch).Milliseconds()

	// Clock rollover detection - only for Generate() using real time
	if now < n.time {
		// Only treat as significant clock backward movement if >1ms backwards
		// This avoids false warnings from minor time source variations in tight loops
		if now < n.time-1 {
			// Log significant clock backwards movement (rare and indicates system issues)
			if !n.quietMode {
				n.clockWarningCount++
				log.Printf("ArbiterID Warning: Clock moved backwards significantly. Current time: %d, Last time: %d. Using last time. (Warning #%d)", now, n.time, n.clockWarningCount)
			} else {
				n.clockWarningCount++
			}
		}
		// Always use the last time when clock appears to go backwards
		now = n.time
	}

	// Handle sequence rollover with real time - can wait and advance
	if now == n.time {
		n.seq = (n.seq + 1) & SeqMax
		if n.seq == 0 {
			// Sequence exhausted, need to wait for next millisecond
			originalTime := n.time
			attempts := 0
			for now <= originalTime {
				attempts++
				if attempts > maxRolloverWaitAttempts {
					if !n.quietMode {
						log.Printf("ArbiterID Critical: Clock appears stuck at %dms after %d attempts. Node ID: %d", now, attempts, n.node)
					}
					return 0, fmt.Errorf("%w: clock stuck at %dms after %d attempts from %dms",
						ErrClockNotAdvancing, now, attempts, originalTime)
				}

				time.Sleep(rolloverWaitCheckInterval)
				// Get fresh time and check if it has advanced
				freshTime := time.Now().UTC().Sub(n.epoch).Milliseconds()
				if freshTime > originalTime {
					now = freshTime
					break
				}
				now = freshTime
			}
		}
	} else {
		n.seq = 0
	}

	return n.generateInternal(idType, now)
}

// GenerateWithTimestamp creates a new unique ID with the given type and specific timestamp.
// This method does NOT include clock rollover detection - it uses the provided timestamp as-is.
// Use this for testing or when you need precise timestamp control.
func (n *Node) GenerateWithTimestamp(idType IDType, timestamp time.Time) (ID, error) {
	if uint16(idType) > TypeMax {
		return 0, fmt.Errorf("%w: got %d, max %d", ErrInvalIDType, idType, TypeMax)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	now := timestamp.UTC().Sub(n.epoch).Milliseconds()

	// Handle sequence management for fixed timestamp
	if now == n.time {
		n.seq = (n.seq + 1) & SeqMax
		if n.seq == 0 {
			// Sequence exhausted - cannot advance time with fixed timestamp
			return 0, fmt.Errorf("%w: sequence exhausted for timestamp %dms, cannot advance time with fixed timestamp",
				ErrClockNotAdvancing, now)
		}
	} else {
		n.seq = 0
	}

	return n.generateInternal(idType, now)
}

// generateInternal handles the core ID generation logic.
// Assumes sequence management and time advancement have been handled by the caller.
// The 'now' parameter should be the timestamp in milliseconds since epoch.
func (n *Node) generateInternal(idType IDType, now int64) (ID, error) {

	n.time = now

	if now > TimestampMax {
		if !n.quietMode {
			log.Printf("ArbiterID Critical: Timestamp %dms has overflowed TimestampMax %dms. Node ID: %d", now, TimestampMax, n.node)
		}
		return 0, fmt.Errorf("arbiterid: timestamp %dms has overflowed maximum %dms (Epoch %s, ~69 years)",
			now, TimestampMax, n.epoch.Format(time.RFC3339))
	}

	id := ID(
		(int64(idType) << TypeShift) |
			(now << TimeShift) |
			(n.node << NodeShift) |
			n.seq,
	)

	if n.strictMonotonicityChecks && id <= n.lastID {
		if !n.quietMode {
			log.Printf("ArbiterID Critical: Monotonicity violation. New ID %d <= Last ID %d. Node ID: %d. Time: %d, Seq: %d", id, n.lastID, n.node, n.time, n.seq)
		}
		return 0, fmt.Errorf("%w: new ID %d (%s) <= last ID %d (%s). Time: %dms, Seq: %d",
			ErrMonotonicityViolation, id, id.TimeISO(), n.lastID, n.lastID.TimeISO(), n.time, n.seq)
	}

	n.lastID = id
	return id, nil
}

// GenerateSimple is a convenience method that generates an ID and panics on error.
func (n *Node) GenerateSimple(idType IDType) ID {
	id, err := n.Generate(idType)
	if err != nil {
		if !n.quietMode {
			log.Panicf("ArbiterID: Failed to generate ID for type %d: %v", idType, err)
		}
		panic(fmt.Sprintf("ArbiterID: Failed to generate ID for type %d: %v", idType, err))
	}
	return id
}

// LastID returns the last ID generated by this node
func (n *Node) LastID() ID {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.lastID
}

// Int64 returns the ID as a raw int64
func (id ID) Int64() int64 {
	return int64(id)
}

// String returns the ID as a decimal string
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// Components extracts and returns all components of the ID.
// Timestamp returned is milliseconds since Unix epoch.
func (id ID) Components() (idType IDType, timestampMillisUnix int64, node int64, seq int64) {
	timestampMillisNodeEpoch := (int64(id) & TimestampMask) >> TimeShift
	timestampMillisUnix = timestampMillisNodeEpoch + Epoch
	idType = IDType(id.Type())
	return idType, timestampMillisUnix, id.Node(), id.Seq()
}

// Type returns the type component of the ID as int64.
func (id ID) Type() int64 {
	return (int64(id) & TypeMask) >> TypeShift
}

// Time returns the timestamp in Unix milliseconds
func (id ID) Time() int64 {
	return ((int64(id) & TimestampMask) >> TimeShift) + Epoch
}

// TimeTime returns the timestamp as a time.Time object in UTC
func (id ID) TimeTime() time.Time {
	return time.Unix(0, id.Time()*int64(time.Millisecond)).UTC()
}

// TimeISO returns the timestamp in ISO 8601 format (UTC)
func (id ID) TimeISO() string {
	return id.TimeTime().Format(time.RFC3339Nano)
}

// Node returns the node component of the ID
func (id ID) Node() int64 {
	return (int64(id) & NodeMask) >> NodeShift
}

// Seq returns the sequence component of the ID
func (id ID) Seq() int64 {
	return int64(id) & SeqMask
}

// ParseString converts a decimal string to an ID
func ParseString(s string) (ID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("arbiterid: failed to parse decimal string '%s': %w", s, err)
	}
	return ID(i), nil
}

// Base2 returns the ID as a base2 string (binary representation)
func (id ID) Base2() string {
	return fmt.Sprintf("%063b", int64(id)) // Pad to 63 bits for consistency
}

// ParseBase2 converts a Base2 string into an ID
func ParseBase2(s string) (ID, error) {
	i, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("arbiterid: failed to parse base2 string '%s': %w", s, err)
	}
	return ID(i), nil
}

// Base32 returns the ID as a base32 string.
func (id ID) Base32() string {
	if id == 0 {
		return string(encodeBase32Map[0])
	}
	n := uint64(id)
	buf := make([]byte, 13) // Max 13 chars for 63 bits (63/5 ~ 12.6)
	i := 12
	for n > 0 {
		buf[i] = encodeBase32Map[n%32]
		n /= 32
		i--
	}
	return string(buf[i+1:])
}

// ParseBase32 converts a base32 string to an ID
func ParseBase32(s string) (ID, error) {
	var val uint64
	if len(s) == 0 {
		return 0, fmt.Errorf("%w: input string is empty", ErrInvalidBase32)
	}
	if len(s) > 13 {
		return 0, fmt.Errorf("%w: input string '%s' too long (max 13 chars)", ErrInvalidBase32, s)
	}
	for i := 0; i < len(s); i++ {
		char := s[i]
		decodedByte := decodeBase32Map[char]
		if decodedByte == 0xFF {
			return 0, fmt.Errorf("%w: invalid char '%c' in '%s'", ErrInvalidBase32, char, s)
		}
		if val > (math.MaxUint64-uint64(decodedByte))/32 {
			return 0, fmt.Errorf("%w: value '%s' overflows uint64", ErrInvalidBase32, s)
		}
		val = val*32 + uint64(decodedByte)
	}
	if val > math.MaxInt64 { // Ensure it fits in positive int64
		return 0, fmt.Errorf("%w: value '%s' overflows positive int64", ErrInvalidBase32, s)
	}
	return ID(val), nil
}

// Base58 returns the ID as a base58 string.
func (id ID) Base58() string {
	if id == 0 {
		return string(encodeBase58Map[0])
	}
	n := uint64(id)
	buf := make([]byte, 11) // Max 11 chars for 63 bits (63/log2(58) ~ 10.7)
	i := 10
	for n > 0 {
		buf[i] = encodeBase58Map[n%58]
		n /= 58
		i--
	}
	return string(buf[i+1:])
}

// ParseBase58 converts a base58 string to an ID
func ParseBase58(s string) (ID, error) {
	var val uint64
	if len(s) == 0 {
		return 0, fmt.Errorf("%w: input string is empty", ErrInvalidBase58)
	}
	if len(s) > 11 {
		return 0, fmt.Errorf("%w: input string '%s' too long (max 11 chars)", ErrInvalidBase58, s)
	}
	for i := 0; i < len(s); i++ {
		char := s[i]
		decodedByte := decodeBase58Map[char]
		if decodedByte == 0xFF {
			return 0, fmt.Errorf("%w: invalid char '%c' in '%s'", ErrInvalidBase58, char, s)
		}
		if val > (math.MaxUint64-uint64(decodedByte))/58 {
			return 0, fmt.Errorf("%w: value '%s' overflows uint64", ErrInvalidBase58, s)
		}
		val = val*58 + uint64(decodedByte)
	}
	if val > math.MaxInt64 { // Ensure it fits in positive int64
		return 0, fmt.Errorf("%w: value '%s' overflows positive int64", ErrInvalidBase58, s)
	}
	return ID(val), nil
}

// Base64 returns the ID as a URL-safe base64 string.
func (id ID) Base64() string {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(id))
	return base64.RawURLEncoding.EncodeToString(buf[:])
}

// ParseBase64 converts a URL-safe base64 string to an ID.
func ParseBase64(s string) (ID, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, fmt.Errorf("arbiterid: failed to decode base64 string '%s': %w", s, err)
	}
	if len(b) != 8 {
		return 0, fmt.Errorf("%w: decoded data len %d for '%s'", ErrBase64InvalidLength, len(b), s)
	}
	// The ID is 63-bit, so the MSB of the uint64 must be 0.
	val := binary.BigEndian.Uint64(b)
	if val > math.MaxInt64 {
		return 0, fmt.Errorf("arbiterid: base64 value '%s' (%d) overflows positive int64 (max %d)", s, val, int64(math.MaxInt64))
	}
	return ID(val), nil
}

// MarshalJSON implements json.Marshaler
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(id), 10) + `"`), nil
}

// JSONSyntaxError is returned when an ID cannot be unmarshaled from JSON
type JSONSyntaxError struct{ Original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("arbiterid: invalid ID JSON format: %s", string(j.Original))
}

// UnmarshalJSON implements json.Unmarshaler
func (id *ID) UnmarshalJSON(b []byte) error {
	s := string(b)
	var val int64
	var err error

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		val, err = strconv.ParseInt(s[1:len(s)-1], 10, 64)
	} else {
		// Allow parsing as raw number for flexibility, though string is preferred.
		val, err = strconv.ParseInt(s, 10, 64)
	}

	if err != nil {
		return fmt.Errorf("arbiterid: failed to parse ID from JSON %s: %w", s, err)
	}
	if val < 0 { // Arbiter IDs are positive
		return fmt.Errorf("arbiterid: parsed JSON ID %d is negative, expected positive value from %s", val, s)
	}
	*id = ID(val)
	return nil
}
